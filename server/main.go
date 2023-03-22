package main

import (
	"fmt"
	"log"
	"net"
)

var openConnections = make(map[net.Conn]bool)
var newConnection = make(chan net.Conn)
var deadConnection = make(chan net.Conn)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}

go func () {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		openConnections[conn] = true
		newConnection <- conn
	}
}()

fmt.Println(<- newConnection)
}
