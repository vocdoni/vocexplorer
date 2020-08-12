package components

import (
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/tendermint/tendermint/rpc/client/http"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/dbapi"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// BlockTxsDashboardView renders the dashboard landing page
type BlockTxsDashboardView struct {
	vecty.Core
	blockIndex          int
	blockRefresh        chan int
	disableBlocksUpdate bool
	disableTxsUpdate    bool
	quitCh              chan struct{}
	t                   *rpc.TendermintInfo
	tClient             *http.HTTP
	txIndex             int
	txRefresh           chan int
}

// Render renders the BlockTxsDashboardView component
func (dash *BlockTxsDashboardView) Render() vecty.ComponentOrHTML {
	if dash != nil && dash.tClient != nil && dash.t != nil && dash.t.ResultStatus != nil {
		return elem.Main(
			elem.Div(
				elem.Div(vecty.Markup(vecty.Class("card-col-3")),
					vecty.Text("Current Block Height: "+util.IntToString(dash.t.ResultStatus.SyncInfo.LatestBlockHeight)),
				),
				vecty.If(int(dash.t.ResultStatus.SyncInfo.LatestBlockHeight)-dash.t.TotalBlocks > 1,
					elem.Div(vecty.Markup(vecty.Class("card-col-3")),
						vecty.Text("Still Syncing With Gateway... "+util.IntToString(dash.t.TotalBlocks)+" Blocks Stored")),
				),
				elem.Div(
					vecty.Markup(vecty.Class("card-col-3")),
					vecty.Text("Total Txs: "+util.IntToString(dash.t.TotalTxs)),
				),
			),
			vecty.Markup(vecty.Class("home")),
			&BlockList{
				t:             dash.t,
				refreshCh:     dash.blockRefresh,
				disableUpdate: &dash.disableBlocksUpdate,
			},
			&TxList{
				t:             dash.t,
				refreshCh:     dash.txRefresh,
				disableUpdate: &dash.disableTxsUpdate,
			},
		)
	}
	return vecty.Text("Connecting to blockchain clients")
}

// InitBlockTxsDashboardView initializes the blocks & transactions view (to be splitted)
func InitBlockTxsDashboardView(t *rpc.TendermintInfo, BlockTxsDashboardView *BlockTxsDashboardView, cfg *config.Cfg) *BlockTxsDashboardView {
	// Init tendermint client
	tClient := rpc.StartClient(cfg.TendermintHost)
	if tClient == nil {
		return BlockTxsDashboardView
	}
	BlockTxsDashboardView.tClient = tClient
	BlockTxsDashboardView.t = t
	BlockTxsDashboardView.quitCh = make(chan struct{})
	BlockTxsDashboardView.blockRefresh = make(chan int, 50)
	BlockTxsDashboardView.txRefresh = make(chan int, 50)
	BlockTxsDashboardView.blockIndex = 0
	BlockTxsDashboardView.txIndex = 0
	BlockTxsDashboardView.disableBlocksUpdate = false
	BlockTxsDashboardView.disableTxsUpdate = false
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
	rpc.UpdateTendermintInfo(d.tClient, d.t)
	d.t.TotalBlocks = int(dbapi.GetBlockHeight()) - 1
	d.t.TotalTxs = int(dbapi.GetTxHeight()) - 1
	updateBlocks(d, util.Max(d.t.TotalBlocks-d.blockIndex, config.ListSize))
	updateTxs(d, util.Max(d.t.TotalTxs-d.txIndex, config.ListSize))
	vecty.Rerender(d)
	for {
		select {
		case <-d.quitCh:
			ticker.Stop()
			fmt.Println("Gateway connection closed")
			return
		case <-ticker.C:
			rpc.UpdateTendermintInfo(d.tClient, d.t)
			d.t.TotalBlocks = int(dbapi.GetBlockHeight()) - 1
			d.t.TotalTxs = int(dbapi.GetTxHeight()) - 1
			if !d.disableBlocksUpdate {
				updateBlocks(d, util.Max(d.t.TotalBlocks-d.blockIndex, config.ListSize))
			}
			if !d.disableTxsUpdate {
				updateTxs(d, util.Max(d.t.TotalTxs-d.txIndex, config.ListSize))
			}
			vecty.Rerender(d)
		case i := <-d.blockRefresh:
		blockloop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-d.blockRefresh:
				default:
					break blockloop
				}
			}
			d.blockIndex = i
			oldBlocks := d.t.TotalBlocks
			d.t.TotalBlocks = int(dbapi.GetBlockHeight()) - 1
			if i < 1 {
				oldBlocks = d.t.TotalBlocks
			}
			updateBlocks(d, util.Max(oldBlocks-d.blockIndex, config.ListSize))

			vecty.Rerender(d)
		case i := <-d.txRefresh:
		txloop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-d.txRefresh:
				default:
					break txloop
				}
			}
			d.txIndex = i
			oldTxs := d.t.TotalTxs
			d.t.TotalTxs = int(dbapi.GetTxHeight()) - 1
			if i < 1 {
				oldTxs = d.t.TotalTxs
			}
			updateTxs(d, util.Max(oldTxs-d.txIndex, config.ListSize))

			vecty.Rerender(d)
		}
	}
}

func updateBlocks(d *BlockTxsDashboardView, index int) {
	log.Infof("Getting Blocks from index %d", util.IntToString(index))
	list := dbapi.GetBlockList(index)
	for i := len(list)/2 - 1; i >= 0; i-- {
		opp := len(list) - 1 - i
		list[i], list[opp] = list[opp], list[i]
	}
	d.t.BlockList = list
}

func updateTxs(d *BlockTxsDashboardView, index int) {
	log.Infof("Getting Txs from index %d", util.IntToString(index))
	list := dbapi.GetTxList(index)
	for i := len(list)/2 - 1; i >= 0; i-- {
		opp := len(list) - 1 - i
		list[i], list[opp] = list[opp], list[i]
	}
	d.t.TxList = list
}
