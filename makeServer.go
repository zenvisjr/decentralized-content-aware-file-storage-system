package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zenvisjr/distributed-file-storage-system/p2p"
)

type ServerConfig struct {
	Port           string   `json:"port"`
	BootstrapNodes []string `json:"peers"`
	KeyPath        string   `json:"key_path"`
}

func loadConfig(path string) ([]ServerConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg []ServerConfig
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}
	// fmt.Println(cfg)
	return cfg, nil
}

func makeServer(configFile *ServerConfig) *FileServer {
	tcpops := p2p.TCPTransportOps{
		ListenerPortAddr: configFile.Port,
		ShakeHands:       p2p.DefensiveHandshakeFunc,
		Decoder:          &p2p.DefaultDecoder{},
		OnPeer:           nil,
	}

	tcpTransport, _ := p2p.NewTCPTransport(tcpops)

	// var key []byte

	//checking if key is provided by user via config file, if not then generate a new key and save it to the key path
	// if _, err := os.Stat(configFile.KeyPath); err == nil {
	// 	key, _ = os.ReadFile(configFile.KeyPath)
	// } else {
	// 	key = newEncryptionKey()
	// 	_ = os.WriteFile(configFile.KeyPath, key, 0644)
	// }

	//for windows as we cant start folder name with : so we need to remove it

	newAddr := strings.ReplaceAll(configFile.Port, ":", "_")
	fileServerOps := FileServerOps{
		RootStorage:       newAddr + "_gyattt",
		PathTransformFunc: CryptoPathTransformFunc,
		Transort:          tcpTransport,
		BootstrapNodes:    configFile.BootstrapNodes,
		EncKey:            newEncryptionKey(),
		// EncKey:            []byte("yokoso"),
	}
	// fmt.Println(fileServerOps)

	newFileServer, _ := NewFileServer(fileServerOps)

	//assigning the onpeer func made in server.go to the TCPtransport in tcp_transport.go
	tcpTransport.OnPeer = newFileServer.OnPeer

	return newFileServer
}

// func makeServer(port string, bootstrapNodes ...string) *FileServer {
// 	tcpops := p2p.TCPTransportOps{
// 		ListenerPortAddr: port,
// 		ShakeHands:       p2p.NOHandshakeFunc,
// 		Decoder:          &p2p.DefaultDecoder{},
// 		OnPeer:           nil,
// 	}

// 	tcpTransport, _ := p2p.NewTCPTransport(tcpops)

// 	//for windows as we cant start folder name with : so we need to remove it

// 	newAddr := strings.ReplaceAll(port, ":", "_")
// 	fileServerOps := FileServerOps{
// 		RootStorage:       newAddr + "_gyattt",
// 		PathTransformFunc: CryptoPathTransformFunc,
// 		Transort:          tcpTransport,
// 		BootstrapNodes:    bootstrapNodes,
// 		EncKey:            newEncryptionKey(),
// 	}
// 	// fmt.Println(fileServerOps)

// 	newFileServer, _ := NewFileServer(fileServerOps)

// 	//assigning the onpeer func made in server.go to the TCPtransport in tcp_transport.go
// 	tcpTransport.OnPeer = newFileServer.OnPeer

// 	return newFileServer
// }

func completeServerSetup() map[string]*FileServer {
	configPath := flag.String("config", "startServerConfig.json", "path of config file to create servers")
	flag.Parse()

	configs, err := loadConfig(*configPath)
	if err != nil {
		log.Fatal("Error loading config file", err)
	}

	//storing each server according to its port
	servers := make(map[string]*FileServer)
	for _, cfg := range configs {
		// fmt.Println("Bootstrap nodes", cfg.BootstrapNodes)
		server := makeServer(&cfg)
		servers[cfg.Port] = server
		go func(server *FileServer) {
			log.Fatal(server.Start())
		}(server)
		time.Sleep(50 * time.Millisecond)
	}

	return servers
}

func runCommandLoop(fs *FileServer) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(">>> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}

		cmd := strings.ToLower(args[0])

		switch cmd {
		case "store":
			fmt.Println(args)
			if len(args) != 2 {
				fmt.Println("Usage: store <file-path>")
				continue
			}
			f, err := os.Open(args[1])
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			key := filepath.Base(args[1])
			// f := bytes.NewReader([]byte("hello watashino soul society"))
			// key := "hello"
			if err := fs.Store(key, f); err != nil {
				fmt.Println("Error storing data", err)
			}
			time.Sleep(2 * time.Second)

		case "get":
			if len(args) != 2 {
				fmt.Println("Usage: get <filename>")
				continue
			}
			key := args[1]
			rd, fileLoc, err := fs.Get(key)
			if err != nil {
				fmt.Println("Error getting file:", err)
			} else {
				fmt.Println("File stored at:", fileLoc)
			}

			ext := getExtension(key)
			if len(ext) == 0 {
				n, err := io.ReadAll(rd)

				if err != nil {
					log.Fatal("Error reading data ", err)
				}
				fmt.Println(string(n))

			}

		case "delete":
			if len(args) != 2 {
				fmt.Println("Usage: delete <filename>")
				continue
			}
			key := args[1]
			err := fs.Delete(key)
			if err != nil {
				fmt.Println("Error deleting file:", err)
			}

		case "deletelocal":
			if len(args) != 2 {
				fmt.Println("Usage: deletelocal <filename>")
				continue
			}
			key := args[1]
			err := fs.DeleteLocal(key)
			if err != nil {
				fmt.Println("Error deleting local file:", err)
			}

		case "quit":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Unknown command. Supported: store, get, delete, deletelocal, quit")
		}
		time.Sleep(50 * time.Millisecond)
	}
}
