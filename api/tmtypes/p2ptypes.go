package tmtypes

import (
	"github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/version"
)

// ProtocolVersion contains the protocol versions for the software.
type ProtocolVersion struct {
	P2P   version.Protocol `json:"p2p"`
	Block version.Protocol `json:"block"`
	App   version.Protocol `json:"app"`
}

// DefaultNodeInfo is the basic node information exchanged
// between two peers during the Tendermint P2P handshake.
type DefaultNodeInfo struct {
	ProtocolVersion ProtocolVersion `json:"protocol_version"`

	// Authenticate
	// TODO: replace with NetAddress
	DefaultNodeID ID     `json:"id"`          // authenticated identifier
	ListenAddr    string `json:"listen_addr"` // accepting incoming

	// Check compatibility.
	// Channels are HexBytes so easier to read as JSON
	Network  string         `json:"network"`  // network/chain ID
	Version  string         `json:"version"`  // major.minor.revision
	Channels bytes.HexBytes `json:"channels"` // channels this node knows about

	// ASCIIText fields
	Moniker string               `json:"moniker"` // arbitrary moniker
	Other   DefaultNodeInfoOther `json:"other"`   // other application specific data
}

// DefaultNodeInfoOther is the misc. applcation specific data
type DefaultNodeInfoOther struct {
	TxIndex    string `json:"tx_index"`
	RPCAddress string `json:"rpc_address"`
}

// ID returns the node's peer ID.
func (info DefaultNodeInfo) ID() ID {
	return info.DefaultNodeID
}

// ID is a hex-encoded crypto.Address
type ID string
