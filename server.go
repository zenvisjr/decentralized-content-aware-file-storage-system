package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/zenvisjr/distributed-file-storage-system/p2p"
)

// FileServerOps represents the options for a file server.
type FileServerOps struct {
	ID                string
	RootStorage       string
	PathTransformFunc PathTransformFunc
	Transort          p2p.Transport
	BootstrapNodes    []string
	EncKey            []byte
}

// FileServer represents a file server.
type FileServer struct {
	FileServerOps
	store        *Store
	quitch       chan struct{}
	peerLock     sync.Mutex
	peers        map[string]p2p.Peer
	notFoundChan chan struct{}

	// incomingStreamChan chan p2p.Peer
}

// NewFileServer creates a new FileServer instance.
func NewFileServer(ops FileServerOps) (*FileServer, error) {
	storeOps := &StoreOps{
		Root:              ops.RootStorage,
		PathTransformFunc: ops.PathTransformFunc,
	}
	// if len(ops.StorageRoot) == 0 {
	// 	ops.StorageRoot = defaultRoot
	// }
	if len(ops.ID) == 0 {
		ops.ID = generateID()
	}

	return &FileServer{
		FileServerOps: ops,
		store:         NewStore(*storeOps),
		quitch:        make(chan struct{}),
		peerLock:      sync.Mutex{},
		peers:         make(map[string]p2p.Peer),
		notFoundChan:  make(chan struct{}, 100),
	}, nil
}

type Message struct {
	// From    string
	Payload any
}

type MessageStoreFile struct {
	Key string
	ID  string
	// Ext  string
	Size      int64
	Signature []byte
}

type MessageGetFile struct {
	Key string
	ID  string
}

type MessageDeleteFile struct {
	Key string
	ID  string
}

type MessageGetFileNotFound struct {
	Key string
	ID  string
}

func (f *FileServer) Store(key string, r io.Reader) error {

	//1. store the file to disk
	//2. broadcast the file to other peers in the network

	var (
		fileBuffer bytes.Buffer
		tee        = io.TeeReader(r, &fileBuffer)
	)
	// fileHash, size, err := hashFileContent(fileBuffer)
	// if err != nil {
	// 	return err
	// }

	// if keyExist, ok := f.store.HashMap[fileHash]; !ok {
	// f.store.HashMap[fileHash] = key

	fmt.Println("Storing file to disk")
	size, err := f.store.Write(f.ID, key, tee)
	if err != nil {
		return err
	}
	fmt.Printf("Original file size stored to disk: %d bytes\n", size)

	// } else {
	// 	fmt.Println("File already exists with key", keyExist)
	// }

	privKey, err := p2p.LoadPrivateKey()
	if err != nil {
		return err
	}

	var teeBuf bytes.Buffer
	// Now read from fileBuffer (convert to reader) and tee to teeBuf
	fileReader := bytes.NewReader(fileBuffer.Bytes())
	teeReader := io.TeeReader(fileReader, &teeBuf)

	// fmt.Printf("Original file size before encryption: %d bytes\n", len(fileBuffer.Bytes()))
	// fmt.Println("Original file hash:", sha256.Sum256(fileBuffer.Bytes()))

	_, err = copyEncrypt(&teeBuf, teeReader, f.EncKey)
	if err != nil {
		return err
	}
	// fmt.Printf("Size returned by copyEncrypt: %d bytes\n", n)
	// fmt.Printf("Actual teeBuf size: %d bytes\n", len(teeBuf.Bytes()))
	// fmt.Printf("Encrypted data hash (IV + encrypted): %x\n", sha256.Sum256(teeBuf.Bytes()))

	signature, err := signSignature(teeBuf.Bytes(), privKey) // fileBuffer already contains full file
	if err != nil {
		return err
	}
	//store the signature of encrypted file in the map
	sigKey := hashKey(key) + getExtension(key)
	f.store.SaveSignature(sigKey, signature)
	// fmt.Printf("YOYOYO [%s] Signature stored in map for key [%s]", f.Transort.ListenAddr(), sigKey)

	msg := Message{
		Payload: MessageStoreFile{
			Key:  sigKey,
			Size: int64(len(teeBuf.Bytes())),
			ID:   f.ID,
			// Ext:  getExtension(key),
			Signature: signature,
		},
	}

	fmt.Printf("Broadcasting file to other peers")
	fmt.Printf("Message size: %d bytes\n", msg.Payload.(MessageStoreFile).Size)

	if err := f.broadcast(&msg); err != nil {
		return err
	}

	time.Sleep(50 * time.Millisecond)

	peers := []io.Writer{}
	for _, peer := range f.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	fmt.Println("Streaming file to other peers")
	mw.Write([]byte{p2p.IncommingStream})
	// Create a fresh reader from the original file data for streaming
	// Stream the encrypted data (IV + encrypted data) to peers
	encryptedReader := bytes.NewReader(teeBuf.Bytes())

	nw, err := io.Copy(mw, encryptedReader)
	if err != nil {
		return err
	}

	fmt.Printf("[%s] Recieved and written (%d) bytes to disk\n", f.Transort.ListenAddr(), nw)

	return nil
}

// Get fetches the file from disk if present locally else fetches from the network from its peers
// we continue reading from each peer until we successfully fetch a file
// once fetched from any peer we store it to disk and return it
// if we are unable to fetch the file from any peer we return error
// if file is not present with any peer we return error
func (f *FileServer) Get(key string) (io.Reader, string, error) {
	if f.store.Has(f.ID, key) {
		fmt.Printf("[%s] have file [%s], serving from local disk\n", f.Transort.ListenAddr(), key)
		_, r, err, fileLocation := f.store.Read(f.ID, key)
		if err != nil {
			return nil, "", err
		}
		fmt.Printf("[%s] File location %s\n", f.Transort.ListenAddr(), fileLocation)

		if rc, ok := r.(io.ReadCloser); ok {
			defer rc.Close() // ✅ this ensures file is closed after use
		}

		return r, fileLocation, nil
	}
	fmt.Printf("[%s] dont have file [%s], fetching from network...\n", f.Transort.ListenAddr(), key)

	msg := Message{
		Payload: MessageGetFile{
			Key: hashKey(key) + getExtension(key),
			ID:  f.ID,
		},
	}
	err := f.broadcast(&msg)
	if err != nil {
		return nil, "", err
	}

	time.Sleep(500 * time.Millisecond)
	noOfPeers := len(f.peers)
	notFound := 0 //tracks peers that explicitly said “I don’t have it.”
	received := 0 //counts any kind of peer response — stream or not-found

	//why recieved? so that we can count all responses and if responses == no of peers
	//it means we were unable to fetch the file from any peer

	for {
		fmt.Println("receiving from peer")

		select {
		case <-f.notFoundChan:
			fmt.Println("Peer responded: not found")
			notFound++
			received++ // we will count response even if stream fails later
			if notFound == noOfPeers {
				return nil, "", fmt.Errorf("[%s] and all its peers don't have file [%s]", f.Transort.ListenAddr(), key)
			}

		case streamPeer := <-f.Transort.ConsumeStream():
			fmt.Println("Peer responded: stream")
			received++
			// Only process if this is the same peer
			// if streamPeer.RemoteAddr().String() != peer.RemoteAddr().String() {
			// 	fmt.Printf("Skipping peer %s, stream came from %s\n", peer.RemoteAddr().String(), streamPeer.RemoteAddr().String())
			// 	continue
			// }

			var (
				fileSize int64
				sigLen   int64
			)

			if err := binary.Read(streamPeer, binary.LittleEndian, &fileSize); err != nil {
				fmt.Println("failed to read fileSize from stream, skipping peer")
				continue
			}

			if err := binary.Read(streamPeer, binary.LittleEndian, &sigLen); err != nil {
				fmt.Println("failed to read sigLen from stream, skipping peer")
				continue
			}

			signature := make([]byte, sigLen)
			if _, err := io.ReadFull(streamPeer, signature); err != nil {
				fmt.Println("failed to read signature from stream, skipping peer")
				continue
			}
			fmt.Printf(" [%s] -> [%s]", msg.Payload.(MessageGetFile).Key, signature)

			encData := make([]byte, fileSize)
			if _, err := io.ReadFull(streamPeer, encData); err != nil {
				fmt.Println("failed to read encrypted data from stream, skipping peer")
				continue
			}

			peerAddr := streamPeer.RemoteAddr().String()
			peerPublicKey, ok := p2p.GetPeerPublicKey(peerAddr)
			if !ok {
				fmt.Printf("Peer %s public key not found, skipping\n", peerAddr)
				continue
			}

			if err := verifySignature(encData, signature, peerPublicKey); err != nil {
				fmt.Println("failed to verify signature, skipping peer")
				continue
			}
			fmt.Printf("Signature verified from peer %s\n", peerAddr)

			encryptedReader := bytes.NewReader(encData)

			n, err := f.store.WriteDecrypted(f.ID, key, f.EncKey, encryptedReader)
			if err != nil {
				fmt.Println("failed to write file to disk after successful verification, skipping peer")
				continue
			}
			fmt.Printf("[%s] received (%d) bytes from peer %s\n", f.Transort.ListenAddr(), n, streamPeer.RemoteAddr().String())
			streamPeer.CloseStream()

			//now we have the file on disk, if we are unable to read it we will simply return error
			_, r, err, fileLocation := f.store.Read(f.ID, key)
			if err != nil {
				return nil, "", err
			}
			// fmt.Printf("[%s] File location %s\n", f.Transort.ListenAddr(), fileLocation)
			return r, fileLocation, nil

		case <-time.After(2 * time.Second):
			// Timeout waiting for this peer — skip
			if received == noOfPeers {
				return nil, "", fmt.Errorf("[%s] peers were not able to fetch file [%s]", f.Transort.ListenAddr(), key)
			}
			continue
		}

	}
	// return nil, "", nil
}

// Delete delets the file locally if its present according to key
// and also from all the peers in the network
func (f *FileServer) Delete(key string) error {
	if f.store.Has(f.ID, key) {
		fmt.Printf("[%s] have file [%s], deleting from local disk\n", f.Transort.ListenAddr(), key)
		if err := f.store.Delete(f.ID, key); err != nil {
			return err
		}
		f.store.DeleteSignature(key)
		fmt.Println("deleted signature from local disk")
	}

	fmt.Printf("[%s] will now delete files from all its peers\n", f.Transort.ListenAddr())

	msg := Message{
		Payload: MessageDeleteFile{
			Key: hashKey(key) + getExtension(key),
			ID:  f.ID,
		},
	}

	err := f.broadcast(&msg)
	if err != nil {
		return err
	}

	time.Sleep(500 * time.Millisecond)

	return nil
}

// DeleteLocal deletes the file locally if its present according to key
func (f *FileServer) DeleteLocal(key string) error {
	if f.store.Has(f.ID, key) {
		fmt.Printf("[%s] have file [%s], deleting from local disk\n", f.Transort.ListenAddr(), key)
		err := f.store.Delete(f.ID, key)
		if err != nil {
			return err
		}
		f.store.DeleteSignature(key)
	}
	return nil
}

// Start starts the file server.
func (f *FileServer) Start() error {
	if err := f.Transort.ListenAndAccept(); err != nil {
		return err
	}
	// fmt.Println("File server started")

	if len(f.BootstrapNodes) > 0 {
		fmt.Println("Bootstraping network......")
		err := f.BootstrapNetwork()
		if err != nil {
			return err
		}
	}
	f.loop()

	return nil
}

func (f *FileServer) BootstrapNetwork() error {
	for _, addr := range f.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}

		go func(addr string) {
			if err := f.Transort.Dial(addr); err != nil {
				log.Printf("failed to connect to bootstrap node %s", addr)
			} else {
				fmt.Println("Connected to bootstrap node", addr)
			}
		}(addr)
	}
	return nil
}

func (f *FileServer) loop() {
	// fmt.Println("File server loop started")
	defer func() {
		log.Println("File server stopped due to error or user QUIT action")
		f.Transort.Close()
	}()

	for {
		select {
		case rpc := <-f.Transort.Consume():
			// fmt.Println("Received RPC", rpc)
			// fmt.Println("RPC Payload", rpc.Payload[0])
			// if rpc.Payload[0] == p2p.IncommingStream {
			// 	fmt.Println("RPC Payload", rpc.Payload[0])
			// if rpc.Stream {
			// 	peer := f.peers[rpc.From]
			// f.incomingStreamChan <- f.peers[rpc.From]
			// 	continue
			// }
			// fmt.Printf("Received %+v\n", msg)
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				// log.Println("decoding error", err)
				continue
			}
			if err := f.HandleMessage(rpc.From, &msg); err != nil {
				log.Println("handling type of message error", err)
				continue
			}

		case <-f.quitch:
			fmt.Println("Quitting")
			return
		}
	}
}

func (f *FileServer) HandleMessage(from string, msg *Message) error {
	switch v := msg.Payload.(type) {
	case MessageStoreFile:
		return f.handleMessageStoreFile(from, &v)
	case MessageGetFile:
		return f.handleMessageGetFile(from, &v)
	case MessageDeleteFile:
		return f.handleMessageDeleteFile(from, &v)
	case MessageGetFileNotFound:
		return f.handleMessageGetFileNotFound()

	}
	return nil
}

func (f *FileServer) handleMessageGetFile(from string, msg *MessageGetFile) error {
	// fmt.Println("executing handleMessageGetFile")
	if !f.store.Has(msg.ID, msg.Key) {
		nack := Message{
			Payload: MessageGetFileNotFound{
				Key: msg.Key,
				ID:  msg.ID,
			},
		}
		peer := f.peers[from]
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(nack)
		if err != nil {
			return err
		}
		peer.Send([]byte{p2p.IncommingMessage})
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
		return fmt.Errorf("[%s] need to serve file (%s) to peer %s but it does not exist on the network", f.Transort.ListenAddr(), msg.Key, from)

	}

	fmt.Printf("[%s] fetching file (%s) from network and sending it via wire to peer %s\n", f.Transort.ListenAddr(), msg.Key, from)

	fileSize, rd, err, _ := f.store.Read(msg.ID, msg.Key)
	if err != nil {
		return err
	}

	r, ok := rd.(io.ReadCloser)
	if ok {
		// fmt.Println("Closing reader")
		defer r.Close()
	}

	peer, ok := f.peers[from]
	if !ok {
		return fmt.Errorf("peer %s could not be found in map", from)
	}

	signature, err := f.store.GetSignature(msg.Key)
	if err != nil {
		return fmt.Errorf("signature for file %s could not be found in map on peer [%s]", msg.Key, from)
	}
	fmt.Println("Signature for file in handleMessageGetFile", msg.Key, "is", signature)
	// First send the "incomingStream" byte to the peer and then we can send
	// the file size as an int64.

	peer.Send([]byte{p2p.IncommingStream})
	binary.Write(peer, binary.LittleEndian, fileSize)
	binary.Write(peer, binary.LittleEndian, int64(len(signature)))
	peer.Write(signature)

	// fmt.Println("Received RPC", rpc)
	// fmt.Println("RPC Payload", rpc.Payload[0])
	// if rpc.Payload[0] == p2p.IncommingStream {
	// 	fmt.Println("RPC Payload", rpc.Payload[0])
	// peer := f.peers[from]

	// 	continue
	// }

	n, err := io.Copy(peer, rd)
	if err != nil {
		return err
	}
	// f.incomingStreamChan <- peer
	fmt.Printf("[%s] Written (%d) bytes to peer %s over network\n", f.Transort.ListenAddr(), n, from)
	return nil

}

func (f *FileServer) handleMessageStoreFile(from string, msg *MessageStoreFile) error {
	// panic("not implemented")
	// fmt.Printf("Received data message %+v from %s\n", msg, from)
	peer, ok := f.peers[from]
	if !ok {
		return fmt.Errorf("peer %s could not be found", from)
	}

	lr := io.LimitReader(peer, msg.Size)
	fmt.Println("size of limit reader", msg.Size)

	fmt.Println("starting")
	// var decryptedFile bytes.Buffer
	// if _, err := copyDecrypt(&decryptedFile, lr, f.EncKey); err != nil {
	// 	return fmt.Errorf("decryption failed: %w", err)
	// }

	// if err := verifySignature(decryptedFile.Bytes(), msg.Signature, peerPublicKey); err != nil {
	// 	log.Printf("Signature verification failed from peer %s: %v", from, err)
	// 	return fmt.Errorf("signature verification failed: %w", err)
	// }
	// fmt.Printf("Signature verified from peer %s\n", from)
	// panic("not implemented")
	time.Sleep(1 * time.Second)
	n, err := f.store.Write(msg.ID, msg.Key, lr)
	time.Sleep(1 * time.Second)
	// panic("not implemented")
	fmt.Printf("[%s] Write completed.\n", f.Transort.ListenAddr())
	if err != nil {
		return err
	}

	// panic("not implemented")
	fmt.Printf("[%s] Written (%d) bytes to disk\n", f.Transort.ListenAddr(), n)

	_, reader, err, _ := f.store.Read(msg.ID, msg.Key)
	if err != nil {
		return err
	}

	//******while deleting files i was getting error, file is used by other process
	//so i need to close the reader to free up the file
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	// fmt.Println("Read")

	// Read all encrypted data into buffer for signature verification
	var encryptedData bytes.Buffer
	if _, err := io.Copy(&encryptedData, reader); err != nil {
		return fmt.Errorf("failed to read encrypted data from disk: %w", err)
	}
	// fmt.Println("Encrypted data length", len(encryptedData.Bytes()))
	// fmt.Printf("Encrypted data hash: %x\n", sha256.Sum256(encryptedData.Bytes()))
	// var decryptedData bytes.Buffer
	// if _, err := copyDecrypt(&decryptedData, reader, f.EncKey); err != nil {
	// 	return fmt.Errorf("decryption failed: %w", err)
	// }
	// fmt.Println("Decrypted data length", len(decryptedData.Bytes()))
	// fmt.Println("Decrypted data hash", sha256.Sum256(decryptedData.Bytes()))

	// Step 4: Get public key and verify
	peerPublicKey, ok := p2p.GetPeerPublicKey(from)
	if !ok {
		f.DeleteLocal(msg.Key)
		return fmt.Errorf("peer %s public key not found", from)
	}
	fmt.Println("Peer public key", peerPublicKey)
	if err := verifySignature(encryptedData.Bytes(), msg.Signature, peerPublicKey); err != nil {
		log.Printf("Signature verification failed from peer %s: %v", from, err)
		f.DeleteLocal(msg.Key)
		return fmt.Errorf("signature verification failed: %w", err)
	}

	f.store.SaveSignature(msg.Key, msg.Signature)

	fmt.Printf("Signature verified from peer %s\n", from)
	// panic("not implemented")
	peer.CloseStream()
	return nil
}

func (f *FileServer) handleMessageDeleteFile(from string, msg *MessageDeleteFile) error {
	// fmt.Printf("Received delete message %+v from %s\n", msg, from)
	_, ok := f.peers[from]
	if !ok {
		log.Printf("Rejecting delete from unknown peer %s", from)
		return fmt.Errorf("unauthorized delete request from unknown peer")
	}

	if f.store.Has(msg.ID, msg.Key) {
		fmt.Printf("deleteing [%s] file from peer [%s] \n", msg.Key, from)
		if err := f.store.Delete(msg.ID, msg.Key); err != nil {
			return err
		}
	}
	f.store.DeleteSignature(msg.Key)
	fmt.Println("deleted signature from peer", from)
	return nil
}

func (f *FileServer) handleMessageGetFileNotFound() error {
	// use a sync.Map or channel to track NACKs — see below for setup
	f.notFoundChan <- struct{}{}
	return nil
}

//utility functions

func (f *FileServer) Stop() {
	close(f.quitch)
}

func (f *FileServer) OnPeer(p p2p.Peer) error {
	f.peerLock.Lock()
	defer f.peerLock.Unlock()

	f.peers[p.RemoteAddr().String()] = p
	log.Println("New peer connected", p.RemoteAddr().String())
	return nil

}

func (f *FileServer) broadcast(msg *Message) error {

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}

	for _, p := range f.peers {
		fmt.Println("Broadcasting message to peer", p.RemoteAddr().String())
		// p.Send([]byte{p2p.IncommingMessage})
		full := append([]byte{p2p.IncommingMessage}, buf.Bytes()...)
		// p.Send(full)

		if err := p.Send(full); err != nil {
			return err
		}
	}
	return nil
}

//What it does:
// Collects all peer connections in a slice
// Wraps them using io.MultiWriter individually using ...
// Encodes the Payload once using gob
// Sends the encoded message to all peers
