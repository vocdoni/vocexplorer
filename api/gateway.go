package api

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"

	"context"
	"encoding/json"
	"fmt"
	"time"

	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"nhooyr.io/websocket"
)

// InitGateway initializes a connection with the gateway
func InitGateway(host string) (*GatewayClient, context.CancelFunc) {
	// Init Gateway client
	fmt.Printf("connecting to %s\n", host)
	gwClient, cancel, err := New(host)
	if err != nil {
		log.Error(err)
		for i := 0; i < 10; i++ {
			gwClient, cancel, err = New(host)
			if err != nil {
				log.Error(err)
			} else {
				break
			}
		}
	}
	if err != nil {
		return nil, cancel
	}
	return gwClient, cancel
}

// New starts a connection with the given endpoint address. From unreleased go-dvote/client
func New(addr string) (*GatewayClient, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	conn, _, err := websocket.Dial(ctx, addr, nil)
	if err != nil {
		return nil, cancel, err
	}
	return &GatewayClient{Addr: addr, Conn: conn, Ctx: ctx}, cancel, nil
}

// PingGateway pings the gateway host
func PingGateway(host string) bool {
	if strings.HasPrefix(host, "ws://") {
		host = host[5:]
	}
	pingClient := http.Client{
		Timeout: 5 * time.Second,
	}
	for i := 0; ; i++ {
		if i > 10 {
			return false
		}
		resp, err := pingClient.Get("http://" + host + "/ping")
		if err != nil {
			log.Debug(err.Error())
			time.Sleep(2 * time.Second)
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Debug(err.Error())
			time.Sleep(time.Second)
			continue
		}
		if string(body) != "pong" {
			log.Warn("Gateway ping not yet available")
		} else {
			return true
		}
	}
}

// GatewayClient holds an API websocket api.
type GatewayClient struct {
	Addr string
	Conn *websocket.Conn
	Ctx  context.Context
}

// GetEntityCount gets number of entities
func (c *GatewayClient) GetEntityCount() (int64, error) {
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
func (c *GatewayClient) GetProcessCount() (int64, error) {
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

// GetProcessKeys gets process keys
func (c *GatewayClient) GetProcessKeys(pid string) (*Pkeys, error) {
	var req MetaRequest
	req.Method = "getProcessKeys"
	req.ProcessID = pid
	// req.EntityID = eid
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get keys for process %s: (%s)", pid, resp.Message)
	}
	return &Pkeys{
		Pub:  resp.EncryptionPublicKeys,
		Priv: resp.EncryptionPrivKeys,
		Comm: resp.CommitmentKeys,
		Rev:  resp.RevealKeys}, nil
}

// GetGatewayInfo gets gateway info
func (c *GatewayClient) GetGatewayInfo() ([]string, int32, bool, error) {
	var req MetaRequest
	req.Method = "getGatewayInfo"
	req.Timestamp = int32(time.Now().Unix())

	resp, err := c.Request(req)
	if err != nil {
		return nil, 0, false, err
	}
	if !resp.Ok {
		return nil, 0, false, fmt.Errorf("cannot get gateway infos")
	}
	return resp.APIList, resp.Health, resp.Ok, nil
}

// GetBlockStatus gets latest block status for blockchain
func (c *GatewayClient) GetBlockStatus() (*[5]int32, int32, int64, error) {
	var req MetaRequest
	req.Method = "getBlockStatus"
	req.Timestamp = int32(time.Now().Unix())

	resp, err := c.Request(req)
	if err != nil {
		return nil, 0, 0, err
	}
	if !resp.Ok {
		return nil, 0, 0, fmt.Errorf("cannot get gateway infos")
	}
	return resp.BlockTime, resp.BlockTimestamp, *resp.Height, nil
}

// GetFinalProcessList gets list of finished processes on the Vochain
func (c *GatewayClient) GetFinalProcessList(from int64) ([]string, error) {
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
func (c *GatewayClient) GetLiveProcessList(from int64) ([]string, error) {
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
func (c *GatewayClient) GetScrutinizerEntities(from string) ([]string, error) {
	var req MetaRequest
	req.Method = "getScrutinizerEntities"
	req.Timestamp = int32(time.Now().Unix())
	req.FromID = from
	req.ListSize = 64

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
func (c *GatewayClient) GetProcessList(entityID string, from string) ([]string, error) {
	var req MetaRequest
	req.Method = "getProcessList"
	req.Timestamp = int32(time.Now().Unix())
	req.EntityID = entityID
	req.FromID = from
	req.ListSize = 64

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
func (c *GatewayClient) GetEnvelopeHeight(processID string) (int64, error) {
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
func (c *GatewayClient) GetEnvelopeList(processID string, from int64) ([]string, error) {
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
func (c *GatewayClient) GetProcessResults(processID string) (string, string, [][]uint32, error) {
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
func (c *GatewayClient) GetEnvelopeStatus(nullifier, processID string) (bool, error) {
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
func (c *GatewayClient) GetEnvelope(processID, nullifier string) (string, error) {
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
func (c *GatewayClient) Request(req MetaRequest) (*MetaResponse, error) {
	method := req.Method
	req.Timestamp = int32(time.Now().Unix())
	reqInner, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}

	reqOuter := RequestMessage{
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

	if err := c.Conn.Write(ctx, websocket.MessageText, reqBody); err != nil {
		return nil, fmt.Errorf("error: %s: %v", method, err)
	}
	// fmt.Println("sent request: " + req.Method)
	_, message, err := c.Conn.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}

	var respOuter ResponseMessage
	if err := json.Unmarshal(message, &respOuter); err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}
	for respOuter.ID != reqOuter.ID {
		fmt.Printf("%s: %v\n", method, "request ID doesn't match\n")
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
func (c *GatewayClient) Close() {
	err := c.Conn.Close(websocket.StatusNormalClosure, "")
	if err != nil {
		log.Error(err)
	} else {
		fmt.Println("Closed websocket connection")
	}
}
