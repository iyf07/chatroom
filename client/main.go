package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gookit/color"
	"io"
	"log"
	"net"
	"os"
	"strings"
	time2 "time"
)

func main() {
	// Parse command-line flag
	var host = flag.String("host", "86.51", "host-ip")
	var port = flag.String("port", "8080", "port")
	var name = flag.String("name", "User", "name")
	var message = flag.String("message", "Hello", "message")
	flag.Parse()

	// Get client IP address
	clientIP, err := getIPAddress()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the TCP server
	serverIP := fmt.Sprintf("192.168.%s:%s", *host, *port)
	connection, err := net.Dial("tcp", serverIP)
	if err != nil {
		log.Fatal(err)
	}
	defer func(connection net.Conn) {
		err = connection.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(connection)

	fmt.Println("Connected to the TCP server " + serverIP)
	fmt.Println("Your IP address is " + clientIP)

	initialMessage := fmt.Sprintf("%s(%s) joined: %s\n", *name, clientIP, *message)
	_, err = connection.Write([]byte(initialMessage))
	if err != nil {
		log.Fatal(err)
	}
	go read(connection)
	write(connection, *name, clientIP)
}

func read(connection net.Conn) {
	for {
		reader := bufio.NewReader(connection)
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			err := connection.Close()
			if err != nil {
				return
			}
			fmt.Println("Connection closed")
			os.Exit(0)
		}
		color.Magenta.Println(message)
	}
}

func write(connection net.Conn, name string, clientIP string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		time := time2.Now().String()[11:19]
		message = fmt.Sprintf("%s %s(%s): %s\n", time, name, clientIP, strings.Trim(message, " \r\n"))
		_, err = connection.Write([]byte(message))
		if err != nil {
			break
		}
	}
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
