package components

import (
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/tendermint/tendermint/rpc/client/http"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/dbapi"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// BlockTxsDashboardView renders the dashboard landing page
type BlockTxsDashboardView struct {
	vecty.Core
	t             *rpc.TendermintInfo
	tClient       *http.HTTP
	quitCh        chan struct{}
	refreshCh     chan int
	blockIndex    int
	disableUpdate bool
}

// Render renders the BlockTxsDashboardView component
func (dash *BlockTxsDashboardView) Render() vecty.ComponentOrHTML {
	if dash != nil && dash.tClient != nil && dash.t != nil && dash.t.ResultStatus != nil {
		return elem.Main(
			elem.Div(
				elem.Div(vecty.Markup(vecty.Class("card-col-3")),
					vecty.Text("Current Block Height: "+util.IntToString(dash.t.ResultStatus.SyncInfo.LatestBlockHeight)),
				),
				vecty.If(int(dash.t.ResultStatus.SyncInfo.LatestBlockHeight)-dash.t.TotalBlocks > 200,
					elem.Div(vecty.Markup(vecty.Class("card-col-3")),
						vecty.Text("Still Syncing With Gateway... "+util.IntToString(dash.t.TotalBlocks)+" Blocks Stored")),
				),
				elem.Div(
					vecty.Markup(vecty.Class("card-col-3")),
					vecty.Text("Total Txs: "+util.IntToString(dash.t.TxCount)),
				),
			),
			vecty.Markup(vecty.Class("home")),
			&BlockList{
				t:             dash.t,
				refreshCh:     dash.refreshCh,
				disableUpdate: &dash.disableUpdate,
			},
		)
	}
	return vecty.Text("Connecting to blockchain clients")
}

func initBlockTxsDashboardView(t *rpc.TendermintInfo, BlockTxsDashboardView *BlockTxsDashboardView, cfg *config.Cfg) *BlockTxsDashboardView {
	// Init tendermint client
	tClient := rpc.StartClient(cfg.TendermintHost)
	if tClient == nil {
		return BlockTxsDashboardView
	}
	BlockTxsDashboardView.tClient = tClient
	BlockTxsDashboardView.t = t
	BlockTxsDashboardView.quitCh = make(chan struct{})
	BlockTxsDashboardView.refreshCh = make(chan int, 50)
	BlockTxsDashboardView.blockIndex = 0
	BlockTxsDashboardView.disableUpdate = false
	BeforeUnload(func() {
		close(BlockTxsDashboardView.quitCh)
	})
	go updateAndRenderBlockTxsDashboard(BlockTxsDashboardView, cfg)
	return BlockTxsDashboardView
}

func updateAndRenderBlockTxsDashboard(d *BlockTxsDashboardView, cfg *config.Cfg) {
	ticker := time.NewTicker(time.Duration(cfg.RefreshTime) * time.Second)
	// Wait for data structs to load
	for d == nil || d.t == nil {
	}
	rpc.UpdateTendermintInfo(d.tClient, d.t, d.blockIndex)
	d.t.TotalBlocks = int(dbapi.GetBlockHeight())
	updateBlocks(d, util.Max(d.t.TotalBlocks-d.blockIndex-config.ListSize+1, 1))
	vecty.Rerender(d)
	for {
		select {
		case <-d.quitCh:
			ticker.Stop()
			fmt.Println("Gateway connection closed")
			return
		case <-ticker.C:
			rpc.UpdateTendermintInfo(d.tClient, d.t, d.blockIndex)
			d.t.TotalBlocks = int(dbapi.GetBlockHeight()) - 1
			if !d.disableUpdate {
				updateBlocks(d, util.Max(d.t.TotalBlocks-d.blockIndex-config.ListSize+1, 1))
			}
			vecty.Rerender(d)
		case i := <-d.refreshCh:
		loop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-d.refreshCh:
				default:
					break loop
				}
			}
			d.blockIndex = i
			oldBlockTxs := d.t.TotalBlocks
			d.t.TotalBlocks = int(dbapi.GetBlockHeight()) - 1
			if i < 1 {
				oldBlockTxs = d.t.TotalBlocks
			}
			updateBlocks(d, util.Max(oldBlockTxs-d.blockIndex-config.ListSize+1, 1))
			vecty.Rerender(d)
		}
	}
}

func updateBlocks(d *BlockTxsDashboardView, index int) {
	fmt.Println("Getting BlockTxs from index " + util.IntToString(index))
	d.t.BlockList = dbapi.GetBlockList(index)
}
