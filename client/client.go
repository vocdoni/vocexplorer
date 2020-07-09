package client

import (
	"context"
	"math/rand"

	"encoding/json"
	"fmt"
	"time"

	"gitlab.com/vocdoni/vocexplorer/util"
	"nhooyr.io/websocket"
)

// New starts a connection with the given endpoint address. From unreleased go-dvote/client
func New(addr string) (*Client, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	conn, _, err := websocket.Dial(ctx, addr, nil)
	if err != nil {
		return nil, cancel, err
	}
	return &Client{Addr: addr, Conn: conn, Ctx: ctx}, cancel, nil
}

// GetGatewayInfo calls gateway api and returns gateway info
func (c *Client) GetGatewayInfo() (*MetaResponse, error) {
	var req MetaRequest
	req.Method = "getGatewayInfo"
	req.Timestamp = int32(time.Now().Unix())

	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get gateway infos")
	}
	return resp, nil
}

// Request makes a request to the previously connected endpoint
func (c *Client) Request(req MetaRequest) (*MetaResponse, error) {
	method := req.Method
	req.Timestamp = int32(time.Now().Unix())
	reqInner, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}

	reqOuter := RequestMessage{
		// ID: "req-2345679",
		ID:          fmt.Sprintf("%d", rand.Intn(1000)),
		MetaRequest: reqInner,
	}
	reqBody, err := json.Marshal(reqOuter)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}

	// Set context for request
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	fmt.Println("request:", reqBody)
	if err := c.Conn.Write(ctx, websocket.MessageText, reqBody); err != nil {
		return nil, fmt.Errorf("Error: %s: %v", method, err)
	}
	fmt.Println("sent request")
	_, message, err := c.Conn.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}
	fmt.Println("response:", message)
	var respOuter ResponseMessage
	if err := json.Unmarshal(message, &respOuter); err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}
	if respOuter.ID != reqOuter.ID {
		return nil, fmt.Errorf("%s: %v", method, "request ID doesn't match")
	}
	if respOuter.Signature == "" {
		return nil, fmt.Errorf("%s: empty signature in response: %s", method, message)
	}
	var respInner MetaResponse
	if err := json.Unmarshal(respOuter.MetaResponse, &respInner); err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}
	return &respInner, nil
}

// Close closes given websocket connection
func (c *Client) Close() {
	err := c.Conn.Close(websocket.StatusNormalClosure, "")
	if !util.ErrPrint(err) {
		fmt.Println("Closed websocket connection")
	}
}
