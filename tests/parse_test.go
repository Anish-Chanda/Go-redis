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
