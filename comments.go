package main

//Get(key string)

//Get attempts to retrieve a file identified by the key.
// If the file is available locally, it is returned from disk.
// Otherwise, a broadcast request is sent to all peers in the network.
// The function waits for responses from peers:
//   - If a peer responds with the file stream, it writes the file to disk and returns it.
//   - If a peer responds with a "not found" message, it is counted toward the failure threshold.
// The function continues listening until:
//   - The file is successfully retrieved from any peer, or
//   - All peers either fail or report the file as not found.
// In both failure scenarios, an appropriate error is returned.


	// var decryptedBuf bytes.Buffer

	// if _, err := copyDecrypt(&decryptedBuf, lr, f.EncKey); err != nil {
	// 	return err
	// }
	// fmt.Println("decrypted")
	// hash, _, err := hashFileContent(&decryptedBuf)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println("hashed")
	// if ok := f.Transort.CheckFileHashMap(hash); ok {
	// 	return fmt.Errorf("[%s] file %s already exists on the peer [%s]", f.Transort.ListenAddr(), hash, from)
	// }
	// fmt.Println("checked")
	// f.Transort.AddFileHashMap(hash, msg.Key)
	// fmt.Println("added")
	// n, err := f.store.Write(msg.ID, msg.Key, bytes.NewReader(decryptedBuf.Bytes()))


	// fileHash, size, err := hashFileContent(fileBuffer)
	// if err != nil {
	// 	return err
	// }

	// if keyExist, ok := f.store.HashMap[fileHash]; !ok {
	// 	f.store.HashMap[fileHash] = key

	// 	fmt.Println("Storing file to disk")
	// 	_, err := f.store.Write(f.ID, key, tee)
	// 	if err != nil {
	// 		return err
	// 	}

	// } else {
	// 	fmt.Println("File already exists with key", keyExist)
	// }