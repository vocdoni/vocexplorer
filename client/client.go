package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"go.vocdoni.io/dvote/log"
	"go.vocdoni.io/dvote/types"
)

type Client struct {
	Address string
	http    *http.Client
}

func New(address string) (*Client, error) {
	cli := &Client{Address: address}
	if strings.HasPrefix(address, "http") {
		tr := &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    10 * time.Second,
			DisableCompression: false,
		}
		cli.http = &http.Client{Transport: tr, Timeout: time.Second * 5}
	} else {
		return nil, fmt.Errorf("address is not http: %s", address)
	}
	return cli, nil
}

func (c *Client) Close() {
	if c.http != nil {
		c.http.CloseIdleConnections()
	}
}

// Request makes a request to the previously connected endpoint
func (c *Client) Request(req types.MetaRequest) (*types.MetaResponse, error) {
	method := req.Method
	req.Timestamp = int32(time.Now().Unix())
	reqInner, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}

	reqOuter := types.RequestMessage{
		ID:          fmt.Sprintf("%d", rand.Intn(1000)),
		MetaRequest: reqInner,
	}
	reqBody, err := json.Marshal(reqOuter)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}

	log.Debugf("request: %s", reqBody)

	resp, err := c.http.Post(c.Address, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	message, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	log.Debugf("response: %s", message)
	var respOuter types.ResponseMessage
	if err := json.Unmarshal(message, &respOuter); err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}
	if respOuter.ID != reqOuter.ID {
		return nil, fmt.Errorf("%s: %v", method, "request ID doesn't match")
	}
	if len(respOuter.Signature) == 0 {
		return nil, fmt.Errorf("%s: empty signature in response: %s", method, message)
	}
	var respInner types.MetaResponse
	if err := json.Unmarshal(respOuter.MetaResponse, &respInner); err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}
	return &respInner, nil
}
