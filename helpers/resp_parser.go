package helpers

import (
	"bytes"
	"fmt"

	"github.com/anish-chanda/goredis/types"
)

func RespParser(buffer []byte) (data types.RespData, err error) {
	// *2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n ECHO hey
	// *1\r\n$4\r\nPING\r\n PING

	//check the data type
	switch buffer[0] {
	case '+':
		fmt.Println("single string")
		return types.RespData{
			Command: parseSimpleString(buffer),
			Args:    nil,
		}, nil
	default:
		fmt.Println("def")
	}
	return types.RespData{}, nil
}

func parseSimpleString(buffer []byte) string {
	//find index of the terminator
	i := bytes.Index(buffer, []byte("\r\n"))
	return string(buffer[1:i])
}
