package components

import (
	"fmt"

	"github.com/gopherjs/vecty"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/rpc"
)

// BlockchainVersion is a tiny component showing the blockchain we're connected to and its version
type BlockchainVersion struct {
	vecty.Core
	T *rpc.TendermintInfo
}

//Render renders the BlockchainVersion component
func (b *BlockchainVersion) Render() vecty.ComponentOrHTML {
	return &bootstrap.Alert{
		Contents: fmt.Sprintf(
			"Connected to blockchain \"<i>%s</i>\" (version %s)",
			b.T.Genesis.ChainID,
			b.T.ResultStatus.NodeInfo.Version,
		),
	}
}
