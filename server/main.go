package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

var openConnections = make(map[net.Conn]bool)
var newConnection = make(chan net.Conn)
var deadConnection = make(chan net.Conn)
var initialMessages []string

func main() {
	// Parse command-line flag
	var port = flag.String("port", "8080", "port")
	flag.Parse()

	// Get server IP address
	serverIP, err := getIPAddress()
	if err != nil {
		log.Fatal(err)
	}

	// Create server
	ln, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("TCP server listens on " + serverIP + ":" + *port)

	// Accept incoming connection requests
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Fatal(err)
			}
			openConnections[conn] = true
			newConnection <- conn
		}
	}()

	// Listen for messages on channels
	for {
		select {
		case conn := <-newConnection:
			broadcastInitialMessages(conn)
			go broadcastMessage(conn)

		case conn := <-deadConnection:
			for item := range openConnections {
				if item == conn {
					break
				}
			}
			delete(openConnections, conn)
		}
	}
}

func broadcastInitialMessages(conn net.Conn) {
	// Read new initial message
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	// Send old initial messages to the current connection
	for _, initialMessage := range initialMessages {
		_, err := conn.Write([]byte(initialMessage))
		if err != nil {
			break
		}
	}

	// Send new initial message to all open connections
	for openConnection := range openConnections {
		_, err := openConnection.Write([]byte(message))
		if err != nil {
			break
		}
	}

	// Store the new initial message
	initialMessages = append(initialMessages, message)
}

func broadcastMessage(conn net.Conn) {
	for {
		// Read message
		reader := bufio.NewReader(conn)
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		// Send message to all open connections
		for openConnection := range openConnections {
			_, err := openConnection.Write([]byte(message))
			if err != nil {
				break
			}
		}
	}
	deadConnection <- conn
}

func getIPAddress() (string, error) {
	// Get all network interfaces
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	// Loop through the network interfaces and print the IPv4 addresses
	for _, netInterface := range netInterfaces {

		// Ignore loopback and other non-physical interfaces
		if netInterface.Flags&net.FlagLoopback != 0 || netInterface.Flags&net.FlagUp == 0 {
			continue
		}

		// Get a list of unicast addresses for the interface
		addrs, err := netInterface.Addrs()
		if err != nil {
			fmt.Println("Failed to get addresses for interface", netInterface.Name, ":", err)
			continue
		}

		// Loop through the unicast addresses and print the IPv4 addresses
		for _, addr := range addrs {

			// Check if the address is an IPv4 address
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", nil
}
