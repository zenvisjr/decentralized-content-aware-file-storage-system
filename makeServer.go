package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
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
		ShakeHands:       p2p.NOHandshakeFunc,
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
