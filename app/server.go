package main

import (
	"fmt"
	"net"
	"os"

	"github.com/anish-chanda/goredis/helpers"
)

var tempStore map[string]interface{} = make(map[string]interface{})

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
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			return
		}

		command := string(buffer[:n])
		fmt.Println("Received command ", command)

		data, err := helpers.RespParser(buffer)
		if err != nil {
			fmt.Println("Error parsing command: ", err.Error())
		}
		if data.Command == "PING" {
			conn.Write([]byte("+PONG\r\n"))
		} else if data.Command == "ECHO" {
			response := fmt.Sprintf("$%d\r\n%s\r\n", len(data.Args[0]), data.Args[0])
			conn.Write([]byte(response))
		} else if data.Command == "SET" {
			//the key will be arg[0] and val will be arg[1]
			tempStore[data.Args[0]] = data.Args[1]
			conn.Write([]byte("+OK\r\n"))
		} else if data.Command == "GET" {
			res, ok := tempStore[data.Args[0]]
			if !ok {
				conn.Write([]byte("$-1\r\n"))
			} else {
				response := fmt.Sprintf("$%d\r\n%s\r\n", len(res.(string)), res)
				conn.Write([]byte(response))
			}
		}
	}
}
