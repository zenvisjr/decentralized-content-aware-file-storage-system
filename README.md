<div id="top">

<!-- HEADER STYLE: COMPACT -->
<!-- <img src="readmeai/assets/logos/purple.svg" width="30%" align="left" style="margin-right: 15px"> -->

# DECENTRALIZED-CONTENT-AWARE-SECURE-FILE-STORAGE-SYSTEM
<em>Decentralized. Secure. Scalable. Distributed File Storage.</em>

<!-- BADGES -->
<!-- <img src="https://img.shields.io/github/license/zenvisjr/decentralized-secure-content-aware-file-storage-system?style=flat-square&logo=opensourceinitiative&logoColor=white&color=00ADD8" alt="license"> -->
<img src="https://img.shields.io/github/last-commit/zenvisjr/decentralized-secure-content-aware-file-storage-system?style=flat-square&logo=git&logoColor=white&color=00ADD8" alt="last-commit">
<img src="https://img.shields.io/github/languages/top/zenvisjr/decentralized-secure-content-aware-file-storage-system?style=flat-square&color=00ADD8" alt="repo-top-language">
<img src="https://img.shields.io/github/languages/count/zenvisjr/decentralized-secure-content-aware-file-storage-system?style=flat-square&color=00ADD8" alt="repo-language-count">
<img src="https://img.shields.io/badge/Go-00ADD8.svg?style=flat-square&logo=Go&logoColor=white" alt="Go">


<br clear="left"/>

## ☀️ Table of Contents

- [☀️ Table of Contents](#️-table-of-contents)
- [🌞 Overview](#-overview)
- [🔥 Features](#-features)
- [🌅 Project Structure](#-project-structure)
  - [🌄 Project Index](#-project-index)
- [🚀 Getting Started](#-getting-started)
  - [🌟 Prerequisites](#-prerequisites)
  - [⚡ Installation & Setup](#-installation--setup)
- [🔆 Usage](#-usage)
- [🏗️ System Architecture Overview](#️-system-architecture-overview)
- [🧩 Architecture Flow](#-architecture-flow)
- [🤝 Contributing](#-contributing)


---

## 🌞 Overview

The **Decentralized Secure Content-Aware File Storage System** is a peer-to-peer file storage solution designed for secure, efficient, and modular file management. It facilitates encrypted file transfer, decentralized storage, and resilient retrieval across interconnected nodes, making it suitable for systems where availability, integrity, and confidentiality are critical.

### Core Capabilities

* **🔒 Secure Data Handling**
  Implements AES encryption and RSA-based digital signature verification to ensure data confidentiality and authenticity during storage and retrieval.

* **🔄 Decentralized Architecture**
  Employs a peer-to-peer transport layer, removing dependency on central servers and enabling fault tolerance through distributed redundancy.

* **📂 Content-Aware File Storage**
  Each file is uniquely hashed and stored with awareness of its origin and integrity, enabling duplicate detection and secure content validation.

* **🧠 File Extension Preservation**
  Maintains original file extensions post-storage and during retrieval, ensuring format compatibility and accurate handling.

* **🚀 Intelligent File Retrieval**
  Utilizes local-first file lookup with peer fallback, optimizing retrieval speed and reducing network load.

* **🛠️ Configurable Startup**
  JSON-based configuration allows dynamic setup of node roles, ports, and peer links for flexible deployment.

* **👩‍💻 Automated Development Workflow**
  Makefile-based automation for building, running, and testing simplifies development and CI integration.


---

## 🔥 Features

| 🔧 Component                     | Description                                                                                                                                                |
| -------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| ⚙️ **Architecture**              | Peer-to-peer decentralized design eliminates single points of failure and enables robust, distributed file storage.                                        |
| 🛡️ **Security**                 | AES encryption and RSA-based digital signatures ensure confidentiality, integrity, and authenticity of stored content.                                     |
| 📂 **Content Awareness**         | Files are uniquely hashed and tracked with preserved extensions to enable content validation and duplication checks.                                       |
| 🚀 **Performance & Scalability** | Go's compiled nature, modular codebase, and decentralized design ensure high performance and easy scalability.                                             |
| 🧪 **Testing**                   | Unit tests using `testify` cover encryption, storage, and network components to ensure reliability.                                                        |
| 🔩 **Automation**                | Makefile-based build and run automation simplifies development and testing workflows.                                                                      |
| 🔌 **Custom P2P Protocol**       | Implemented a `p2p` package, a TCP-based peer discovery, handshake authentication using RSA keys, message streaming, and file transfer across nodes. |
| 🧱 **Modular Design**            | Clean separation of concerns across crypto, storage, network transport, and command handlers for easier maintenance and extensibility.                     |

---


## 🌅 Project Structure

```sh
└── decentralized-secure-content-aware-file-storage-system/
    ├── p2p
    │   ├── authentication.go
    │   ├── encoding.go
    │   ├── handshake.go
    │   ├── message.go
    │   ├── tcp_transport.go
    │   ├── tcp_transport_test.go
    │   └── transport.go
    ├── Makefile
    ├── README.md
    ├── crypto.go
    ├── crypto_test.go
    ├── go.mod
    ├── go.sum
    ├── main.go
    ├── makeServer.go
    ├── server.go
    ├── startServerConfig.json
    ├── store.go
    └── store_test.go
```

### 🌄 Project Index

<details open>
	<summary><b><code>DECENTRALIZED-SECURE-CONTENT-AWARE-FILE-STORAGE-SYSTEM/</code></b></summary>
	<!-- __root__ Submodule -->
	<details>
		<summary><b>__root__</b></summary>
		<blockquote>
			<div class='directory-path' style='padding: 8px 0; color: #666;'>
				<code><b>⦿ __root__</b></code>
			<table style='width: 100%; border-collapse: collapse;'>
			<thead>
				<tr style='background-color: #f8f9fa;'>
					<th style='width: 30%; text-align: left; padding: 8px;'>File Name</th>
					<th style='text-align: left; padding: 8px;'>Summary</th>
				</tr>
			</thead>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/makeServer.go'>makeServer.go</a></b></td>
					<td style='padding: 8px;'>Initializes and manages decentralized file server instances. Loads config files, sets up peers, starts the server, and handles commands like store, get, and delete within the peer-to-peer network.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/crypto.go'>crypto.go</a></b></td>
					<td style='padding: 8px;'>Provides cryptographic utilities for AES encryption/decryption, SHA-1 hashing, random ID/key generation, and RSA-SHA256 digital signature handling for data integrity and verification.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/go.mod'>go.mod</a></b></td>
					<td style='padding: 8px;'>Defines module dependencies and Go version for consistent builds.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/server.go'>server.go</a></b></td>
					<td style='padding: 8px;'>
					- Defines the <code>FileServer</code> struct to manage storage, peer connections, and network operations.<br>
					- Configures server options using <code>FileServerOps</code> (ID, root path, transport, encryption key, peers, etc.).<br>
					- Handles peer-to-peer communication using the custom <code>p2p</code> package.<br>
					- Processes file-related operations: store, retrieve, and delete.<br>
					- Manages file signature mapping and not-found response tracking.<br>
					- Initializes the file server via <code>NewFileServer()</code> constructor.<br>
					</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/Makefile'>Makefile</a></b></td>
					<td style='padding: 8px;'>Automates build, run, and test procedures for the Go codebase.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/startServerConfig.json'>startServerConfig.json</a></b></td>
					<td style='padding: 8px;'>Establishes the initial configuration for the network of servers, specifying the unique port, connected peers, and corresponding key path for each server.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/store.go'>store.go</a></b></td>
					<td style='padding: 8px;'>
					- Handles local file storage operations including write, read, and delete.<br>
					- Uses SHA-1 hash of the file name for consistent disk path generation.<br>
					- Preserves original file extensions during storage and retrieval.<br>
					- Manages RSA signature storage and validation alongside file content.<br>
					- Provides functions to check file existence and clear the storage directory.<br>
					- Supports writing decrypted files to disk after retrieval from peers.<br>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/main.go'>main.go</a></b></td>
					<td style='padding: 8px;'>- Acts as the entry point of the application.<br>
					- Parses command-line arguments to execute operations like <code>store</code>, <code>get</code>, <code>delete</code>, and <code>deletelocal</code>.<br>
					- Generates RSA key pairs if not already present.<br>
					- Loads configuration from <code>startServerConfig.json</code> to set up server and peers.<br>
					- Initializes and starts the <code>FileServer</code> instance.<br>
					- Facilitates interaction with the distributed system via CLI.<br>
				</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/go.sum'>go.sum</a></b></td>
					<td style='padding: 8px;'>It records the exact versions of dependencies used in the project and ensures that the project consistently uses the same versions of dependencies across various environments.</td>
				</tr>
			</table>
		</blockquote>
	</details>
	<!-- p2p Submodule -->
	<details>
		<summary><b>p2p</b></summary>
		<blockquote>
			<div class='directory-path' style='padding: 8px 0; color: #666;'>
				<code><b>⦿ p2p</b></code>
			<table style='width: 100%; border-collapse: collapse;'>
			<thead>
				<tr style='background-color: #f8f9fa;'>
					<th style='width: 30%; text-align: left; padding: 8px;'>File Name</th>
					<th style='text-align: left; padding: 8px;'>Summary</th>
				</tr>
			</thead>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/p2p/handshake.go'>handshake.go</a></b></td>
					<td style='padding: 8px;'>- Implements the handshake protocol between peers during connection establishment.<br>
					- Exchanges node IDs and public RSA keys to verify identity and enable secure communication.<br>
					- Prevents self-connections and duplicate connections using defensive checks.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/p2p/message.go'>message.go</a></b></td>
					<td style='padding: 8px;'>- Defines communication structures within the peer-to-peer network module<br>- It establishes data types for Incoming messages and streams, facilitating the exchange of arbitrary data between peers<br>- The flexibility of the RPC structure allows for varied data payloads and supports both regular and streaming data transmission modes.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/p2p/transport.go'>transport.go</a></b></td>
					<td style='padding: 8px;'>- Transport.go in the p2p package facilitates communication between nodes in a P2P network<br>- It defines interfaces for peers and transport systems, establishing methods for handling connections, sending data, and managing file hashmaps<br>- Its functionality supports various transport forms such as TCP, UDP, and WebSockets.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/p2p/tcp_transport.go'>tcp_transport.go</a></b></td>
					<td style='padding: 8px;'>- Implements the TCP-based transport layer for peer communication.<br>
					- Defines the <code>TCPTransport</code> struct to listen for and manage TCP connections.<br>
					- Handles incoming connections, performs handshake, and manages active peers.<br>
					- Defines the <code>TCPPeer</code> struct representing a connected remote peer.<br>
					- Supports message sending, peer removal, and stream handling over TCP sockets.<br>
					- Integrates with the <code>handshake.go</code> and <code>message.go</code> components for secure messaging.<br>
				</td>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/p2p/encoding.go'>encoding.go</a></b></td>
					<td style='padding: 8px;'>- Encoding.go in the p2p package provides interfaces and implementations for decoding incoming data in the peer-to-peer network communication<br>- It includes a default decoder and a GOB decoder, each offering a method to read from a stream and decode incoming RPC messages, handling both standard data and stream data.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/p2p/authentication.go'>authentication.go</a></b></td>
					<td style='padding: 8px;'>- The p2p/authentication.go is instrumental in the peer-to-peer networks security ecosystem<br>- It ensures the creation and availability of RSA key pairs for secure communication<br>- The functionality includes generating a new key pair if one doesnt exist, and loading an existing public or private key when necessary.</td>
				</tr>
			</table>
		</blockquote>
	</details>
</details>

---

Thanks! Based on your updated `Makefile` and `startServerConfig.json` structure, here's the revised and complete **🚀 Getting Started** section with everything integrated:

---

## 🚀 Getting Started

### 🌟 Prerequisites

Make sure you have the following installed:

* **Go** (v1.20 or higher) — [Install Go](https://golang.org/dl/)
* **Git** — to clone the repository

---

### ⚡ Installation & Setup

Follow these steps to set up and run the system locally:

1. **Clone the repository:**

   ```sh
   git clone https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system
   ```

2. **Navigate to the project repository:**

   ```sh
   cd decentralized-secure-content-aware-file-storage-system
   ```

2. **Build the project:**

   Use the provided `Makefile`:

   ```sh
   make build
   ```

   This will compile the binary to `./endgame/yeahbaby`.

3. **Define peer configuration:**

   Modify or create a `startServerConfig.json` file with the structure below:

   ```json
   [
     {
       "port": ":3000",
       "peers": ["", ":8000"],
       "key_path": "3000.key"
     },
     {
       "port": ":4000",
       "peers": [":3000"],
       "key_path": "4000.key"
     },
     {
       "port": ":8000",
       "peers": [],
       "key_path": "8000.key"
     },
     {
       "port": ":5000",
       "peers": [":3000", ":4000", ":8000"],
       "key_path": "5000.key"
     }
   ]
   ```

   Each server entry defines:

   * Listening port
   * Peers it connects to
   * Path to its session key file (otherwise auto-generated each time if not stored)

4. **Run the server:**

   ```sh
   make run
   ```

5. **Run all tests:**

   ```sh
   make test
   ```

---


Here’s a complete **🔆 Usage** section to showcase the CLI commands supported by your system, with clear examples and expected behavior:

---

## 🔆 Usage

After running the server (`make run`), you can use the following CLI operations to interact with the decentralized file storage system:

### 📥 Store a File

Encrypts and stores the specified file on the local node and broadcasts it to peers.

```sh
> store <file_path>
```

**Example:**

```sh
> store samples/image.png
```

---

### 📤 Retrieve a File

Retrieves the file from the local store or fetches it from peers if not found locally. Decrypts and restores the original file with preserved extension.

```sh
> get <file_key>
```

**Example:**

```sh
> get image.png
```

---

### ❌ Delete from All Peers

Deletes the file from the local node and all connected peers. Also removes its associated digital signature.

```sh
> delete <file_key>
```

---

### 🗑️ Delete Locally Only

Removes the file and its signature only from the current node without affecting peers.

```sh
> deletelocal <file_key>
```

---

### 🚪 Quit the Server Loop
Exits the infinite input loop and gracefully shuts down the server.

```sh
> quit
```


---

Great. Let's start with the **complete architecture flow** — explained step-by-step with clear module interactions:

---

## 🏗️ System Architecture Overview

### 🔄 Overall Flow – How Components Interact

```plaintext
               ┌──────────────────────┐
               │  startServerConfig   │
               │  (JSON config file)  │
               └─────────┬────────────┘
                         │
                         ▼
               ┌──────────────────────┐
               │   makeServer.go      │
               │ - Initializes        │
               │   FileServer via     │
               │   NewFileServer()    │
               └─────────┬────────────┘
                         │
                         ▼
               ┌──────────────────────┐
               │     main.go          │
               │ - Parses config      │
               │ - Starts server      │
               │ - Accepts CLI input  │
               └─────────┬────────────┘
                         │
                         ▼
               ┌────────────────────────────┐
               │        server.go           │
               │ - Defines FileServer logic │
               │ - Handles file operations  │
               │ - Manages peer map         │
               └─────────┬────────────┬─────┘
                         │            │
                         ▼            ▼
         ┌────────────────────┐   ┌────────────────────┐
         │     store.go       │   │     crypto.go       │
         │ - Local storage    │   │ - AES encryption    │
         │ - SHA1 pathing     │   │ - RSA signature     │
         └────────┬───────────┘   └────────┬────────────┘
                  │                        │
                  ▼                        ▼
        ┌─────────────────────┐   ┌─────────────────────────┐
        │     p2p package     │   │   File Signature Map     │
        │ - TCP transport     │   │ - In-memory tracking     │
        │ - Handshake auth    │   │   of signatures per file │
        │ - Peer messaging    │   └─────────────────────────┘
        └────────┬────────────┘
                 ▼
      ┌──────────────────────────┐
      │      message.go          │
      │ - Defines RPC structure  │
      │ - Supports stream flags  │
      └──────────────────────────┘
```


---

## 🧩 Architecture Flow

**main.go** is the entry point:

* Calls `EnsureKeyPair()` → checks/generates `private.pem` and `public.pem`
* Calls `completeServerSetup()` → loads all servers from `startServerConfig.json`

**completeServerSetup()** (in makeServer.go):

* Loads JSON configs for each server (port, peer list, key path)
* For each config:

  * Calls `makeServer()`
  * Inside `makeServer()`:

    * Initializes TCP transport using `p2p.NewTCPTransport()`
    * Creates a `FileServer` using `NewFileServer()` with:

      * Port
      * Peer addresses
      * AES encryption key
      * Path transformation logic
    * Assigns the file server’s `OnPeer()` handler to the TCP transport
    * Starts the server via `server.Start()` (in a goroutine)
* After starting the server, we execute `runCommandLoop()` to accept CLI input in a infinite loop until `quit` is entered.

**Result:** All servers are initialized, their peers are connected via TCP, and ready to accept commands.

---



## 🤝 Contributing


<details closed>
<summary>Contributing Guidelines</summary>

1. **Fork the Repository**: Start by forking the project repository to your github account.
2. **Clone Locally**: Clone the forked repository to your local machine using a git client.
   ```sh
   git clone https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system
   ```
3. **Create a New Branch**: Always work on a new branch, giving it a descriptive name.
   ```sh
   git checkout -b new-feature-x
   ```
4. **Make Your Changes**: Develop and test your changes locally.
5. **Commit Your Changes**: Commit with a clear message describing your updates.
   ```sh
   git commit -m 'Implemented new feature x.'
   ```
6. **Push to github**: Push the changes to your forked repository.
   ```sh
   git push origin new-feature-x
   ```
7. **Submit a Pull Request**: Create a PR against the original project repository. Clearly describe the changes and their motivations.
8. **Review**: Once your PR is reviewed and approved, it will be merged into the main branch. Congratulations on your contribution!
</details>

<details closed>
<summary>Contributor Graph</summary>
<br>
<p align="left">
   <a href="https://github.com{/zenvisjr/decentralized-secure-content-aware-file-storage-system/}graphs/contributors">
      <img src="https://contrib.rocks/image?repo=zenvisjr/decentralized-secure-content-aware-file-storage-system">
   </a>
</p>
</details>

---


## 👨‍💻 Author

**Ayush Rai**  
📧 Email: [ayushrai.cse@gmail.com](mailto:ayushrai.cse@gmail.com)

---




<div align="right">

[![][back-to-top]](#top)

</div>


[back-to-top]: https://img.shields.io/badge/-BACK_TO_TOP-151515?style=flat-square


