package main

import (
	"bytes"
	"testing"
)

func TestProtocol(t *testing.T) {
	msg := "*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$3\r\nbar\r\n"
	cmd, err := parseCommand(bytes.NewBufferString(msg))
	if err != nil {
		t.Fatal(err)
	}

	setCmd, ok := cmd.(*SetCommand)
	if !ok {
		t.Fatalf("command = '%s'; want 'SET'", "idk")
	}

	if setCmd.Key != "mykey" {
		t.Errorf("key = '%s'; want 'mykey'", setCmd.Key)
	}
	if bytes.Equal(setCmd.Val, []byte("bar")) {
		t.Errorf("value = '%s'; want 'bar'", setCmd.Key)
	}
}
