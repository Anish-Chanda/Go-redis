package helpers

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/anish-chanda/goredis/types"
)

func RespParser(buffer []byte) (data types.RespData, err error) {
	//check the data type
	switch buffer[0] {
	case '+': //simple string
		return types.RespData{
			Command: parseSimpleString(buffer),
			Args:    nil,
		}, nil

	case '-': //simple error, essentially the same as a simple string
		return parseError(buffer), nil

	case '$':
		return parseBulkString(buffer)

	case '*': //arrays
		return parseArray(buffer)
	default:
		fmt.Println("def")
	}
	return types.RespData{}, nil
}

func parseBulkString(buffer []byte) (types.RespData, error) {
	lengthBytes := bytes.Index(buffer, []byte("\r\n"))
	if lengthBytes == -1 {
		return types.RespData{}, fmt.Errorf("malformed bulk string")
	}

	//calc length
	len, err := strconv.Atoi(string(buffer[1:lengthBytes]))
	if err != nil {
		return types.RespData{}, err
	}
	//TODO: if the len is -1 then its a null bulk string
	startOfData := lengthBytes + 2
	endOfData := startOfData + len

	//TODO:check if buffer is long enough

	data := string(buffer[startOfData:endOfData])

	return types.RespData{
		Command: data,
		Args:    []string{},
	}, nil
}

func parseArrayElement(buffer []byte) (res string, rem []byte, err error) {
	//TODO: handle other types
	switch buffer[0] {
	case '+': //simple string
		nextClrfIn := bytes.Index(buffer, []byte("\r\n"))
		return string(buffer[1:nextClrfIn]), buffer[nextClrfIn+2:], nil

	case '$':
		nextClrfIn := bytes.Index(buffer, []byte("\r\n"))
		dataLen, err := strconv.Atoi(string(buffer[1:nextClrfIn]))
		if err != nil {
			return "", []byte{}, err
		}

		buffer = buffer[nextClrfIn+2:]
		//fmt.Println("next bulk string, ", )
		fmt.Println("remaining buffer, ", string(buffer))
		return string(buffer[:dataLen]), buffer[dataLen+2:], nil
		//return parseBulkString(buffer)

	// case '*': //arrays
	// 	parseArray(buffer)
	default:
		fmt.Println("default")
		return "", []byte{}, nil
	}

}

func parseArray(buffer []byte) (types.RespData, error) {

	// calc number of elements
	length, err := strconv.Atoi(string(buffer[1:bytes.Index(buffer, []byte("\r\n"))]))
	buffer = buffer[4:]
	fmt.Println("length is, ", length)
	if err != nil {
		return types.RespData{}, fmt.Errorf("malformed array: invalid length")
	}

	fmt.Println("next is,", string(buffer))

	elements := make([]string, 0, length)

	for i := 0; i < length; i++ {
		res, rem, err := parseArrayElement(buffer)
		buffer = rem
		if err != nil {
			return types.RespData{}, err
		}
		elements = append(elements, res)

		fmt.Println("item, ", res)
		fmt.Println("remaining buffer, ", string(rem))
	}

	return types.RespData{Command: elements[0], Args: elements[1:]}, nil
}

func parseSimpleString(buffer []byte) string {
	//find index of the terminator
	i := bytes.Index(buffer, []byte("\r\n"))
	return string(buffer[1:i])
}

func parseError(buffer []byte) types.RespData {
	//find index of the terminator
	i := bytes.Index(buffer, []byte("\r\n"))
	msg := string(buffer[1:i])
	return types.RespData{
		Command: "ERROR",
		Args:    []string{msg},
	}
}
