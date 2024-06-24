package main

import (
	"context"
	"fmt"
	"net"

	"github.com/tidwall/resp"
)

func main() {
	ctx := context.Background()

	client := NewClient(":5000")

	for i := range 10 {
		client.Set(
			ctx,
			fmt.Sprintf("key_%d", i),
			[]byte(fmt.Sprintf("val_%d", i)),
		)
	}

	for i := range 10 {
		key := fmt.Sprintf("key_%d", i)
		val, _ := client.Get(ctx, key)
		fmt.Printf("'%s' => '%s'\n", key, val)
	}
}

type Client struct {
	url string
}

func NewClient(url string) *Client {
	return &Client{url: url}
}

func (c *Client) Set(ctx context.Context, key string, val []byte) error {
	conn, err := net.Dial("tcp", c.url)
	if err != nil {
		return err
	}
	defer conn.Close()

	wr := resp.NewWriter(conn)
	if err := wr.WriteArray([]resp.Value{
		resp.StringValue("SET"),
		resp.StringValue(key),
		resp.BytesValue(val),
	}); err != nil {
		return err
	}

	return nil
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	conn, err := net.Dial("tcp", c.url)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	wr := resp.NewWriter(conn)
	if err := wr.WriteArray([]resp.Value{
		resp.StringValue("GET"),
		resp.StringValue(key),
	}); err != nil {
		return nil, err
	}

	p := make([]byte, 1024)
	n, err := conn.Read(p)

	return p[:n], nil
}
