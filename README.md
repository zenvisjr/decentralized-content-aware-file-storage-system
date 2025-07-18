<div id="top">

<!-- HEADER STYLE: COMPACT -->
<!-- <img src="readmeai/assets/logos/purple.svg" width="30%" align="left" style="margin-right: 15px"> -->

# DECENTRALIZED-CONTENT-AWARE-SECURE-FILE-STORAGE-SYSTEM
<em>Decentralized. Secure. Scalable. Redefining Distributed File Storage.</em>

<!-- BADGES -->
<!-- <img src="https://img.shields.io/github/license/zenvisjr/decentralized-secure-content-aware-file-storage-system?style=flat-square&logo=opensourceinitiative&logoColor=white&color=00ADD8" alt="license"> -->
<img src="https://img.shields.io/github/last-commit/zenvisjr/decentralized-secure-content-aware-file-storage-system?style=flat-square&logo=git&logoColor=white&color=00ADD8" alt="last-commit">
<img src="https://img.shields.io/github/languages/top/zenvisjr/decentralized-secure-content-aware-file-storage-system?style=flat-square&color=00ADD8" alt="repo-top-language">
<img src="https://img.shields.io/github/languages/count/zenvisjr/decentralized-secure-content-aware-file-storage-system?style=flat-square&color=00ADD8" alt="repo-language-count">
<img src="https://img.shields.io/badge/Go-00ADD8.svg?style=flat-square&logo=Go&logoColor=white" alt="Go">


<br clear="left"/>

## ☀️ Table of Contents

- [☀ ️ Table of Contents](#table-of-contents)
- [🌞 Overview](#overview)
- [🔥 Features](#features)
- [🌅 Project Structure](#project-structure)
    - [🌄 Project Index](#project-index)
- [🚀 Getting Started](#getting-started)
    - [🌟 Prerequisites](#-prerequisites)
    - [⚡ Installation](#-installation)
    - [🔆 Usage](#-usage)
    - [🌠 Testing](#-testing)
- [🌻 Roadmap](#-roadmap)
- [🤝 Contributing](#-contributing)
- [📜 License](#-license)
- [✨ Acknowledgments](#-acknowledgments)



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
					<td style='padding: 8px;'>- MakeServer.go is primarily responsible for setting up and managing a decentralized file server in this distributed file storage system<br>- It facilitates the creation of servers, handling of configuration files, and execution of user commands such as storing, retrieving, and deleting files<br>- This forms a vital part of the projects peer-to-peer communication infrastructure.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/crypto.go'>crypto.go</a></b></td>
					<td style='padding: 8px;'>- Crypto.go serves as the cryptographic core of the project, providing functionality for secure data handling<br>- It facilitates the generation of random IDs and encryption keys, offers SHA-1 hashing, and is instrumental in both encrypting and decrypting streams of data<br>- Additionally, it implements signature functionalities using RSA and SHA256 for data verification.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/go.mod'>go.mod</a></b></td>
					<td style='padding: 8px;'>- Go.mod in the project distributed-file-storage-system defines the module path and specifies the Go language version and required dependencies<br>- It aids in ensuring consistent, reproducible builds by pinning specific versions of dependencies, including go-spew, go-difflib, testify and yaml.v3.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/server.go'>server.go</a></b></td>
					<td style='padding: 8px;'>- Server.go<code>The </code>server.go<code> file is integral to the operation of the Distributed File Storage System<br>- This file primarily serves as the backbone of the file server within the system, and it is responsible for managing the distributed file storage architecture.The key purpose of </code>server.go<code> is to define the operations and structure of the file server, and to initialize new file server instances<br>- The file server is designed to handle storage operations, peer-to-peer communication, and encryption key management.The </code>FileServerOps<code> struct is a collection of options for the file server, including its ID, root storage location, path transformation function, transport protocol, bootstrap nodes, and encryption key.The </code>FileServer<code> struct represents the file server itself, incorporating the server's options, a store for the server's files, a quit channel, a lock for managing peer access, a map of peers, and a channel for handling not found errors.The </code>NewFileServer<code> function is used to create a new instance of FileServer<br>- It takes a </code>FileServerOps<code> struct as an argument and returns a new </code>FileServer<code> instance.This file is part of the main package and leverages the </code>p2p` package for peer-to-peer communication capabilities<br>- It is central to the functioning of the Distributed File Storage System and plays a critical role in the overall codebase architecture.Please refer to the individual function documentation for more specific details on the implementation and use of each function.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/store_test.go'>store_test.go</a></b></td>
					<td style='padding: 8px;'>- Store_test.go conducts unit tests for validating the functionalities of a storage system<br>- It checks the transformation of file paths using a cryptographic function, file storage, and retrieval operations<br>- Furthermore, it ensures the system correctly handles file deletions, and the storage system can be properly cleared without errors.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/crypto_test.go'>crypto_test.go</a></b></td>
					<td style='padding: 8px;'>- Crypto_test.go serves as a crucial testing unit within the codebase, focusing on validating the encryption and decryption functions<br>- It creates a mock scenario to simulate the process of copying and encrypting a large file, followed by its decryption, ensuring the integrity and effectiveness of these procedures.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/Makefile'>Makefile</a></b></td>
					<td style='padding: 8px;'>- Makefile serves as a script for automating the build, run, and test procedures of the Go codebase<br>- It compiles the Go source code into a binary file, executes the binary, and runs unit tests across all packages in the project, providing a streamlined development workflow.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/startServerConfig.json'>startServerConfig.json</a></b></td>
					<td style='padding: 8px;'>- StartServerConfig.json establishes the initial configuration for the network of servers, specifying the unique port, connected peers, and corresponding key path for each server<br>- Its role within the entire codebase is to dictate the fundamental topology and security parameters of the server network.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/store.go'>store.go</a></b></td>
					<td style='padding: 8px;'>- Store.go serves as the main storage module, facilitating file operations such as reading, writing, and deleting files on disk<br>- It organizes files using a SHA1 hashing mechanism and supports storing and retrieving file signatures<br>- Additionally, it provides a means to check file existence, clear the root directory, and handle decrypted file writing.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/main.go'>main.go</a></b></td>
					<td style='padding: 8px;'>- Main.go serves as the central orchestrator for a distributed file storage system<br>- It manages key operations such as key pair generation, server setup, file handling, and command execution<br>- The file facilitates interactions with multiple servers, ensuring data storage and retrieval while maintaining a clear root across the system.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/go.sum'>go.sum</a></b></td>
					<td style='padding: 8px;'>- Go.sum manages the specific versions of dependencies used in the project<br>- It ensures that the project consistently uses the same versions of dependencies across various environments<br>- This helps in maintaining stability, reproducibility, and integrity in the application by preventing unnoticed updates or changes in the dependencies.</td>
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
					<td style='padding: 8px;'>- The <code>handshake.go</code> file in the <code>p2p</code> directory manages peer-to-peer network handshakes in the application<br>- It performs the exchange of network node IDs, and also the exchange of public keys for securing communication<br>- The file defines functions for both simple and defensive handshakes, handling potential issues such as self-connections and handshake timeouts.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/p2p/message.go'>message.go</a></b></td>
					<td style='padding: 8px;'>- In the broader context of the project, <code>message.go</code> defines communication structures within the peer-to-peer network module<br>- It establishes data types for Incoming messages and streams, facilitating the exchange of arbitrary data between peers<br>- The flexibility of the RPC structure allows for varied data payloads and supports both regular and streaming data transmission modes.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/p2p/transport.go'>transport.go</a></b></td>
					<td style='padding: 8px;'>- Transport.go in the p2p package facilitates communication between nodes in a P2P network<br>- It defines interfaces for peers and transport systems, establishing methods for handling connections, sending data, and managing file hashmaps<br>- Its functionality supports various transport forms such as TCP, UDP, and WebSockets.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/p2p/tcp_transport.go'>tcp_transport.go</a></b></td>
					<td style='padding: 8px;'>- TCPTransport, located at p2p/tcp_transport.go, establishes and manages TCP connections for a peer-to-peer network<br>- It handles incoming messages, manages peer connections, and maintains a hash map of files<br>- The TCPPeer struct represents the remote node over a TCP connection<br>- These components together facilitate efficient network communication and data transfers.</td>
				</tr>
				<tr style='border-bottom: 1px solid #eee;'>
					<td style='padding: 8px;'><b><a href='https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/master/p2p/tcp_transport_test.go'>tcp_transport_test.go</a></b></td>
					<td style='padding: 8px;'>- TCPTransportTest validates the functionality of the TCP Transport protocol in a peer-to-peer(P2P) system<br>- The test ensures that a TCP connection can be successfully established, data can be sent, and the connection can be closed without errors<br>- Its a crucial component of the project, ensuring robust networking in the P2P system.</td>
				</tr>
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

## 🚀 Getting Started

### 🌟 Prerequisites

This project requires the following dependencies:

- **Programming Language:** Go
- **Package Manager:** Go modules

### ⚡ Installation

Build decentralized-secure-content-aware-file-storage-system from the source and intsall dependencies:

1. **Clone the repository:**

    ```sh
    ❯ git clone https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system
    ```

2. **Navigate to the project directory:**

    ```sh
    ❯ cd decentralized-secure-content-aware-file-storage-system
    ```

3. **Install the dependencies:**

<!-- SHIELDS BADGE CURRENTLY DISABLED -->
	<!-- [![go modules][go modules-shield]][go modules-link] -->
	<!-- REFERENCE LINKS -->
	<!-- [go modules-shield]: https://img.shields.io/badge/Go-00ADD8.svg?style={badge_style}&logo=go&logoColor=white -->
	<!-- [go modules-link]: https://golang.org/ -->

	**Using [go modules](https://golang.org/):**

	```sh
	❯ go build
	```

### 🔆 Usage

Run the project with:

**Using [go modules](https://golang.org/):**
```sh
go run {entrypoint}
```

### 🌠 Testing

Decentralized-secure-content-aware-file-storage-system uses the {__test_framework__} test framework. Run the test suite with:

**Using [go modules](https://golang.org/):**
```sh
go test ./...
```

---

## 🌻 Roadmap

- [X] **`Task 1`**: <strike>Implement feature one.</strike>
- [ ] **`Task 2`**: Implement feature two.
- [ ] **`Task 3`**: Implement feature three.

---

## 🤝 Contributing

- **💬 [Join the Discussions](https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/discussions)**: Share your insights, provide feedback, or ask questions.
- **🐛 [Report Issues](https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/issues)**: Submit bugs found or log feature requests for the `decentralized-secure-content-aware-file-storage-system` project.
- **💡 [Submit Pull Requests](https://github.com/zenvisjr/decentralized-secure-content-aware-file-storage-system/blob/main/CONTRIBUTING.md)**: Review open PRs, and submit your own PRs.

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

## 📜 License

Decentralized-secure-content-aware-file-storage-system is protected under the [LICENSE](https://choosealicense.com/licenses) License. For more details, refer to the [LICENSE](https://choosealicense.com/licenses/) file.

---

## ✨ Acknowledgments

- Credit `contributors`, `inspiration`, `references`, etc.

<div align="right">

[![][back-to-top]](#top)

</div>


[back-to-top]: https://img.shields.io/badge/-BACK_TO_TOP-151515?style=flat-square


---
