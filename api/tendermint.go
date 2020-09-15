package api

import (
	"time"

	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/api/rpc"
)

//StartTendermint starts the tendermint client
func StartTendermint(host string, conns int) (*rpc.TendermintRPC, bool) {
	for i := 0; ; i++ {
		if i > 20 {
			return nil, false
		}
		hostCopy := string([]byte(host))
		tmClient := StartTendermintClient(hostCopy, conns)
		if tmClient == nil {
			time.Sleep(1 * time.Second)
			continue
		} else {
			return tmClient, true
		}
	}
}

// StartTendermintClient initializes an http tendermint api client on websockets
func StartTendermintClient(host string, conns int) *rpc.TendermintRPC {
	log.Infof("connecting to %s with %d connections", host, conns)
	tClient, err := rpc.InitTendermintRPC(host, conns)
	if err != nil {
		log.Warn(err.Error())
		return nil
	}
	return tClient
}

// PingTendermint pings the tendermint client and returns true if ok
func PingTendermint(t *rpc.TendermintRPC) bool {
	if t == nil {
		return false
	}
	status, err := t.Status()
	if err != nil || status == nil {
		return false
	}
	return true
}

// GetHealth calls the tendermint Health api
func GetHealth(t *rpc.TendermintRPC) *coretypes.ResultStatus {
	if t == nil {
		return nil
	}
	status, err := t.Status()
	if err != nil {
		log.Error(err)
		return nil
	}
	return status
}

// GetGenesis gets the first block
func GetGenesis(t *rpc.TendermintRPC) *types.GenesisDoc {
	if t == nil {
		return nil
	}
	result, err := t.Genesis()
	if err != nil {
		log.Error(err)
		return nil
	}
	return result.Genesis
}

// GetBlock returns the contents of one block
func GetBlock(t *rpc.TendermintRPC, height int64) *coretypes.ResultBlock {
	if t == nil {
		return nil
	}
	block, err := t.Block(&height)
	if err != nil {
		log.Error(err)
		return nil
	}
	return block
}

// GetTransaction gets a transaction by hash
func GetTransaction(t *rpc.TendermintRPC, hash []byte) *coretypes.ResultTx {
	if t == nil {
		return nil
	}
	res, err := t.Tx(hash, false)
	if err != nil {
		log.Error(err)
		return nil
	}
	return res
}
