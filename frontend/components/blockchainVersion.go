package components

import (
	"fmt"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/rpc"
)

// BlockchainVersion is a tiny component showing the blockchain we're connected to and its version
type BlockchainVersion struct {
	vecty.Core
	T *rpc.TendermintInfo
}

func (b *BlockchainVersion) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class("alert", "alert-success")),
		vecty.Markup(
			vecty.UnsafeHTML(
				fmt.Sprintf(
					"Connected to blockchain \"<i>%s</i>\" (version %s)",
					b.T.Genesis.ChainID,
					b.T.ResultStatus.NodeInfo.Version,
				),
			),
		),
	)
}
