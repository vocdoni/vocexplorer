package client

import (
	"context"

	"nhooyr.io/websocket"
)

// Client holds an API websocket client. From unreleased go-dvote/client
type Client struct {
	Addr string
	Conn *websocket.Conn
}

// New starts a connection with the given endpoint address. From unreleased go-dvote/client
func New(addr string) (*Client, error) {
	conn, _, err := websocket.Dial(context.TODO(), addr, nil)
	if err != nil {
		return nil, err
	}
	return &Client{Addr: addr, Conn: conn}, nil
}
