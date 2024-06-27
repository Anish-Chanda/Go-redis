package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	//listen to tcp reqs on 6379
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on port 6379")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		//return default pong for now
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Handling new connection")

	for {
		buffer := make([]byte, 256)
		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			return
		}
		conn.Write([]byte("+PONG\r\n"))
	}
}
