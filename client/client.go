package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"gitlab.com/vocdoni/vocexplorer/logger"
	"go.vocdoni.io/dvote/httprouter/jsonrpcapi"
	"nhooyr.io/websocket"
)

type Client struct {
	Address string
	ws      *websocket.Conn
	http    *http.Client
}

// New starts a connection with the given endpoint address.
// Supported protocols are ws(s):// and http(s)://
func New(addr string) (*Client, error) {
	cli := &Client{Address: addr}
	var err error
	if strings.HasPrefix(addr, "ws") {
		logger.Info(fmt.Sprintf("Connecting to gateway: %v", addr))
		cli.ws, _, err = websocket.Dial(context.Background(), addr, nil)
		if err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(addr, "http") {
		tr := &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    10 * time.Second,
			DisableCompression: false,
		}
		cli.http = &http.Client{Transport: tr, Timeout: time.Second * 20}
		if cli.http == nil {
			return nil, fmt.Errorf("unable to connect to %s", addr)
		}
	} else {
		return nil, fmt.Errorf("address is not websockets nor http: %s", addr)
	}
	return cli, nil
}

func (c *Client) Close() error {
	var err error
	if c.ws != nil {
		err = c.ws.Close(websocket.StatusNormalClosure, "")
	}
	if c.http != nil {
		c.http.CloseIdleConnections()
	}
	return err
}

// Request makes a request to the previously connected endpoint
func (c *Client) Request(req APIrequest) (*APIresponse, error) {
	if c == nil {
		return nil, fmt.Errorf("unable to make request %s: client not connected", req.Method)
	}
	method := req.Method
	req.Timestamp = int32(time.Now().Unix())
	reqInner, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}

	reqOuter := jsonrpcapi.RequestMessage{
		ID:         fmt.Sprintf("%d", rand.Intn(1000)),
		MessageAPI: reqInner,
	}
	reqBody, err := json.Marshal(reqOuter)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}

	message := []byte{}
	if c.ws != nil {
		tctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		if err := c.ws.Write(tctx, websocket.MessageText, reqBody); err != nil {
			return nil, fmt.Errorf("%s: %v", method, err)
		}
		_, message, err = c.ws.Read(tctx)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", method, err)
		}
	}
	if c.http != nil {
		resp, err := c.http.Post(c.Address, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, err
		}
		message, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		resp.Body.Close()
	}
	var respOuter jsonrpcapi.ResponseMessage
	if err := json.Unmarshal(message, &respOuter); err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}
	if respOuter.ID != reqOuter.ID {
		return nil, fmt.Errorf("%s: %v", method, "request ID doesn't match")
	}
	if len(respOuter.Signature) == 0 {
		return nil, fmt.Errorf("%s: empty signature in response: %s", method, message)
	}
	var respInner APIresponse
	if err := json.Unmarshal(respOuter.MessageAPI, &respInner); err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}
	return &respInner, nil
}
