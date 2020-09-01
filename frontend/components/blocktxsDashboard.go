package components

import (
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/api"
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
	blockIndex int
	txIndex    int
}

// Render renders the BlockTxsDashboardView component
func (dash *BlockTxsDashboardView) Render() vecty.ComponentOrHTML {
	if dash != nil && len(store.Blocks.Blocks) > 0 {
		return Container(
			renderGatewayConnectionBanner(),
			renderServerConnectionBanner(),
			&LatestBlocksWidget{},
			&BlockList{},
			&TxList{},
			&BlockchainInfo{},
		)
	}
	return &bootstrap.Alert{
		Contents: "Connecting to blockchain clients",
		Type:     "warning",
	}
}

// UpdateAndRenderBlockTxsDashboard keeps the block transactions dashboard updated
func UpdateAndRenderBlockTxsDashboard(d *BlockTxsDashboardView) {
	actions.EnableUpdates()
	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	updateBlockTxsDashboard(d)
	for {
		select {
		case <-store.RedirectChan:
			fmt.Println("Redirecting...")
			ticker.Stop()
			return
		case <-ticker.C:
			updateBlockTxsDashboard(d)
		case i := <-store.Blocks.Pagination.PagChannel:
		blockloop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-store.Blocks.Pagination.PagChannel:
				default:
					break blockloop
				}
			}
			d.blockIndex = i
			oldBlocks := store.Blocks.Count
			newHeight, _ := api.GetBlockHeight()
			dispatcher.Dispatch(&actions.BlocksHeightUpdate{Height: int(newHeight) - 1})
			if i < 1 {
				oldBlocks = store.Blocks.Count
			}
			updateBlocks(d, util.Max(oldBlocks-d.blockIndex, config.ListSize))
		case i := <-store.Transactions.Pagination.PagChannel:
		txloop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-store.Transactions.Pagination.PagChannel:
				default:
					break txloop
				}
			}
			d.txIndex = i
			oldTxs := store.Transactions.Count
			newHeight, _ := api.GetTxHeight()
			dispatcher.Dispatch(&actions.SetTransactionCount{Count: int(newHeight) - 1})
			if i < 1 {
				oldTxs = store.Transactions.Count
			}
			updateTxs(d, util.Max(oldTxs-d.txIndex, config.ListSize))
		}
	}
}

func updateBlockTxsDashboard(d *BlockTxsDashboardView) {
	dispatcher.Dispatch(&actions.GatewayConnected{Connected: api.PingGateway(store.Config.GatewayHost)})
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.Ping()})

	actions.UpdateCounts()
	rpc.UpdateBlockchainStatus(store.TendermintClient)
	if !store.Blocks.Pagination.DisableUpdate {
		updateBlocks(d, util.Max(store.Blocks.Count-d.blockIndex, config.ListSize))
	}
	if !store.Transactions.Pagination.DisableUpdate {
		updateTxs(d, util.Max(store.Transactions.Count-d.txIndex, config.ListSize))
	}
}

func updateBlocks(d *BlockTxsDashboardView, index int) {
	log.Infof("Getting Blocks from index %d", util.IntToString(index))
	list, ok := api.GetBlockList(index)
	if ok {
		reverseBlockList(&list)
		dispatcher.Dispatch(&actions.SetBlockList{BlockList: list})
	}
}

func updateTxs(d *BlockTxsDashboardView, index int) {
	log.Infof("Getting Txs from index %d", util.IntToString(index))
	list, ok := api.GetTxList(index)
	if ok {
		reverseTxList(&list)
		dispatcher.Dispatch(&actions.SetTransactionList{TransactionList: list})
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
