package components

import (
	"fmt"
	"strconv"
	"syscall/js"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/tendermint/tendermint/rpc/client/http"
	"gitlab.com/vocdoni/vocexplorer/rpc"
)

// Stats renders the stats pane
type StatsView struct {
	vecty.Core
	t  *rpc.TendermintInfo
	gw *client.GatewayInfo
	c  *http.HTTP
}

// Render renders the StatsView component
func (b *StatsView) Render() vecty.ComponentOrHTML {
	return elem.Div(
		renderTendermintInfo(b.t),
	)
}

func renderTendermintInfo(t *rpc.TendermintInfo) vecty.ComponentOrHTML {
	if t != nil {
		return elem.Div(
			renderStatus(t),
		)
	}
	return elem.Div(vecty.Text("No info struct rendered"))
}

func renderStatus(t *rpc.TendermintInfo) vecty.ComponentOrHTML {
	if t.Status != nil {
		sync := t.Status.SyncInfo
		// valid := t.Status.ValidatorInfo
		return elem.Div(
			// vecty.If(sync != nil, elem.UnorderedList{
			elem.UnorderedList(
				elem.ListItem(vecty.Text("Latest Block Hash: "+sync.LatestBlockHash.String())),
				elem.ListItem(vecty.Text("Latest App Hash: "+sync.LatestAppHash.String())),
				elem.ListItem(vecty.Text("Latest Block Height: "+strconv.Itoa(int(sync.LatestBlockHeight)))),
				elem.ListItem(vecty.Text("Latest Block Time: "+sync.LatestBlockTime.String())),
			),
		// ),
		)
	}
	return elem.Div(vecty.Text("Waiting for tendermint blockchain info..."))
}

func initStatsView(t *rpc.TendermintInfo) *StatsView {
	js.Global().Set("tendermint", true)
	c, err := rpc.InitClient()
	if err != nil {
		js.Global().Get("alert").Invoke("Unable to connect to Tendermint client. Please see readme file")
		return nil
	}
	// var t *rpc.TendermintInfo
	var StatsView StatsView
	StatsView.c = c
	StatsView.t = t
	go updateAndRenderStats(&StatsView)
	return &StatsView
}

func updateAndRenderStats(bv *StatsView) {
	for js.Global().Get("tendermint").Bool() {
		fmt.Println("Getting tendermint info")
		rpc.UpdateBlockInfo(bv.c, bv.t)
		time.Sleep(5 * time.Second)
		vecty.Rerender(bv)
	}
	fmt.Println("Closing tendermint updater")
}
