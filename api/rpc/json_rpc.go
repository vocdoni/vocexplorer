package rpc

import (
	"context"
	"encoding/json"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
	"nhooyr.io/websocket"
)

var id int32
var cdc *amino.Codec

func init() {
	initCdc()
}

// InitTendermintRPC initializes a TendermintRPC client
func InitTendermintRPC(host string, conns int) (*TendermintRPC, error) {
	t := new(TendermintRPC)
	for i := 0; i < conns; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		c, _, err := websocket.Dial(ctx, host, &websocket.DialOptions{})
		if err != nil {
			return nil, err
		}
		// Set readLimit to the maximum read size, from tendermint/p2p/conn/connection.go
		c.SetReadLimit(22020096)
		t.AddConnection(c)
	}
	result := new(coretypes.ResultStatus)
	_, err := t.call("status", nil, result)
	if err != nil {
		return nil, errors.Wrap(err, "status")
	}

	return t, nil
}

// Status calls the tendermint status api
func (t *TendermintRPC) Status() (*coretypes.ResultStatus, error) {
	result := new(coretypes.ResultStatus)
	_, err := t.call("status", nil, result)
	if err != nil {
		return nil, errors.Wrap(err, "status")
	}
	return result, nil
}

// Genesis calls the tendermint Genesis api
func (t *TendermintRPC) Genesis() (*coretypes.ResultGenesis, error) {
	result := new(coretypes.ResultGenesis)
	_, err := t.call("genesis", nil, result)
	if err != nil {
		return nil, errors.Wrap(err, "genesis")
	}
	return result, nil
}

// Block calls the tendermint Block api
func (t *TendermintRPC) Block(height *int64) (*coretypes.ResultBlock, error) {
	result := new(coretypes.ResultBlock)
	params := map[string]interface{}{
		"height": height,
	}
	_, err := t.call("block", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "block")
	}
	return result, nil
}

// Tx calls the tendermint Tx api
func (t *TendermintRPC) Tx(hash []byte, prove bool) (*coretypes.ResultTx, error) {
	result := new(coretypes.ResultTx)
	params := map[string]interface{}{
		"hash":  hash,
		"prove": prove,
	}
	_, err := t.call("tx", params, result)
	if err != nil {
		return nil, errors.Wrap(err, "tx")
	}
	return result, nil
}

// Validators calls the tendermint Validators api
func (t *TendermintRPC) Validators(height *int64, page, perPage int) (*coretypes.ResultValidators, error) {
	result := new(coretypes.ResultValidators)
	params := map[string]interface{}{
		"height":   height,
		"page":     page,
		"per_page": perPage,
	}
	_, err := t.call("validators", params, result)
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

func (t *TendermintRPC) call(method string, params map[string]interface{}, result interface{}) (interface{}, error) {
	myID := rpctypes.JSONRPCIntID(atomic.AddInt32(&id, 1))
	var err error
	done := make(chan struct{})
	go t.request(method, params, myID, result, &err, done)
	<-done
	return result, err
}

func (t *TendermintRPC) request(method string, params map[string]interface{}, myID rpctypes.JSONRPCIntID, response interface{}, err *error, done chan struct{}) {
	var req []byte
	req, *err = bundleRequest(method, params, myID)
	if *err != nil {
		close(done)
		return
	}
	p := t.GetConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	var msg []byte
	msg, *err = p.WriteRead(ctx, req)
	p.Release()
	if *err != nil {
		close(done)
		return
	}
	response, *err = UnmarshalResponseBytes(cdc, msg, myID, response)
	if *err != nil {
		close(done)
		return
	}
	close(done)
}

func initCdc() {
	cdc = amino.NewCodec()
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PubKeyEd25519{},
		ed25519.PubKeyAminoName, nil)
}
