package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/anish-chanda/goredis/helpers"
	"github.com/anish-chanda/goredis/types"
)

var tempStore map[string]types.Store = make(map[string]types.Store)

var (
	port               *int
	replicaOf          *string
	role               string
	master_replid      string
	master_repl_offset int
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func main() {

	//read port argument
	port = flag.Int("port", 6379, "port where the server will listen")
	//read replica of argument
	replicaOf = flag.String("replicaof", "", "replica")

	flag.Parse()

	//find server role
	if *replicaOf == "" {
		role = "master"
	} else {
		//initial handshake with master
		role = "slave"
		handleHandshake()
	}

	fmt.Println("Role: ", role)

	//set master_replid
	master_replid = genAlphaNumString()
	fmt.Println("Master Replication ID: ", master_replid)

	//initial master_repl_offset
	master_repl_offset = 0

	//listen to tcp reqs on 6379
	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(*port))
	if err != nil {
		fmt.Println("Failed to bind to port " + strconv.Itoa(*port))
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on port " + strconv.Itoa(*port))

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

		//TODO: change to a switch statement
		if data.Command == "PING" {
			conn.Write([]byte("+PONG\r\n"))
		} else if data.Command == "ECHO" {
			response := fmt.Sprintf("$%d\r\n%s\r\n", len(data.Args[0]), data.Args[0])
			conn.Write([]byte(response))
		} else if data.Command == "SET" {
			handleSET(conn, data)
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

			//INFO command with only replication
		} else if data.Command == "INFO" {
			handleINFO(conn, data)
		}
	}
}

func handleSET(conn net.Conn, data types.RespData) {
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
}

func handleHandshake() {
	//ping master
	fmt.Println("Initiating Handshake with master")
	addr := strings.Replace(*replicaOf, " ", ":", 1)
	fmt.Println("ping addr: ", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Error connecting to master: ", err.Error())
		os.Exit(1)
	}
	conn.Write([]byte("*1\r\n$4\r\nPING\r\n"))
	fmt.Println("Pinged Master")

	//check if ping was successfull
	buffer := make([]byte, 256)
	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}
	data, err := helpers.RespParser(buffer)
	if err != nil {
		fmt.Println("Error parsing ping resp: ", err.Error())
		os.Exit(1)
	}
	//continue if ping was successfull
	if data.Command == "PONG" {
		// send port in replconf
		conn.Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$14\r\nlistening-port\r\n$4\r\n6380\r\n"))

		data := readConn(conn)
		//continue if replconf 1 was successfull
		if data.Command == "OK" {
			// send capabilities
			conn.Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$4\r\ncapa\r\n$6\r\npsync2\r\n"))
			data = readConn(conn)

			if data.Command != "OK" {
				fmt.Println("Error in handshake with master")
				os.Exit(1)
			} else {
				fmt.Println("Handshake (replconf sharing) with master successfull, continuing...")

				//send PSYNC
				//ask for replication id and send offset
				conn.Write([]byte("*3\r\n$5\r\nPSYNC\r\n$1\r\n?\r\n$2\r\n-1\r\n"))

				data = readConn(conn)
				fmt.Println("PSYNC respond", data.Command)
			}
		}

	}

}

func handleINFO(conn net.Conn, data types.RespData) {
	if len(data.Args) < 1 {
		conn.Write([]byte("$-1\r\n"))
		return
	}
	if data.Args[0] != "replication" {
		conn.Write([]byte("$-1\r\n"))
		return
	}
	res := `# Replication
	role:` + role + `
	master_replid:` + master_replid + `
	master_repl_offset:` + strconv.Itoa(master_repl_offset)

	//check what role this server has
	response := fmt.Sprintf("$%d\r\n%s\r\n", len(res), res)
	conn.Write([]byte(response))
}

func readConn(conn net.Conn) types.RespData {
	buffer := make([]byte, 256)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		return types.RespData{}
	}

	data, err := helpers.RespParser(buffer)
	if err != nil {
		fmt.Println("Error parsing command: ", err.Error())
		return types.RespData{}
	}
	return data
}

func genAlphaNumString() string {
	b := make([]byte, 40)
	for i := 0; i < 40; i++ {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
