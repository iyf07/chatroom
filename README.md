# chatroom
This is a Golang app that creates a chatroom and allows communications among devices in the same local network.

## Usage
### `go run ./server/main.go` 
Creates a TCP server.
- `-port`: port that the TCP server runs on. Default is `8080`.
### `go run ./client/main.go` 
Creates a client and connects to the TCP server.
- `-host`: host component of the TCP server IP address. Default is `86.51`.
- `-port`: port that the TCP server runs on. Default is `8080`.
- `-name`: client username. Default is `User`.
- `-message`: the initial message a client is sending. Default is `Hello`.

## Example
```
Client1(192.168.86.51) joined: Message1

Client2(192.168.86.26) joined: Message2

23:21:41 Client1(192.168.86.51): Hello from Client 1!

23:21:55 Client2(192.168.86.26): Hello from Client 2!
```