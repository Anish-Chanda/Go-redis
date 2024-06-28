package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/anish-chanda/goredis/helpers"
	"github.com/anish-chanda/goredis/types"
)

var tempStore map[string]types.Store = make(map[string]types.Store)

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
			var exp int64
			//check if exp was set
			if len(data.Args) > 3 && data.Args[2] == "px" {
				if len(data.Args) < 4 {
					//invalid length
				}
				exp, _ = strconv.ParseInt(data.Args[3], 10, 64)

				//TODO: handle errors

				fmt.Println("SETTING EXPIRY CURRENT TIME, ", time.Now().String())
				fmt.Println("SETTING EXPIRY, ", time.Now().Add(time.Duration(exp)*time.Millisecond).String())
				//the key will be arg[0] and val will be arg[1]
				tempStore[data.Args[0]] = types.Store{
					Value: data.Args[1],
					Exp:   time.Now().Add(time.Millisecond * time.Duration(exp)),
				}
				conn.Write([]byte("+OK\r\n"))
			} else {
				tempStore[data.Args[0]] = types.Store{
					Value: data.Args[1],
					Exp:   time.Time{},
				}
				conn.Write([]byte("+OK\r\n"))
			}

		} else if data.Command == "GET" {
			res, ok := tempStore[data.Args[0]]
			if !ok {
				conn.Write([]byte("$-1\r\n"))
			} else {
				//check if exp is in past
				fmt.Println("EXPIRY, ", res.Exp.String())
				fmt.Println("TIME NOW, ", time.Now().String())
				fmt.Println("VAL", res.Value)
				if res.Exp != (time.Time{}) && res.Exp.Before(time.Now()) {
					conn.Write([]byte("$-1\r\n"))
				} else {
					response := fmt.Sprintf("$%d\r\n%s\r\n", len(res.Value), res.Value)
					conn.Write([]byte(response))
				}

			}
		}
	}
}
