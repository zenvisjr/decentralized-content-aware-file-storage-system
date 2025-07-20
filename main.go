package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/zenvisjr/distributed-file-storage-system/p2p"
)

func init() {
	registerAll()
}

func main() {

	if err := p2p.EnsureKeyPair(); err != nil {
		log.Fatalf("Key generation failed: %v", err)
	}

	fmt.Println("LETS START COOKING")

	servers := completeServerSetup()

	s3 := servers[":5000"]
	s2 := servers[":4000"]
	s1 := servers[":3000"]
	s4 := servers[":8000"]

	time.Sleep(500 * time.Millisecond)

	s3.store.ClearRoot()
	s2.store.ClearRoot()
	s1.store.ClearRoot()
	s4.store.ClearRoot()

	// runCommandLoop(s3)

	key := "bill.pdf"
	f, err := os.Open(key)
	if err != nil {
		log.Fatal(err)
	}

	// key := "aizen"
	// // // f := bytes.NewReader([]byte("hello watashino soul society"))
	if err := s3.Store(key, f); err != nil {
		fmt.Println("Error storing data", err)
	}
	fmt.Println("---------------------------------------------------------------------------------------")

	// time.Sleep(500 * time.Millisecond)
	// if err := s3.Delete(key); err != nil {
	// 	fmt.Println("Error deleting data", err)
	// }

	// time.Sleep(500 * time.Millisecond)
	// if err := s3.DeleteLocal(key); err != nil {
	// 	fmt.Println("Error deleting data", err)
	// }

	// fmt.Println("---------------------------------------------------------------------------------------")

	// if err := s1.DeleteRemote(key); err != nil {
	// 	fmt.Println("Error deleting data", err)
	// }
	// fmt.Println("---------------------------------------------------------------------------------------")

	// if err := s2.DeleteRemote(key); err != nil {
	// 	fmt.Println("Error deleting data", err)
	// }

	// fmt.Println("---------------------------------------------------------------------------------------")

	// f, err = os.Open(key)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // // key := "aizen"
	// // // // // f := bytes.NewReader([]byte("hello watashino soul society"))
	// if err := s3.Store(key, f); err != nil {
	// 	fmt.Println("Error storing data", err)
	// }

	// // // time.Sleep(500 * time.Millisecond)

	// rd, flLocation, err := s3.Get(key)
	// if err != nil {
	// 	log.Fatal("Error getting data\n", err)
	// }
	// ext := getExtension(key)
	// if len(ext) == 0 {
	// 	n, err := io.ReadAll(rd)

	// 	if err != nil {
	// 		log.Fatal("Error reading data ", err)
	// 	}
	// 	fmt.Println(string(n))
	// }
	// fmt.Println(flLocation)
	fmt.Println("---------------------------------------------------------------------------------------")
}
