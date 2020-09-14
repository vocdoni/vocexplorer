package api

import (
	"time"

	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/api/rpc"
	"nhooyr.io/websocket"
)

//StartTendermint starts the tendermint client
func StartTendermint(host string) (*websocket.Conn, bool) {
	for i := 0; ; i++ {
		if i > 20 {
			return nil, false
		}
		hostCopy := string([]byte(host))
		tmClient := StartTendermintClient(hostCopy)
		if tmClient == nil {
			time.Sleep(1 * time.Second)
			continue
		} else {
			return tmClient, true
		}
	}
}

// StartTendermintClient initializes an http tendermint api client on websockets
func StartTendermintClient(host string) *websocket.Conn {
	log.Infof("connecting to %s", host)
	tClient, err := rpc.NewClient(host)
	if err != nil {
		log.Warn(err.Error())
		return nil
	}
	return tClient
}

// PingTendermint pings the tendermint client and returns true if ok
func PingTendermint(c *websocket.Conn) bool {
	if c == nil {
		return false
	}
	status, err := rpc.Status(c)
	if err != nil || status == nil {
		return false
	}
	return true
}

// GetHealth calls the tendermint Health api
func GetHealth(c *websocket.Conn) *coretypes.ResultStatus {
	if c == nil {
		return nil
	}
	status, err := rpc.Status(c)
	if err != nil {
		log.Error(err)
		return nil
	}
	return status
}

// GetGenesis gets the first block
func GetGenesis(c *websocket.Conn) *types.GenesisDoc {
	if c == nil {
		return nil
	}
	result, err := rpc.Genesis(c)
	if err != nil {
		log.Error(err)
		return nil
	}
	return result.Genesis
}

// GetBlock returns the contents of one block
func GetBlock(c *websocket.Conn, height int64) *coretypes.ResultBlock {
	if c == nil {
		return nil
	}
	block, err := rpc.Block(c, &height)
	if err != nil {
		log.Error(err)
		return nil
	}
	return block
}

// GetTransaction gets a transaction by hash
func GetTransaction(c *websocket.Conn, hash []byte) *coretypes.ResultTx {
	if c == nil {
		return nil
	}
	res, err := rpc.Tx(c, hash, false)
	if err != nil {
		log.Error(err)
		return nil
	}
	return res
}
