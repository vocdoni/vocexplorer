package components

import (
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/dbapi"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// BlockTxsDashboardView renders the dashboard landing page
type BlockTxsDashboardView struct {
	vecty.Core
	gatewayConnected    bool
	serverConnected     bool
	blockIndex          int
	blockRefresh        chan int
	disableBlocksUpdate bool
	disableTxsUpdate    bool
	t                   *rpc.TendermintInfo
	txIndex             int
	txRefresh           chan int
}

// Render renders the BlockTxsDashboardView component
func (dash *BlockTxsDashboardView) Render() vecty.ComponentOrHTML {
	if dash != nil && store.Tendermint != nil && dash.t != nil && dash.t.ResultStatus != nil {
		return Container(
			renderGatewayConnectionBanner(dash.gatewayConnected),
			renderServerConnectionBanner(dash.serverConnected),
			&LatestBlocksWidget{
				T: dash.t,
			},
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
			&BlockchainInfo{
				T: dash.t,
			},
		)
	}
	return &bootstrap.Alert{
		Contents: "Connecting to blockchain clients",
		Type:     "warning",
	}
}

// InitBlockTxsDashboardView initializes the blocks & transactions view (to be splitted)
func InitBlockTxsDashboardView(t *rpc.TendermintInfo, BlockTxsDashboardView *BlockTxsDashboardView, cfg *config.Cfg) *BlockTxsDashboardView {
	BlockTxsDashboardView.t = t
	BlockTxsDashboardView.blockRefresh = make(chan int, 50)
	BlockTxsDashboardView.txRefresh = make(chan int, 50)
	BlockTxsDashboardView.blockIndex = 0
	BlockTxsDashboardView.txIndex = 0
	BlockTxsDashboardView.disableBlocksUpdate = false
	BlockTxsDashboardView.disableTxsUpdate = false
	BlockTxsDashboardView.gatewayConnected = true
	BlockTxsDashboardView.serverConnected = true
	BeforeUnload(func() {
		dispatcher.Dispatch(&actions.SignalRedirect{})
	})
	go updateAndRenderBlockTxsDashboard(BlockTxsDashboardView, cfg)
	return BlockTxsDashboardView
}

func updateAndRenderBlockTxsDashboard(d *BlockTxsDashboardView, cfg *config.Cfg) {
	ticker := time.NewTicker(time.Duration(cfg.RefreshTime) * time.Second)
	updateBlockTxsDashboard(d)
	vecty.Rerender(d)
	for {
		select {
		case <-store.RedirectChan:
			fmt.Println("Redirecting...")
			ticker.Stop()
			return
		case <-ticker.C:
			updateBlockTxsDashboard(d)
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
			newHeight, _ := dbapi.GetBlockHeight()
			d.t.TotalBlocks = int(newHeight) - 1
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
			newHeight, _ := dbapi.GetTxHeight()
			d.t.TotalTxs = int(newHeight) - 1
			if i < 1 {
				oldTxs = d.t.TotalTxs
			}
			updateTxs(d, util.Max(oldTxs-d.txIndex, config.ListSize))

			vecty.Rerender(d)
		}
	}
}

func updateBlockTxsDashboard(d *BlockTxsDashboardView) {
	if !rpc.Ping(store.TendermintClient) {
		d.gatewayConnected = false
	} else {
		d.gatewayConnected = true
	}
	if !dbapi.Ping() {
		d.serverConnected = false
	} else {
		d.serverConnected = true
	}
	updateHeight(d.t)
	rpc.UpdateTendermintInfo(store.TendermintClient, d.t)
	if !d.disableBlocksUpdate {
		updateBlocks(d, util.Max(d.t.TotalBlocks-d.blockIndex, config.ListSize))
	}
	if !d.disableTxsUpdate {
		updateTxs(d, util.Max(d.t.TotalTxs-d.txIndex, config.ListSize))
	}
}

func updateBlocks(d *BlockTxsDashboardView, index int) {
	log.Infof("Getting Blocks from index %d", util.IntToString(index))
	list, ok := dbapi.GetBlockList(index)
	if ok {
		reverseBlockList(&list)
		d.t.BlockList = list
	}
}

func updateTxs(d *BlockTxsDashboardView, index int) {
	log.Infof("Getting Txs from index %d", util.IntToString(index))
	list, ok := dbapi.GetTxList(index)
	if ok {
		reverseTxList(&list)
		d.t.TxList = list
	}
}

func reverseBlockList(list *[config.ListSize]*types.StoreBlock) {
	for i := len(list)/2 - 1; i >= 0; i-- {
		opp := len(list) - 1 - i
		list[i], list[opp] = list[opp], list[i]
	}
}

func reverseTxList(list *[config.ListSize]*types.SendTx) {
	for i := len(list)/2 - 1; i >= 0; i-- {
		opp := len(list) - 1 - i
		list[i], list[opp] = list[opp], list[i]
	}
}
