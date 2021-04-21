package components

import (
	"fmt"

	"github.com/hexops/vecty"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

// BlockchainVersion is a tiny component showing the blockchain we're connected to and its version
type BlockchainVersion struct {
	vecty.Core
}

//Render renders the BlockchainVersion component
func (b *BlockchainVersion) Render() vecty.ComponentOrHTML {
	return &bootstrap.Alert{
		Contents: fmt.Sprintf(
			"Connected to blockchain \"<i>%s</i>\"",
			store.Stats.ChainID,
		),
	}
}
