package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// defaultRoot is the default root folder used if no Root is provided
const defaultRoot = "gyattt"

// PathKey represents the file path and file name of a file.
type PathKey struct {
	FilePath string
	FileName string
}

type PathTransformFunc func(string) PathKey

// DefaultPathTransformFunc is used when no PathTransformFunc is provided.
var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		FilePath: key,
		FileName: key,
	}
}

// CryptoPathTransformFunc transforms the given path into a PathKey containing path where the file
// will be stored and filename using SHA-1 hash.
func CryptoPathTransformFunc(path string) PathKey {
	hasher := sha1.New()
	hasher.Write([]byte(path))
	h := hasher.Sum(nil)
	hash := hex.EncodeToString(h)

	blockSize := 5
	sliceLen := len(hash) / blockSize
	newPath := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		newPath[i] = string(hash[from:to])
	}
	return PathKey{FilePath: strings.Join(newPath, "/"), FileName: hash}

}

// StoreOps represents the options for a store.
type StoreOps struct {
	//Root is the root folder containing all the folders/files of the system
	Root              string
	PathTransformFunc PathTransformFunc
	
}

// Store represents a store.
type Store struct {
	StoreOps
	//map to store hash of file -> key
	HashMap map[string]string
	signatureMap map[string][]byte
}

// NewStore creates a new store
func NewStore(ops StoreOps) *Store {
	if ops.PathTransformFunc == nil {
		ops.PathTransformFunc = DefaultPathTransformFunc
	}
	if len(ops.Root) == 0 {
		ops.Root = defaultRoot
	}
	
	return &Store{
		StoreOps: ops,
		HashMap: make(map[string]string),
		signatureMap: make(map[string][]byte),
	}
}
func(s *Store) SaveSignature(key string, sig []byte) {
	s.signatureMap[key] = sig
}
func(s *Store) GetSignature(key string) ([]byte, error) {
	sig, ok := s.signatureMap[key]
	if !ok {
		return nil, errors.New("signature not found")
	}
	return sig, nil
}

func(s *Store) DeleteSignature(key string) {
	delete(s.signatureMap, key)
}

// Checks whether the given key maps to an existing file on disk under the store’s directory.
func (s *Store) Has(id string, key string) bool {

	ext := getExtension(key)
	key = removeExtension(key)
	// fmt.Println("Key in has", key)

	pathkey := s.PathTransformFunc(key)
	// fmt.Println("PathKey in has", pathkey)
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathkey.PathAndFileName())
	if len(ext) != 0 {
		fullPathWithRoot += ext
	}
	// fmt.Println("Full path in has", fullPathWithRoot)

	_, err := os.Stat(fullPathWithRoot)
	return !errors.Is(err, os.ErrNotExist)
}

// Reads the file from disk and returns the reader
func (s *Store) Read(id string, key string) (int64, *os.File, error, string) {

	//we dont need to read the file in buffer as file can be very large and
	//it can take time copying it in buffer so we returned the reader directl

	return s.readStream(id, key)
}

// Writes the file to disk
func (s *Store) Write(id, key string, w io.Reader) (int64, error, *os.File) {
	// fmt.Println("executing write")
	// panic("not implemented")
	return s.writeStream(id, key, w)
}

// Deletes the file from disk
func (s *Store) Delete(id string, key string) error {

	// ext := getExtension(key)
	key = removeExtension(key)
	pathKey := s.PathTransformFunc(key)

	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.FileName)
	}()

	firstPathNameWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FirstPathName())
	// fmt.Println("Delete location", firstPathNameWithRoot)

	return os.RemoveAll(firstPathNameWithRoot)
}

// Clears the root directory
func (s *Store) ClearRoot() error {
	return os.RemoveAll(s.Root)
}

// Decrypts the file and write it to disk
func (s *Store) WriteDecrypted(id, key string, encKey []byte, w io.Reader) (int, error) {
	return s.writeDecryptedStream(id, key, encKey, w)
}

//utility functions

func removeExtension(key string) string {
	return strings.TrimSuffix(key, getExtension(key))

}
func (s *Store) openFileForWriting(id, key string) (*os.File, error) {
	// fmt.Println("executing writeStream")
	ext := getExtension(key)
	key = removeExtension(key)
	// fmt.Println("Key in write", key)
	//creating pathname using transform function
	pathKey := s.PathTransformFunc(key)
	// fmt.Println("PathKey in write", pathKey)

	FilePathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FilePath)

	// fmt.Println("FilePathWithRoot", FilePathWithRoot)
	//creating directory recursively with pathname
	err := os.MkdirAll(FilePathWithRoot, 0755)
	if err != nil {
		fmt.Println("Error creating directory", err)
		return nil, err
	}
	// fmt.Println("Directory created")

	//creating file name with pathname and filename and extension if available
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.PathAndFileName())
	if len(ext) != 0 {
		fullPathWithRoot += ext
	} //creating the actual file
	// if len(ext) == 0 {
	return os.Create(fullPathWithRoot)
	// }
	// return os.Create(fullPathWithRoot + ext)
}

func (s *Store) readStream(id string, key string) (int64, *os.File, error, string) {
	ext := getExtension(key)
	key = removeExtension(key)
	// fmt.Println("Key in read", key)

	patkey := s.PathTransformFunc(key)
	// fmt.Println("PathKey in read", patkey)
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, patkey.PathAndFileName())

	if len(ext) != 0 {
		fullPathWithRoot += ext
	}

	fd, err := os.Open(fullPathWithRoot)
	if err != nil {
		return 0, nil, err, ""
	}
	f, err := fd.Stat()
	if err != nil {
		return 0, nil, err, ""
	}
	fileLocation, _ := filepath.Abs(fullPathWithRoot)
	// fmt.Println("File location", fileLocation)

	return f.Size(), fd, nil, fileLocation

}

func (s *Store) writeStream(id, key string, r io.Reader) (int64, error, *os.File) {

	fd, err := s.openFileForWriting(id, key)
	if err != nil {
		return 0, err, nil
	}
	// defer fd.Close()
	// fmt.Println("File created")

	//writing to disk on the file location
	//remember our data is in the buf as we already copied it
	n, err := io.Copy(fd, r)
	if err != nil {
		return 0, err, nil
	}
	// fmt.Println("Data written to file")

	//logging the written bytes
	log.Printf("Written (%d) bytes to disk\n", n)
	// panic("not implemented")
	return n, nil, fd
}

func (s *Store) writeDecryptedStream(id, key string, encKey []byte, r io.Reader) (int, error) {
	fd, err := s.openFileForWriting(id, key)
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	return copyDecrypt(fd, r, encKey)

}

func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.FilePath, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func (p PathKey) PathAndFileName() string {
	return fmt.Sprintf("%s/%s", p.FilePath, p.FileName)
}

func getExtension(key string) string {
	return filepath.Ext(key)
}
