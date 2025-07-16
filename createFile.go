package main

import (
	"log"
	"os"
)

func createFile() {
	//creating a big file of 1gb
	f, err := os.Create("bigfile.dat")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	buf := make([]byte, 1024*1024) // 1 MB
	for i := range buf {
		buf[i] = byte(i % 256)
	}

	log.Println("Writing file...")

	for i := 0; i < 500; i++ { // 1024 * 1MB = 1 GB
		if _, err := f.Write(buf); err != nil {
			log.Fatal(err)
		}
	}
}
