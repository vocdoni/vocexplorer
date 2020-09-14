package rpc

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
	"nhooyr.io/websocket"
)

var id rpctypes.JSONRPCIntID
var cdc *amino.Codec
var reqMutex *sync.Mutex

func init() {
	initCdc()
}

// NewClient initializes a jsonrpc client
func NewClient(host string) (*websocket.Conn, error) {
	reqMutex = new(sync.Mutex)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, _, err := websocket.Dial(ctx, host, &websocket.DialOptions{})
	if err != nil {
		return nil, err
	}
	// Set readLimit to the maximum read size, from tendermint/p2p/conn/connection.go
	c.SetReadLimit(22020096)
	result := new(coretypes.ResultStatus)
	_, err = call(c, "status", nil, result)
	if err != nil {
		return nil, errors.Wrap(err, "status")
	}

	return c, nil
}

// Status calls the tendermint status api
func Status(c *websocket.Conn) (*coretypes.ResultStatus, error) {
	result := new(coretypes.ResultStatus)
	_, err := call(c, "status", nil, result)
	if err != nil {
		return nil, errors.Wrap(err, "status")
	}
	return result, nil
}

// Genesis calls the tendermint Genesis api
func Genesis(c *websocket.Conn) (*coretypes.ResultGenesis, error) {
	result := new(coretypes.ResultGenesis)
	_, err := call(c, "genesis", nil, result)
	if err != nil {
		return nil, errors.Wrap(err, "genesis")
	}
	return result, nil
}

// Block calls the tendermint Block api
func Block(c *websocket.Conn, height *int64) (*coretypes.ResultBlock, error) {
	result := new(coretypes.ResultBlock)
	params := map[string]interface{}{
		"height": height,
	}
	_, err := call(c, "block", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "block")
	}
	return result, nil
}

// Tx calls the tendermint Tx api
func Tx(c *websocket.Conn, hash []byte, prove bool) (*coretypes.ResultTx, error) {
	result := new(coretypes.ResultTx)
	params := map[string]interface{}{
		"hash":  hash,
		"prove": prove,
	}
	_, err := call(c, "tx", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "tx")
	}
	return result, nil
}

// Validators calls the tendermint Validators api
func Validators(c *websocket.Conn, height *int64, page, perPage int) (*coretypes.ResultValidators, error) {
	result := new(coretypes.ResultValidators)
	params := map[string]interface{}{
		"height":   height,
		"page":     page,
		"per_page": perPage,
	}
	_, err := call(c, "validators", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "validators")
	}
	return result, nil
}

func marshalParams(params map[string]interface{}) ([]byte, error) {
	var paramsMap = make(map[string]json.RawMessage, len(params))
	for name, value := range params {
		valueJSON, err := cdc.MarshalJSON(value)
		if err != nil {
			return []byte{}, err
		}
		paramsMap[name] = valueJSON
	}

	payload, err := json.Marshal(paramsMap) // NOTE: Amino doesn't handle maps yet.
	if err != nil {
		return []byte{}, err
	}
	return payload, nil
}

func bundleRequest(method string, params map[string]interface{}, id rpctypes.JSONRPCIntID) ([]byte, error) {
	if params == nil {
		params = map[string]interface{}{}
	}
	payload, err := marshalParams(params)
	if err != nil {
		return nil, err
	}
	request := rpctypes.RPCRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  payload,
	}
	req, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func call(c *websocket.Conn, method string, params map[string]interface{}, result interface{}) (interface{}, error) {
	id++
	myID := id
	req, err := bundleRequest(method, params, myID)
	if err != nil {
		return nil, err
	}
	reqMutex.Lock()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = c.Write(ctx, websocket.MessageText, req)
	if err != nil {
		return nil, err
	}
	_, msg, err := c.Read(ctx)
	reqMutex.Unlock()
	if err != nil {
		return nil, err
	}
	response, err := UnmarshalResponseBytes(cdc, msg, myID, result)
	if err != nil {
		return nil, err
	}
	return response, err
}

func initCdc() {
	cdc = amino.NewCodec()
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PubKeyEd25519{},
		ed25519.PubKeyAminoName, nil)
}
