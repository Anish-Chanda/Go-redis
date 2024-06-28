package tests

import (
	"reflect"
	"testing"

	"github.com/anish-chanda/goredis/helpers"
)

func TestParseSimplePong(t *testing.T) {
	input := "+PONG\r\n"
	inputBytes := []byte(input)

	respData, err := helpers.RespParser(inputBytes)
	if err != nil {
		t.Fatalf("Failed to parse response: %s", err)
	}

	if !reflect.DeepEqual(respData.Command, "PONG") {
		t.Fatalf("Expected PONG, got %s", respData.Command)
	}
}

func TestArrayParse(t *testing.T) {

	input := "*2\r\n+hi\r\n+ho\r\n"
	inputBytes := []byte(input)

	respData, err := helpers.RespParser(inputBytes)
	if err != nil {
		t.Fatalf("Failed to parse response: %s", err)
	}

	if !reflect.DeepEqual(respData.Command, "hi") {
		t.Fatalf("Expected hi, got %s", respData.Command)
	}

	//test with bulk strings
	input = "*3\r\n$4\r\nECHO\r\n$3\r\nhey\r\n$4\r\nhey2\r\n"
	// fmt.Println("TESTING ECHO")
	inputBytes = []byte(input)

	respData, err = helpers.RespParser(inputBytes)
	if err != nil {
		t.Fatalf("Failed to parse response: %s", err)
	}

	if !reflect.DeepEqual(respData.Command, "ECHO") {
		t.Fatalf("Expected ECHO, got %s", respData.Command)
	}

	if !reflect.DeepEqual(respData.Args[0], "hey") {
		t.Fatalf("Expected hey, got %s", respData.Command)
	}
	if !reflect.DeepEqual(respData.Args[1], "hey2") {
		t.Fatalf("Expected hey, got %s", respData.Command)
	}
}

func TestBulkStringParse(t *testing.T) {
	input := "$5\r\nhello\r\n"
	inputBytes := []byte(input)

	respData, err := helpers.RespParser(inputBytes)
	if err != nil {
		t.Fatalf("Failed to parse response: %s", err)
	}

	if !reflect.DeepEqual(respData.Command, "hello") {
		t.Fatalf("Expected hello, got %s", respData.Command)
	}

	if !reflect.DeepEqual(respData.Args, []string{}) {
		t.Fatalf("Expected [], got %s", respData.Args)
	}
}

func TestEmptyBulkString(t *testing.T) {
	input := "$0\r\n\r\n"
	inputBytes := []byte(input)

	respData, err := helpers.RespParser(inputBytes)
	if err != nil {
		t.Fatalf("Failed to parse response: %s", err)
	}

	if !reflect.DeepEqual(respData.Command, "") {
		t.Fatalf("Expected empty, got %s", respData.Command)
	}

	if !reflect.DeepEqual(respData.Args, []string{}) {
		t.Fatalf("Expected [], got %s", respData.Args)
	}
}
