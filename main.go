package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {

	fmt.Println("LETS START COOKING")
	servers := completeServerSetup()
	s3 := servers[":5000"]
	s2 := servers[":4000"]
	s1 := servers[":3000"]
	s4 := servers[":8000"]

	// s1 := makeServer(":3000", "", ":8000")
	// s2 := makeServer(":4000", ":3000")
	// s3 := makeServer(":5000", ":3000", ":4000", ":8000")
	// s4 := makeServer(":8000", "")

	// go func() {
	// 	log.Fatal(s1.Start())
	// }()

	// time.Sleep(50 * time.Millisecond)

	// go func() {
	// 	log.Fatal(s2.Start())
	// }()

	// time.Sleep(50 * time.Millisecond)

	// go func() {
	// 	log.Fatal(s3.Start())
	// }()

	// go func() {
	// 	log.Fatal(s4.Start())
	// }()

	// time.Sleep(50 * time.Millisecond)

	time.Sleep(1 * time.Second)

	s3.store.ClearRoot()
	s2.store.ClearRoot()
	s1.store.ClearRoot()
	s4.store.ClearRoot()

	// createFile()

	key := "test.zip"

	f, err := os.Open(key)
	if err != nil {
		log.Fatal(err)
	}
	// ext := filepath.Ext(f.Name())
	// size, _ := f.Stat()

	defer f.Close()

	// log.Println("Done opening file")

	// start := time.Now()
	// for i := 0; i < 1; i++ {
	// key := fmt.Sprintf("bill.pdf", i)

	// f := bytes.NewReader([]byte("hello watashino soul society"))
	if err := s3.Store(key, f); err != nil {
		fmt.Println("Error storing data", err)
	}

	time.Sleep(500 * time.Millisecond)
	if err := s3.Delete(key); err != nil {
		fmt.Println("Error deleting data", err)
	}

	// time.Sleep(500 * time.Millisecond)
	// if err := s3.DeleteLocal(key); err != nil {
	// 	fmt.Println("Error deleting data", err)
	// }

	time.Sleep(500 * time.Millisecond)

	rd, fileLocation, err := s3.Get(key)
	if err != nil {
		log.Fatal("Error getting data\n", err)
	}
	ext := getExtension(key)
	if len(ext) == 0 {
		n, err := io.ReadAll(rd)

		if err != nil {
			log.Fatal("Error reading data ", err)
		}
		fmt.Println(string(n))

	}
	fmt.Println("Retrived file location", fileLocation)

}

// fmt.Println("Time taken to store file", time.Since(start))
// }

func init() {
	gob.Register(MessageStoreFile{}) // Register the pointer form too
	gob.Register(MessageGetFile{})
	gob.Register(MessageDeleteFile{})
}
