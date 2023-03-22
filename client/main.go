package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main(){
	connection, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
	}
	defer connection.Close()
	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')
	if err != nil{
		fmt.Println(err)
	}
	username = strings.Trim(username, " \r\n")
	welcomemsg := fmt.Sprintf("welcome")
	fmt.Println(welcomemsg)
}
