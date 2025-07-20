package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
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

// type Pair struct {
// 	Key string
// 	Value string
// }

type Pair struct {
	hasFile bool
	addr    string
}

// FileServer represents a file server.
type FileServer struct {
	FileServerOps
	store                   *Store
	quitch                  chan struct{}
	peerLock                sync.Mutex
	peers                   map[string]p2p.Peer
	notFoundChan            chan struct{}
	auditLogger             *AuditLogger
	storeAckChan            chan MessageStoreAck
	deleteAckChan           chan MessageDeleteAck
	storedFiles             map[string]bool
	duplicationResponseChan chan Pair
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
		// ops.ID = "watashino"
	}

	// auditFileLocation := "audit/" + ops.RootStorage + "/" + ops.ID
	// logFileName := ops.ID + ".log"
	// audit, err := NewAuditLogger(auditFileLocation, logFileName)
	// if err != nil {
	// 	log.Fatalf("Audit logger initialization failed: %v", err)
	// }

	audit, err := simpleNewAuditLogger("audit.log")
	if err != nil {
		log.Fatalf("Audit logger initialization failed: %v", err)
	}

	return &FileServer{
		FileServerOps:           ops,
		store:                   NewStore(*storeOps),
		quitch:                  make(chan struct{}),
		peerLock:                sync.Mutex{},
		peers:                   make(map[string]p2p.Peer),
		notFoundChan:            make(chan struct{}, 100),
		auditLogger:             audit,
		storeAckChan:            make(chan MessageStoreAck, 100),
		deleteAckChan:           make(chan MessageDeleteAck, 100),
		storedFiles:             make(map[string]bool),
		duplicationResponseChan: make(chan Pair, 100),
	}, nil
}

func (f *FileServer) Store(key string, r io.Reader) error {

	//1. store the file to disk
	//2. broadcast the file to other peers in the network

	fmt.Printf("[%s] Starting Store operation for key: %s\n", f.Transort.ListenAddr(), key)
	f.auditLogger.Log("STORE", key, "START", "Initiating store operation")

	var (
		size         int64
		err          error
		fd           *os.File
		fileLocation string
		msg          Message
	)

	if ok := f.store.Has(f.ID, key); ok {
		f.auditLogger.Log("STORE", key, "LOCAL", "ALREADY_STORED")
		_, fd, err, fileLocation = f.store.Read(f.ID, key)
		if err != nil {
			f.auditLogger.Log("GET", key, "LOCAL", "FAIL_READ")
			return err
		}
		fmt.Printf("[%s] File %s already stored locally at location [%s]\n", f.Transort.ListenAddr(), key, fileLocation)
		//**** (IMPORTANT) add this to close hte file we write otherwise it will remain open and when performing delete operation it will throw error
		defer fd.Close()
		goto ENCRYPT_AND_STREAM
	}
	fmt.Printf("[%s] Storing file to disk\n", f.Transort.ListenAddr())
	size, err, fd = f.store.Write(f.ID, key, r)
	defer fd.Close()
	if err != nil {
		f.auditLogger.Log("STORE", key, "LOCAL", "FAILED")
		return err
	}
	fmt.Printf("[%s] File stored locally: %d bytes\n", f.Transort.ListenAddr(), size)
	f.auditLogger.Log("STORE", key, "LOCAL", "SUCCESS")

ENCRYPT_AND_STREAM:

	//check if the file is already stored in the network
	msg = Message{
		Payload: MessageDuplicateCheck{
			Key: hashKey(key) + getExtension(key),
			ID:  f.ID,
		},
	}
	if err := f.broadcast(&msg); err != nil {
		return err
	}

	// Wait 1 sec for all peers to reply
	timeout := time.After(3 * time.Second)
	peersWithNoFile := new([]string)

	waiting := len(f.peers)
	fmt.Println("Waiting for duplicate responses from", waiting, "peers")
END:
	for waiting > 0 {
		// fmt.Printf("[%s] Waiting for duplicate responses from %d peers\n", f.Transort.ListenAddr(), waiting)
		select {
		case pair := <-f.duplicationResponseChan:
			if !pair.hasFile {
				// fmt.Printf("[%s] Received duplicate response from peer %s for file %s\n", f.Transort.ListenAddr(), addr, msg.Payload.(MessageDuplicateCheck).Key)
				fmt.Printf("[%s ]  File is not present with peer [%s]\n", f.Transort.ListenAddr(), pair.addr)
				*peersWithNoFile = append(*peersWithNoFile, pair.addr)
			} else {
				fmt.Printf("[%s] File is present with peer [%s]\n", f.Transort.ListenAddr(), pair.addr)
			}
			waiting--
		case <-timeout:
			fmt.Printf("[%s] Timeout while waiting for duplicate responses\n", f.Transort.ListenAddr())
			goto END
		}
	}
	fmt.Println("peersWithNoFile after ack", len(*peersWithNoFile))
	if len(*peersWithNoFile) == 0 {
		fmt.Printf("[%s] File %s stored on all peers\n", f.Transort.ListenAddr(), key)
		f.auditLogger.Log("STORE", key, "LOCAL", "STORED_ON_ALL_PEERS")
		return nil
	} else if len(*peersWithNoFile) > 0 {
		fmt.Printf("[%s] File %s not stored on %d peers\n", f.Transort.ListenAddr(), key, len(*peersWithNoFile))
		f.auditLogger.Log("STORE", key, "LOCAL", "NOT_STORED_ON_SOME_PEERS")
	} else {
		fmt.Printf("[%s] File %s not stored on any peer\n", f.Transort.ListenAddr(), key)
		f.auditLogger.Log("STORE", key, "LOCAL", "NOT_STORED_ON_ANY_PEER")
	}

	tempFile, err := os.CreateTemp("", "enc_temp_*.bin")
	if err != nil {
		fmt.Printf("[%s] Failed to create temp file locally to store encrypted file\n", f.Transort.ListenAddr())
		f.auditLogger.Log("STORE", key, "LOCAL", "FAIL_CREATE_TEMP_FILE")
		return err
	}
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()

	privKey, err := p2p.LoadPrivateKey()
	if err != nil {
		f.auditLogger.Log("STORE", key, "KEYLOAD", "FAILED")
		return err
	}

	var (
		// fileReader   = bytes.NewReader(fileBuffer.Bytes())
		// encryptedBuf bytes.Buffer
		hasher      = sha256.New()
		multiWriter = io.MultiWriter(tempFile, hasher)
	)

	// fmt.Printf("Original file size before encryption: %d bytes\n", len(fileBuffer.Bytes()))
	// fmt.Println("Original file hash:", sha256.Sum256(fileBuffer.Bytes()))
	fmt.Printf("[%s] Encrypting file...\n", f.Transort.ListenAddr())
	f.auditLogger.Log("STORE", key, "ENCRYPT", "START")

	fd.Seek(0, io.SeekStart) // Rewind before re-reading
	_, err = copyEncrypt(multiWriter, fd, f.EncKey)
	if err != nil {
		f.auditLogger.Log("STORE", key, "ENCRYPT", "FAILED")
		return err
	}
	f.auditLogger.Log("STORE", key, "ENCRYPT", "SUCCESS")
	// fmt.Printf("Size returned by copyEncrypt: %d bytes\n", n)
	// fmt.Printf("Actual teeBuf size: %d bytes\n", len(teeBuf.Bytes()))
	// fmt.Printf("Encrypted data hash (IV + encrypted): %x\n", sha256.Sum256(teeBuf.Bytes()))

	digest := hasher.Sum(nil)
	signature, err := signSignature(digest[:], privKey) // fileBuffer already contains full file
	if err != nil {
		f.auditLogger.Log("STORE", key, "SIGN", "FAILED")
		return err
	}

	// fmt.Printf("[%s] Signature generated: %x\n", f.Transort.ListenAddr(), signature)

	//store the signature of encrypted file in the map
	sigKey := hashKey(key) + getExtension(key)
	f.store.SaveSignature(sigKey, signature)

	fmt.Printf("[%s] Signature generated and saved\n", f.Transort.ListenAddr())
	f.auditLogger.Log("STORE", key, "SIGN", "SUCCESS")
	// fmt.Printf("YOYOYO [%s] Signature stored in map for key [%s]", f.Transort.ListenAddr(), sigKey)
	fStat, err := os.Stat(tempFile.Name())
	if err != nil {
		return err
	}
	msg = Message{
		Payload: MessageStoreFile{
			Key:  sigKey,
			Size: int64(fStat.Size()),
			ID:   f.ID,
			// Ext:  getExtension(key),
			Signature: signature,
		},
	}

	fmt.Printf("[%s] Broadcasting metadata to peers...\n", f.Transort.ListenAddr())
	f.auditLogger.Log("STORE", key, "BROADCAST", "START")

	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	fmt.Println("peersWithNoFile", len(*peersWithNoFile))

	if len(*peersWithNoFile) == 0 {
		fmt.Println("Complete Broadcast")
		if err := f.broadcast(&msg); err != nil {
			f.auditLogger.Log("STORE", key, "BROADCAST", "FAILED")
			return err
		}

	} else {
		fmt.Println("Limited broadcast")
		if err := f.broadcastLimited(*peersWithNoFile, &msg); err != nil {
			f.auditLogger.Log("STORE", key, "BROADCAST", "FAILED")
			return err
		}

	}

	time.Sleep(50 * time.Millisecond)

	peers := []io.Writer{}
	if len(*peersWithNoFile) > 0 {
		for _, peer := range *peersWithNoFile {
			peers = append(peers, f.peers[peer])
		}
	} else {
		for _, peer := range f.peers {
			peers = append(peers, peer)
		}
	}

	mw := io.MultiWriter(peers...)

	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		f.auditLogger.Log("STORE", key, "ENCRYPT", "FAILED")

		return err
	}
	
	var totalPeers int
	if len(*peersWithNoFile) > 0 {
		totalPeers = len(*peersWithNoFile)
	} else {
		totalPeers = len(f.peers)
	}

	fmt.Printf("[%s] Streaming encrypted file to %d peers...\n", f.Transort.ListenAddr(), len(f.peers)-len(*peersWithNoFile))
	f.auditLogger.Log("STORE", key, "STREAM", "START")
	mw.Write([]byte{p2p.IncommingStream})
	// Create a fresh reader from the original file data for streaming
	// Stream the encrypted data (IV + encrypted data) to peers
	// encryptedReader := bytes.NewReader(tempFile)

	nw, err := io.Copy(mw, tempFile)
	if err != nil {
		f.auditLogger.Log("STORE", key, "STREAM", "FAILED")
		return err
	}
	fmt.Printf("[%s] Encrypted file streamed to peers: %d bytes\n", f.Transort.ListenAddr(), nw)

	f.auditLogger.Log("STORE", key, "STREAM", "SUCCESS")
	fmt.Printf("[%s] Recieved and written (%d) bytes to disk\n", f.Transort.ListenAddr(), nw)

	time.Sleep(500 * time.Millisecond)

	successCount := 0
	failures := make(map[string]string)

	var expectedAcks int
	if len(*peersWithNoFile) > 0 {
		expectedAcks = len(*peersWithNoFile)
	} else {
		expectedAcks = len(f.peers)
	}

	timeout = time.After(2 * time.Second)
	timeoutErr := ""

	// ackLoop:
	for expectedAcks > 0 {
		select {
		case ack := <-f.storeAckChan:
			fmt.Printf("[%s] Received store ACK from peer %s\n", f.Transort.ListenAddr(), ack.From)
			// fmt.Printf("[%s] Received error %s\n", f.Transort.ListenAddr(), ack.Err)
			expectedAcks--
			if ack.Err == "" {
				successCount++
				// fmt.Println("successCount", successCount)
			} else {
				failures[ack.From] = ack.Err
			}
		case <-timeout:
			timeoutErr = fmt.Sprintf("[%s] Timeout waiting for store ACKs from all peers.\n", f.Transort.ListenAddr())
			f.auditLogger.Log("REPLICATE", key, "ALL", "ACK_TIMEOUT")
			return errors.New(timeoutErr)
			// break ackLoop
		}
	}
	if successCount == totalPeers {
		fmt.Printf("[%s] File %s successfully replicated to %d peers.\n", f.Transort.ListenAddr(), key, successCount)
		f.auditLogger.Log("REPLICATE", key, "NETWORK", "SUCCESS")
	} else {
		fmt.Printf("[%s] Replication failed on %d peers:\n", f.Transort.ListenAddr(), totalPeers-successCount)
		for peer, reason := range failures {
			fmt.Printf(" - %s: %s\n", peer, reason)
			f.auditLogger.Log("REPLICATE", key, peer, "FAIL: "+reason)
		}
	}

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
		f.auditLogger.Log("GET", key, "LOCAL", "SUCCESS")
		_, r, err, fileLocation := f.store.Read(f.ID, key)
		if err != nil {
			f.auditLogger.Log("GET", key, "LOCAL", "FAIL_READ")
			return nil, "", err
		}
		// fmt.Printf("[%s] File location %s\n", f.Transort.ListenAddr(), fileLocation)

		// if rc, ok := r.(io.ReadCloser); ok {
		defer r.Close() // ✅ this ensures file is closed after use
		// }

		return r, fileLocation, nil
	}
	fmt.Printf("[%s] dont have file [%s], fetching from network...\n", f.Transort.ListenAddr(), key)
	f.auditLogger.Log("GET", key, "LOCAL", "MISS")

	msg := Message{
		Payload: MessageGetFile{
			Key: hashKey(key) + getExtension(key),
			ID:  f.ID,
		},
	}
	f.auditLogger.Log("GET", key, "NETWORK", "START_BROADCAST")
	err := f.broadcast(&msg)
	if err != nil {
		f.auditLogger.Log("GET", key, "NETWORK", "FAIL_BROADCAST")
		return nil, "", err
	}

	time.Sleep(500 * time.Millisecond)
	noOfPeers := len(f.peers)
	notFound := 0 //tracks peers that explicitly said “I don’t have it.”
	received := 0 //counts any kind of peer response — stream or not-found

	//why recieved? so that we can count all responses and if responses == no of peers
	//it means we were unable to fetch the file from any peer

	for {
		fmt.Printf("[%s] Waiting for peer response...\n", f.Transort.ListenAddr())

		select {
		case <-f.notFoundChan:
			fmt.Printf("[%s] A peer responded: file [%s] not found\n", f.Transort.ListenAddr(), key)

			notFound++
			received++ // we will count response even if stream fails later
			if notFound == noOfPeers {
				f.auditLogger.Log("GET", key, "NETWORK", "FAIL_NOT_FOUND_ON_ANY_PEER")
				return nil, "", fmt.Errorf("[%s] and all its peers don't have file [%s]", f.Transort.ListenAddr(), key)
			}

		case streamPeer := <-f.Transort.ConsumeStream():
			fmt.Printf("[%s] Received encrypted file stream from peer %s\n", f.Transort.ListenAddr(), streamPeer.RemoteAddr().String())
			received++
			// Only process if this is the same peer
			// if streamPeer.RemoteAddr().String() != peer.RemoteAddr().String() {
			// 	fmt.Printf("Skipping peer %s, stream came from %s\n", peer.RemoteAddr().String(), streamPeer.RemoteAddr().String())
			// 	continue
			// }

			//we will stream the encrypted file to a temp file
			tempEncFile, err := os.CreateTemp("", "enc_tmp_*.bin")
			if err != nil {
				fmt.Printf("[%s] Failed to create temp file, skipping...\n", f.Transort.ListenAddr())
				f.auditLogger.Log("GET", key, "NETWORK", "FAIL_CREATE_TEMP_FILE")
				continue
			}
			defer func() {
				tempEncFile.Close()
				os.Remove(tempEncFile.Name())
			}()

			var (
				fileSize int64
				sigLen   int64
			)
			peerAddr := streamPeer.RemoteAddr().String()

			if err := binary.Read(streamPeer, binary.LittleEndian, &fileSize); err != nil {
				fmt.Printf("[%s] Failed to read fileSize from peer %s, skipping...\n", f.Transort.ListenAddr(), peerAddr)
				f.auditLogger.Log("GET", key, peerAddr, "FAIL_READ_FILE_SIZE")
				continue
			}

			if err := binary.Read(streamPeer, binary.LittleEndian, &sigLen); err != nil {
				fmt.Printf("[%s] Failed to read sigLen from peer %s, skipping...\n", f.Transort.ListenAddr(), peerAddr)
				f.auditLogger.Log("GET", key, peerAddr, "FAIL_READ_SIG_LEN")
				continue
			}

			signature := make([]byte, sigLen)
			if _, err := io.ReadFull(streamPeer, signature); err != nil {
				fmt.Printf("[%s] Failed to read signature from peer %s, skipping...\n", f.Transort.ListenAddr(), peerAddr)
				f.auditLogger.Log("GET", key, peerAddr, "FAIL_READ_SIGNATURE")
				continue
			}
			// fmt.Printf(" [%s] -> [%s]", msg.Payload.(MessageGetFile).Key, signature)

			// encData := make([]byte, fileSize)
			encData := io.LimitReader(streamPeer, fileSize)
			hasher := sha256.New()
			tee := io.TeeReader(encData, io.MultiWriter(tempEncFile, hasher))

			// var encryptedBuf bytes.Buffer
			if _, err := io.Copy(io.Discard, tee); err != nil {
				fmt.Printf("[%s] Failed to copy encrypted data from peer %s, skipping...\n", f.Transort.ListenAddr(), peerAddr)
				f.auditLogger.Log("GET", key, peerAddr, "FAIL_READ_ENCRYPTED_DATA")
				//remove temp file
				os.Remove(tempEncFile.Name())
				continue
			}

			peerPublicKey, ok := p2p.GetPeerPublicKey(peerAddr)
			if !ok {
				fmt.Printf("[%s] Peer %s public key not found, skipping...\n", f.Transort.ListenAddr(), peerAddr)
				f.auditLogger.Log("GET", key, peerAddr, "FAIL_PUBLIC_KEY_NOT_FOUND")
				continue
			}

			digest := hasher.Sum(nil)
			if err := verifySignature(digest, signature, peerPublicKey); err != nil {
				fmt.Printf("[%s] Signature verification FAILED from peer %s\n", f.Transort.ListenAddr(), peerAddr)
				f.auditLogger.Log("GET", key, peerAddr, "FAIL_SIGNATURE_VERIFICATION")
				os.Remove(tempEncFile.Name())
				continue
			}
			fmt.Printf("[%s] Signature verification PASSED from peer %s\n", f.Transort.ListenAddr(), peerAddr)
			f.auditLogger.Log("GET", key, peerAddr, "SIGNATURE_VERIFICATION_SUCCESS")

			// encryptedReader := bytes.NewReader(encryptedBuf.Bytes())

			// Step 1: Seek back to beginning of temp file
			if _, err := tempEncFile.Seek(0, io.SeekStart); err != nil {
				fmt.Printf("[%s] Failed to rewind temp encrypted file: %v\n", f.Transort.ListenAddr(), err)
				f.auditLogger.Log("GET", key, peerAddr, "FAIL_REWIND_TEMP_FILE")
				os.Remove(tempEncFile.Name())
				continue
			}

			// Step 2: Use it as reader for decryption
			n, err := f.store.WriteDecrypted(f.ID, key, f.EncKey, tempEncFile)
			if err != nil {
				fmt.Printf("[%s] Failed to write decrypted file to disk: %v\n", f.Transort.ListenAddr(), err)
				f.auditLogger.Log("GET", key, peerAddr, "FAIL_WRITE_DECRYPTED_FILE")
				os.Remove(tempEncFile.Name())
				continue
			}
			fmt.Printf("[%s] Wrote %d decrypted bytes to disk from peer %s\n", f.Transort.ListenAddr(), n, peerAddr)
			streamPeer.CloseStream()

			//now we have the file on disk, if we are unable to read it we will simply return error
			_, r, err, fileLocation := f.store.Read(f.ID, key)
			if err != nil {
				f.auditLogger.Log("GET", key, peerAddr, "FAIL_READ_AFTER_WRITE")
				return nil, "", err
			}
			f.auditLogger.Log("GET", key, peerAddr, "SUCCESS")

			// fmt.Printf("[%s] File location %s\n", f.Transort.ListenAddr(), fileLocation)
			return r, fileLocation, nil

		case <-time.After(2 * time.Second):
			// Timeout waiting for this peer — skip
			f.auditLogger.Log("GET", key, "NETWORK", "FAIL_TIMEOUT_ALL_PEERS")

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

	fmt.Printf("[%s] Starting Delete operation for key: %s\n", f.Transort.ListenAddr(), key)
	f.auditLogger.Log("DELETE", key, "START", "Initiating delete operation")

	//first delete the file locally if its present
	f.DeleteLocal(key)

	//now delete the file from all the peers in the network
	fmt.Printf("[%s] Initiating delete broadcast for file [%s] to peers...\n", f.Transort.ListenAddr(), key)
	f.auditLogger.Log("DELETE", key, "NETWORK", "START_BROADCAST")
	msg := Message{
		Payload: MessageDeleteFile{
			Key: hashKey(key) + getExtension(key),
			ID:  f.ID,
		},
	}

	err := f.broadcast(&msg)
	if err != nil {
		f.auditLogger.Log("DELETE", key, "NETWORK", "FAIL_BROADCAST")
		return err
	}

	expectedAcks, totalPeers := len(f.peers), len(f.peers)
	successCount := 0
	failures := make(map[string]string)
	timeout := time.After(5 * time.Second)
	// timeoutErr := ""
	// ackLoop:
	for expectedAcks > 0 {
		select {
		case ack := <-f.deleteAckChan:
			fmt.Printf("[%s] Received delete ACK from peer %s\n", f.Transort.ListenAddr(), ack.From)
			// fmt.Printf("[%s] Received error %s\n", f.Transort.ListenAddr(), ack.Err)
			expectedAcks--
			if ack.Err == "" {
				successCount++
				fmt.Println("successCount", successCount)
			} else {
				failures[ack.From] = ack.Err
			}
		case <-timeout:
			f.auditLogger.Log("DELETE", key, "ALL", "ACK_TIMEOUT")
			fmt.Printf("[%s] Timeout waiting for delete ACKs from all peers.\n", f.Transort.ListenAddr())
			return errors.New("timeout waiting for delete ACKs from all peers")
		}
	}
	if successCount == totalPeers {
		fmt.Printf("[%s] File %s successfully deleted from %d peers.\n", f.Transort.ListenAddr(), key, successCount)
		f.auditLogger.Log("DELETE", key, "NETWORK", "SUCCESS")
	} else {
		fmt.Printf("[%s] Delete failed on %d peers:\n", f.Transort.ListenAddr(), totalPeers-successCount)
		for peer, reason := range failures {
			fmt.Printf(" - %s: %s\n", peer, reason)
			f.auditLogger.Log("DELETE", key, peer, "FAIL: "+reason)
		}
	}

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
		f.auditLogger.Log("DELETE", key, "LOCAL", "SUCCESS")
		fmt.Printf("[%s] deleted signature from local disk\n", f.Transort.ListenAddr())
	} else {
		fmt.Printf("[%s] File [%s] not found locally. Skipping local deletion.\n", f.Transort.ListenAddr(), key)
		f.auditLogger.Log("DELETE", key, "LOCAL", "NOT_FOUND")
	}

	return nil
}

func (f *FileServer) DeleteRemote(key string) error {
	ext := getExtension(key)
	remotekey := hashKey(key)
	if len(ext) != 0 {
		remotekey += ext
	}
	fmt.Printf("[%s] file [%s] hash key [%s]\n", f.Transort.ListenAddr(), key, remotekey)
	if f.store.Has(f.ID, remotekey) {
		fmt.Printf("[%s] have file [%s], deleting from local disk\n", f.Transort.ListenAddr(), key)
		err := f.store.Delete(f.ID, remotekey)
		if err != nil {
			return err
		}
		f.store.DeleteSignature(hashKey(key))
		f.auditLogger.Log("DELETE", key, "LOCAL", "SUCCESS")
		fmt.Printf("[%s] deleted signature from local disk\n", f.Transort.ListenAddr())
	} else {
		fmt.Printf("[%s] File [%s] not found locally. Skipping local deletion.\n", f.Transort.ListenAddr(), key)
		f.auditLogger.Log("DELETE", key, "LOCAL", "NOT_FOUND")
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
	case MessageStoreAck:
		return f.handleMessageStoreAck(from, &v)
	case MessageDeleteAck:
		return f.handleMessageDeleteAck(from, &v)
	case MessageDuplicateCheck:
		return f.handleMessageDuplicateCheck(from, &v)
	case MessageDuplicateResponse:
		return f.handleMessageDuplicateResponse(from, &v)

	}
	return nil
}

func (f *FileServer) handleMessageDuplicateResponse(from string, msg *MessageDuplicateResponse) error {
	// fmt.Println("executing handleMessageDuplicateResponse")
	fmt.Printf("[%s] Received duplicate response from peer %s for file %s\n", f.Transort.ListenAddr(), from, msg.Key)
	fmt.Println("HasIt", msg.HasIt)
	if msg.HasIt {
		f.duplicationResponseChan <- Pair{hasFile: true, addr: from}
		fmt.Printf("DEBUG_FOUND: Successfully sent to duplicationResponseChan: %s\n", from)

	} else {
		f.duplicationResponseChan <- Pair{hasFile: false, addr: from}
		fmt.Printf("DEBUG_NOT_FOUND: Successfully sent to duplicationResponseChan: %s\n", from)
	}
	return nil
}

func (f *FileServer) handleMessageDuplicateCheck(from string, msg *MessageDuplicateCheck) error {
	// fmt.Println("executing handleMessageDuplicateCheck")
	fmt.Printf("[%s] Checking for duplicate file %s on peer %s\n", f.Transort.ListenAddr(), msg.Key, from)
	var HasFile bool = false
	if f.store.Has(msg.ID, msg.Key) {
		HasFile = true
		fmt.Printf("[%s] File %s already exists on peer %s\n", f.Transort.ListenAddr(), msg.Key, from)
	}

	response := Message{
		Payload: MessageDuplicateResponse{
			Key:   msg.Key,
			ID:    msg.ID,
			HasIt: HasFile,
			From:  f.Transort.ListenAddr(),
		},
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(response); err != nil {
		fmt.Println("encoding error", err)
		return err
	}

	// send back as direct message (not stream)
	p := f.peers[from]
	p.Send([]byte{p2p.IncommingMessage})
	return p.Send(buf.Bytes())
}

func (f *FileServer) handleMessageDeleteAck(from string, msg *MessageDeleteAck) error {
	// fmt.Println("executing handleMessageDeleteAck")
	fmt.Printf("Received DELETE ACK from %s for file %s\n", from, msg.Key)
	f.deleteAckChan <- *msg
	return nil
}

func (f *FileServer) handleMessageStoreAck(from string, msg *MessageStoreAck) error {
	// fmt.Println("executing handleMessageStoreAck")
	fmt.Printf("Received STORE ACK from %s for file %s\n", from, msg.Key)
	f.storeAckChan <- *msg
	return nil
}

func (f *FileServer) handleMessageGetFile(from string, msg *MessageGetFile) error {
	// fmt.Println("executing handleMessageGetFile")

	peer, ok := f.peers[from]
	// fmt.Println("YOOOOOOOOOOOOOOOOOO-peer", peer.LocalAddr())
	// fmt.Println("YOOOOOOOOOOOOOOOOOO-peer", peer.RemoteAddr())
	// fmt.Println("YOOOOOOOOOOOOOOOOOO-from", from)
	if !ok {
		f.auditLogger.Log("GET", msg.Key, from, "FAIL_PEER_NOT_FOUND")
		return fmt.Errorf("peer %s could not be found in map", from)
	}

	if !f.store.Has(msg.ID, msg.Key) {
		fmt.Printf("[%s] need to serve file (%s) to peer %s but it does not exist on the network", f.Transort.ListenAddr(), msg.Key, from)
		f.auditLogger.Log("GET", msg.Key, from, "FILE_NOT_FOUND")
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
			f.auditLogger.Log("GET", msg.Key, from, "FAIL_ENCODING_NACK")
			return err
		}
		peer.Send([]byte{p2p.IncommingMessage})
		if err := peer.Send(buf.Bytes()); err != nil {
			f.auditLogger.Log("GET", msg.Key, from, "FAIL_SENDING_NACK")
			return err
		}
		return fmt.Errorf("[%s] need to serve file (%s) to peer %s but it does not exist on the network", f.Transort.ListenAddr(), msg.Key, from)

	}

	fmt.Printf("[%s] Serving file [%s] to peer [%s]\n", f.Transort.ListenAddr(), msg.Key, from)
	f.auditLogger.Log("GET_REQ", msg.Key, from, "FOUND")

	fileSize, rd, err, _ := f.store.Read(msg.ID, msg.Key)
	if err != nil {
		f.auditLogger.Log("GET", msg.Key, from, "FAIL_READING_FILE")
		return err
	}

	// r, ok := rd.(io.ReadCloser)
	// if ok {
	// 	// fmt.Println("Closing reader")
	defer rd.Close()
	// }

	signature, err := f.store.GetSignature(msg.Key)
	if err != nil {
		f.auditLogger.Log("GET", msg.Key, from, "FAIL_GETTING_SIGNATURE")
		return fmt.Errorf("signature for file %s could not be found in map on peer [%s]", msg.Key, from)
	}
	fmt.Printf("[%s] Found signature for [%s], sending metadata to peer [%s]\n", f.Transort.ListenAddr(), msg.Key, from)
	// First send the "incomingStream" byte to the peer and then we can send
	// the file size as an int64.

	// Send metadata: [Stream Signal] [file size] [sig length] [signature]
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
		f.auditLogger.Log("GET_REQ", msg.Key, from, "FAIL_STREAM")
		return err
	}
	// f.incomingStreamChan <- peer
	fmt.Printf("[%s] Sent (%d bytes) of file [%s] to peer [%s]\n", f.Transort.ListenAddr(), n, msg.Key, from)
	f.auditLogger.Log("GET_REQ", msg.Key, from, "SUCCESS")
	return nil

}

func (f *FileServer) handleMessageStoreFile(from string, msg *MessageStoreFile) error {
	// panic("not implemented")
	// fmt.Printf("Received data message %+v from %s\n", msg, from)
	var (
		err error
		n   int64
		fd  *os.File
		// reader io.Reader
	)
	err = nil //default error
	peer, ok := f.peers[from]
	if !ok {
		err = fmt.Errorf("[%s] Peer %s not found in map", f.Transort.ListenAddr(), from)
		f.auditLogger.Log("STORE", msg.Key, from, "FAIL_PEER_NOT_FOUND")
		return err
	}

	fmt.Printf("[%s] Receiving encrypted file %s from peer %s (%d bytes)...\n", f.Transort.ListenAddr(), msg.Key, from, msg.Size)

	hasher := sha256.New()
	lr := io.TeeReader(io.LimitReader(peer, msg.Size), hasher)

	n, err, fd = f.store.Write(msg.ID, msg.Key, lr)

	//**** (IMPORTANT) add this to close hte file we write otherwise it will remain open and when performing delete operation it will throw error
	if fd != nil {
		defer fd.Close()
	}

	// panic("not implemented")
	if err != nil {
		f.auditLogger.Log("STORE", msg.Key, from, "FAIL_WRITE_DISK")
		return err
	}

	// panic("not implemented")
	fmt.Printf("[%s] Written (%d) bytes to disk on peer %s\n", f.Transort.ListenAddr(), n, from)

	// Step 4: Get public key and verify
	peerPublicKey, ok := p2p.GetPeerPublicKey(from)
	if !ok {
		err = fmt.Errorf("[%s] Public key not found for peer [%s] Deleting file", f.Transort.ListenAddr(), from)
		fmt.Printf("[%s] Public key not found for peer [%s]. Deleting file.\n", f.Transort.ListenAddr(), from)

		f.DeleteLocal(msg.Key)
		f.auditLogger.Log("STORE", msg.Key, from, "FAIL_PUBLIC_KEY_NOT_FOUND")
		return err
	}
	// fmt.Println("Peer public key", peerPublicKey)
	// fmt.Printf("[%s] Signature : %x\n", f.Transort.ListenAddr(), msg.Signature)

	digest := hasher.Sum(nil)
	if err = verifySignature(digest, msg.Signature, peerPublicKey); err != nil {
		fmt.Printf("[%s] Signature verification FAILED for file %s from peer %s: %v\n", f.Transort.ListenAddr(), msg.Key, from, err)
		f.DeleteLocal(msg.Key)
		f.auditLogger.Log("STORE", msg.Key, from, "FAIL_INVALID_SIGNATURE")
		return err
	}

	f.store.SaveSignature(msg.Key, msg.Signature)

	fmt.Printf("[%s] Signature VERIFIED and file %s from peer %s stored successfully.\n", f.Transort.ListenAddr(), msg.Key, from)
	f.auditLogger.Log("STORE", msg.Key, from, "SUCCESS")

	peer.CloseStream()

	var errStr string
	if err != nil {
		errStr = err.Error()
	} else {
		errStr = ""
	}
	//sending acknowledge to the sender
	ack := Message{
		Payload: MessageStoreAck{
			Key:  msg.Key,
			From: f.Transort.ListenAddr(),
			Err:  errStr,
		},
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(ack); err != nil {
		fmt.Printf("Failed to encode ACK to peer %s: %v\n", from, err)
		return nil // we return nil because file operation succeeded
	}

	if err := peer.Send([]byte{p2p.IncommingMessage}); err != nil {
		fmt.Printf("Failed to send ACK prefix to peer %s: %v\n", from, err)
		return nil
	}
	if err := peer.Send(buf.Bytes()); err != nil {
		fmt.Printf("Failed to send ACK to peer %s: %v\n", from, err)
		return nil
	}

	return nil
}

func (f *FileServer) handleMessageDeleteFile(from string, msg *MessageDeleteFile) error {
	// fmt.Printf("Received delete message %+v from %s\n", msg, from)
	err := error(nil)
	peer, ok := f.peers[from]
	if !ok {
		f.auditLogger.Log("DELETE", msg.Key, from, "FAIL_UNKNOWN_PEER")
		log.Printf("[%s] Rejecting delete from unknown peer %s", f.Transort.ListenAddr(), from)
		err = fmt.Errorf("unauthorized delete request from unknown peer")
		return err
	}

	if f.store.Has(msg.ID, msg.Key) {
		fmt.Printf("[%s] deleteing [%s] file from peer [%s] \n", f.Transort.ListenAddr(), msg.Key, from)
		if err = f.store.Delete(msg.ID, msg.Key); err != nil {
			f.auditLogger.Log("DELETE", msg.Key, from, "FAIL_DELETE")
		}
	}
	f.auditLogger.Log("DELETE", msg.Key, from, "SUCCESS")
	f.store.DeleteSignature(msg.Key)
	fmt.Printf("[%s] deleted signature from peer [%s] \n", f.Transort.ListenAddr(), from)
	f.auditLogger.Log("DELETE_SIGNATURE", msg.Key, from, "SUCCESS")

	//sending acknowledge to the sender

	var errStr string
	if err != nil {
		errStr = err.Error()
	} else {
		errStr = ""
	}

	ack := Message{
		Payload: MessageDeleteAck{
			Key:  msg.Key,
			From: f.Transort.ListenAddr(),
			Err:  errStr,
		},
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(ack); err != nil {
		fmt.Printf("Failed to encode ACK to peer %s: %v\n", from, err)
		return nil // we return nil because file operation succeeded
	}

	if err := peer.Send([]byte{p2p.IncommingMessage}); err != nil {
		fmt.Printf("Failed to send ACK prefix to peer %s: %v\n", from, err)
		return nil
	}
	if err := peer.Send(buf.Bytes()); err != nil {
		fmt.Printf("Failed to send ACK to peer %s: %v\n", from, err)
		return nil
	}

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
		fmt.Printf("[%s] Broadcasting message to peer %s\n", f.Transort.ListenAddr(), p.RemoteAddr().String())
		// p.Send([]byte{p2p.IncommingMessage})
		full := append([]byte{p2p.IncommingMessage}, buf.Bytes()...)
		// p.Send(full)

		if err := p.Send(full); err != nil {
			return err
		}
	}
	return nil
}

func (f *FileServer) broadcastDuplicateCheck(msg *Message) error {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}

	for _, p := range f.peers {
		fmt.Printf("[%s] Broadcasting duplicate check to peer %s\n", f.Transort.ListenAddr(), p.RemoteAddr().String())
		p.Send([]byte{p2p.IncommingMessage})
		// full := append([]byte{p2p.IncommingMessage}, buf.Bytes()...)
		// p.Send(full)

		if err := p.Send(buf.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

func (f *FileServer) broadcastLimited(peers []string, msg *Message) error {
	// fmt.Println("Limited broadcast")
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}

	for _, addr := range peers {
		fmt.Printf("[%s] Broadcasting message to peer %s\n", f.Transort.ListenAddr(), addr)
		// p.Send([]byte{p2p.IncommingMessage})
		full := append([]byte{p2p.IncommingMessage}, buf.Bytes()...)
		// p.Send(full)
		peer, ok := f.peers[addr]
		if !ok {
			return fmt.Errorf("peer %s not found", addr)
		}
		if err := peer.Send(full); err != nil {
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
