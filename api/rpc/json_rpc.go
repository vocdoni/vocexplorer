package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
	"gitlab.com/vocdoni/go-dvote/log"
	"nhooyr.io/websocket"
)

var id uint64
var cdc = amino.NewCodec()

func init() {
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PubKeyEd25519{},
		ed25519.PubKeyAminoName, nil)
}

// InitTendermintRPC initializes a TendermintRPC client
func InitTendermintRPC(host string, conns int) (*TendermintRPC, error) {
	t := new(TendermintRPC)
	tMutex := new(sync.Mutex)
	done := make(chan error, conns)
	for i := 0; i < conns; i++ {
		go newConnection(done, t, host, tMutex)
	}
	// Sync: wait here for all goroutines to complete
	num := 0
	for err := range done {
		if err != nil {
			log.Warn(err)
		}
		if num >= conns-1 {
			break
		}
		num++
	}
	if t == nil || len(t.Conns) < 1 {
		return nil, fmt.Errorf("cannot connect to websocket client")
	}
	result := new(coretypes.ResultStatus)
	_, err := t.call("status", nil, result)
	if err != nil {
		return nil, errors.Wrap(err, "status")
	}

	return t, nil
}

func newConnection(done chan error, t *TendermintRPC, host string, connMutex *sync.Mutex) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, _, err := websocket.Dial(ctx, host, &websocket.DialOptions{})
	if err != nil {
		done <- err
		return
	}
	// Set readLimit to the maximum read size, from tendermint/p2p/conn/connection.go
	c.SetReadLimit(22020096)
	connMutex.Lock()
	t.AddConnection(c)
	connMutex.Unlock()
	done <- nil
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
			return nil, err
		}
		paramsMap[name] = valueJSON
	}

	payload, err := json.Marshal(paramsMap) // NOTE: Amino doesn't handle maps yet.
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func bundleRequest(method string, params map[string]interface{}, id rpctypes.JSONRPCIntID) ([]byte, error) {
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
	myID := rpctypes.JSONRPCIntID(atomic.AddUint64(&id, 1))
	err := t.request(method, params, myID, result)
	return result, err
}

func (t *TendermintRPC) request(method string, params map[string]interface{}, myID rpctypes.JSONRPCIntID, result interface{}) error {
	req, err := bundleRequest(method, params, myID)
	if err != nil {
		return err
	}
	p := t.GetConnection()
	if p == nil {
		return errors.Errorf("unable to get websocket connection from pool")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	var msg []byte
	msg, err = p.WriteRead(ctx, req)
	p.Release()
	if err != nil {
		return err
	}
	_, err = UnmarshalResponseBytes(cdc, msg, myID, result)
	if err != nil {
		return err
	}
	return nil
}
