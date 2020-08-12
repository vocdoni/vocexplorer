package client

import (
	"math/rand"

	"context"
	"encoding/json"
	"fmt"
	"time"

	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/util"
	"nhooyr.io/websocket"
)

// InitGateway initializes a connection with the gateway
func InitGateway(host string) (*Client, context.CancelFunc) {
	// Init Gateway client
	fmt.Println("connecting to " + host)
	gwClient, cancel, err := New(host)
	if util.ErrPrint(err) {
		return nil, cancel
	}
	return gwClient, cancel
}

// New starts a connection with the given endpoint address. From unreleased go-dvote/client
func New(addr string) (*Client, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	conn, _, err := websocket.Dial(ctx, addr, nil)
	if err != nil {
		return nil, cancel, err
	}
	return &Client{Addr: addr, Conn: conn, Ctx: ctx}, cancel, nil
}

// VochainInfo requests

// GetEntityCount gets number of entities
func (c *Client) GetEntityCount() (int64, error) {
	var req MetaRequest
	req.Method = "getScrutinizerEntityCount"
	req.Timestamp = int32(time.Now().Unix())

	resp, err := c.Request(req)
	if err != nil {
		return 0, err
	}
	if !resp.Ok {
		return 0, fmt.Errorf("cannot get entity count")
	}
	return *resp.Size, nil
}

// GetProcessCount gets number of processes
func (c *Client) GetProcessCount() (int64, error) {
	var req MetaRequest
	req.Method = "getProcessCount"
	req.Timestamp = int32(time.Now().Unix())

	resp, err := c.Request(req)
	if err != nil {
		return 0, err
	}
	if !resp.Ok {
		return 0, fmt.Errorf("cannot get process count")
	}
	return *resp.Size, nil
}

// GetGatewayInfo gets gateway info
func (c *Client) GetGatewayInfo() ([]string, int32, bool, int32, error) {
	var req MetaRequest
	req.Method = "getGatewayInfo"
	req.Timestamp = int32(time.Now().Unix())

	resp, err := c.Request(req)
	if err != nil {
		return nil, 0, false, 0, err
	}
	if !resp.Ok {
		return nil, 0, false, 0, fmt.Errorf("cannot get gateway infos")
	}
	return resp.APIList, resp.Health, resp.Ok, resp.Timestamp, nil
}

// GetBlockStatus gets latest block status for blockchain
func (c *Client) GetBlockStatus() (*[5]int32, int32, int64, bool, error) {
	var req MetaRequest
	req.Method = "getBlockStatus"
	req.Timestamp = int32(time.Now().Unix())

	resp, err := c.Request(req)
	if err != nil {
		return nil, 0, 0, false, err
	}
	if !resp.Ok {
		return nil, 0, 0, false, fmt.Errorf("cannot get gateway infos")
	}
	return resp.BlockTime, resp.BlockTimestamp, *resp.Height, resp.Ok, nil
}

// GetFinalProcessList gets list of finished processes on the Vochain
func (c *Client) GetFinalProcessList(from int64) ([]string, error) {
	var req MetaRequest
	req.Method = "getProcListResults"
	req.Timestamp = int32(time.Now().Unix())
	req.From = from
	req.ListSize = int64(config.ListSize)

	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get gateway infos")
	}
	return resp.ProcessIDs, nil
}

// GetLiveProcessList gets list of live processes on the Vochain
func (c *Client) GetLiveProcessList(from int64) ([]string, error) {
	var req MetaRequest
	req.Method = "getProcListLiveResults"
	req.Timestamp = int32(time.Now().Unix())
	req.From = from
	req.ListSize = int64(config.ListSize)

	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get gateway infos")
	}
	return resp.ProcessIDs, nil
}

// GetScrutinizerEntities gets list of entities indexed by the scrutinizer on the Vochain
func (c *Client) GetScrutinizerEntities(from int64) ([]string, error) {
	var req MetaRequest
	req.Method = "getScrutinizerEntities"
	req.Timestamp = int32(time.Now().Unix())
	req.From = from
	req.ListSize = int64(config.ListSize)

	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get gateway infos")
	}
	return resp.EntityIDs, nil
}

// EntityInfo requests

// GetProcessList gets list of processes for a given entity, starting at from
func (c *Client) GetProcessList(entityID string, from int64) ([]string, error) {
	var req MetaRequest
	req.Method = "getProcessList"
	req.Timestamp = int32(time.Now().Unix())
	req.EntityID = entityID
	req.From = from
	req.ListSize = int64(config.ListSize)

	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get gateway infos")
	}
	return resp.ProcessList, nil
}

// ProcessInfo requests

// GetEnvelopeHeight gets number of envelopes in a given process
func (c *Client) GetEnvelopeHeight(processID string) (int64, error) {
	var req MetaRequest
	req.Method = "getEnvelopeHeight"
	req.ProcessID = processID
	resp, err := c.Request(req)
	if err != nil {
		return 0, err
	}
	if !resp.Ok {
		return 0, fmt.Errorf(resp.Message)
	}
	return *resp.Height, nil
}

// GetEnvelopeList gets list of envelopes in a given process, starting at from
func (c *Client) GetEnvelopeList(processID string, from int64) ([]string, error) {
	var req MetaRequest
	req.Method = "getEnvelopeList"
	req.ProcessID = processID
	req.From = from
	req.ListSize = int64(config.ListSize)
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf(resp.Message)
	}
	return *resp.Nullifiers, nil
}

// GetProcessResults gets the results of a given process
func (c *Client) GetProcessResults(processID string) (string, string, [][]uint32, error) {
	var req MetaRequest
	req.Method = "getResults"
	req.ProcessID = processID
	resp, err := c.Request(req)
	if err != nil {
		return "", "", nil, err
	}
	if !resp.Ok {
		return "", "", nil, fmt.Errorf(resp.Message)
	}
	return resp.Type, resp.State, resp.Results, nil
}

// EnvelopeInfo requests

// GetEnvelopeStatus gets status of given envelope
func (c *Client) GetEnvelopeStatus(nullifier, processID string) (bool, error) {
	var req MetaRequest
	req.Method = "getEnvelopeStatus"
	req.ProcessID = processID
	req.Nullifier = nullifier
	resp, err := c.Request(req)
	if err != nil {
		return false, err
	}
	if !resp.Ok || resp.Registered == nil {
		return false, fmt.Errorf("cannot check envelope (%s)", resp.Message)
	}
	return *resp.Registered, nil
}

// GetEnvelope gets contents of given envelope
func (c *Client) GetEnvelope(processID, nullifier string) (string, error) {
	var req MetaRequest
	req.Method = "getEnvelope"
	req.Timestamp = int32(time.Now().Unix())
	req.ProcessID = processID
	req.Nullifier = nullifier

	resp, err := c.Request(req)
	if err != nil {
		return "", err
	}
	if !resp.Ok {
		return "", fmt.Errorf("cannot get envelope contents")
	}
	return resp.Payload, nil
}

//___________________________________________________________________________

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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// fmt.Println("request:", reqBody)
	if err := c.Conn.Write(ctx, websocket.MessageText, reqBody); err != nil {
		return nil, fmt.Errorf("Error: %s: %v", method, err)
	}
	fmt.Println("sent request: " + req.Method)
	_, message, err := c.Conn.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}

	var respOuter ResponseMessage
	if err := json.Unmarshal(message, &respOuter); err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}
	for respOuter.ID != reqOuter.ID {
		fmt.Printf("%s: %v", method, "request ID doesn't match\n")
		// Try to read & trash one more message so client can catch up
		_, message, err := c.Conn.Read(ctx)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", method, err)
		}
		if err := json.Unmarshal(message, &respOuter); err != nil {
			return nil, fmt.Errorf("%s: %v", method, err)
		}
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
