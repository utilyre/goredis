package main

import (
	"fmt"
	"io"

	"github.com/tidwall/resp"
)

type Command interface{}

type SetCommand struct {
	Key string
	Val []byte
}

type GetCommand struct {
	Key string
}

func parseCommand(r io.Reader) (Command, error) {
	rd := resp.NewReader(r)

	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if v.Type() == resp.Array {
			s := v.Array()
			switch s[0].String() {
			case "SET":
				if len(s) != 3 {
					return nil, fmt.Errorf("command 'SET': invalid number of parameters")
				}

				return &SetCommand{
					Key: s[1].String(),
					Val: s[2].Bytes(),
				}, nil
			case "GET":
				if len(s) != 2 {
					return nil, fmt.Errorf("command 'GET': invalid number of parameters")
				}

				return &GetCommand{Key: s[1].String()}, nil
			}
		}
	}

	return nil, fmt.Errorf("invalid or unknown command")
}
