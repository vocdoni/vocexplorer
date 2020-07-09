package components

import (
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	// "gitlab.com/vocdoni/vocexplorer/rpc"
)

// BlocksView renders the blocks page
type BlocksView struct {
	vecty.Core
	// t *rpc.TendermintInfo
}

// Render renders the BlocksView component
func (b *BlocksView) Render() vecty.ComponentOrHTML {
	js.Global().Set("page", "blocks")
	js.Global().Set("gateway", false)
	return elem.Div(
		&Header{currentPage: "blocks"},
	)
}

// func (t *TendermintInfo) RenderTendermintInfo() vecty.ComponentOrHTML {
// 	if t != nil {
// 		return elem.Div{
// 			t.renderStatus()
// 		}
// 	}
// }

// func renderStatus() vecty.ComponentOrHTML{
// 	if t.Status != nil {
// 		sync := t.Status.SyncInfo
// 		valid := t.Status.ValidatorInfo
// 		return {
// 			vecty.If(sync != nil, elem.UnorderedList{
// 				elem.ListItem(vecty.Text("Latest Block Hash: "+sync.LatestBlockHash.Dump()))
// 				elem.ListItem(vecty.Text("Latest App Hash: "+sync.LatestAppHash.Dump()))
// 				elem.ListItem(vecty.Text("Latest App Hash: "+strconv.Itoa(sync.LatestBlockHeight))
// 				elem.ListItem(vecty.Text("Latest App Hash: "+sync.LatestBlockHeight.String())
// 			})
// 		}
// 	}
// }
