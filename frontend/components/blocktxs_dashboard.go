package components

import (
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
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
	vecty.Mounter
	Rendered bool
}

// Mount is called after the component renders to signal that it can be rerendered safely
func (dash *BlockTxsDashboardView) Mount() {
	if !dash.Rendered {
		dash.Rendered = true
		vecty.Rerender(dash)
	}
}

// Render renders the BlockTxsDashboardView component
func (dash *BlockTxsDashboardView) Render() vecty.ComponentOrHTML {
	if !dash.Rendered {
		return LoadingBar()
	}
	if dash != nil && len(store.Blocks.Blocks) > 0 {
		return Container(
			renderGatewayConnectionBanner(),
			renderServerConnectionBanner(),
			// &LatestBlocksWidget{},
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
			dispatcher.Dispatch(&actions.BlocksIndexChange{Index: i})
			oldBlocks := store.Blocks.Count
			newHeight, _ := api.GetBlockHeight()
			dispatcher.Dispatch(&actions.BlocksHeightUpdate{Height: int(newHeight) - 1})
			if i < 1 {
				oldBlocks = store.Blocks.Count
			}
			updateBlocks(d, util.Max(oldBlocks-store.Blocks.Pagination.Index, config.ListSize))
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
			dispatcher.Dispatch(&actions.TransactionsIndexChange{Index: i})
			oldTxs := store.Transactions.Count
			newHeight, _ := api.GetTxHeight()
			dispatcher.Dispatch(&actions.SetTransactionCount{Count: int(newHeight) - 1})
			if i < 1 {
				oldTxs = store.Transactions.Count
			}
			updateTxs(d, util.Max(oldTxs-store.Transactions.Pagination.Index, config.ListSize))
		}
	}
}

func updateBlockTxsDashboard(d *BlockTxsDashboardView) {
	dispatcher.Dispatch(&actions.GatewayConnected{Connected: api.PingGateway(store.Config.GatewayHost)})
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.Ping()})

	actions.UpdateCounts()
	rpc.UpdateBlockchainStatus(store.TendermintClient)
	if !store.Blocks.Pagination.DisableUpdate {
		updateBlocks(d, util.Max(store.Blocks.Count-store.Blocks.Pagination.Index, config.ListSize))
	}
	if !store.Transactions.Pagination.DisableUpdate {
		updateTxs(d, util.Max(store.Transactions.Count-store.Transactions.Pagination.Index, config.ListSize))
	}
}

func updateBlocks(d *BlockTxsDashboardView, index int) {
	fmt.Printf("Getting Blocks from index %d\n", index)
	list, ok := api.GetBlockList(index)
	if ok {
		reverseBlockList(&list)
		dispatcher.Dispatch(&actions.SetBlockList{BlockList: list})
	}
}

func updateTxs(d *BlockTxsDashboardView, index int) {
	fmt.Printf("Getting Txs from index %d\n", index)
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
