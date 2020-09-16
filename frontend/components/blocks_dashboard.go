package components

import (
	"fmt"
	"log"
	"time"

	"github.com/hexops/vecty"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/proto"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// BlocksDashboardView renders the dashboard landing page
type BlocksDashboardView struct {
	vecty.Core
	vecty.Mounter
	Rendered bool
}

// Mount is called after the component renders to signal that it can be rerendered safely
func (dash *BlocksDashboardView) Mount() {
	if !dash.Rendered {
		dash.Rendered = true
		vecty.Rerender(dash)
	}
}

// Render renders the BlocksDashboardView component
func (dash *BlocksDashboardView) Render() vecty.ComponentOrHTML {
	if !dash.Rendered {
		return LoadingBar()
	}
	if dash != nil {
		return Container(
			renderGatewayConnectionBanner(),
			renderServerConnectionBanner(),
			&BlockList{},
		)
	}
	return &bootstrap.Alert{
		Contents: "Connecting to blockchain clients",
		Type:     "warning",
	}
}

// UpdateBlocksDashboard keeps the blocks dashboard updated
func UpdateBlocksDashboard(d *BlocksDashboardView) {
	dispatcher.Dispatch(&actions.EnableAllUpdates{})

	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	updateBlocksDashboard(d)
	for {
		select {
		case <-store.RedirectChan:
			fmt.Println("Redirecting...")
			ticker.Stop()
			return
		case <-ticker.C:
			updateBlocksDashboard(d)
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
			updateBlocks(d, util.Max(oldBlocks-store.Blocks.Pagination.Index, 1))

		case search := <-store.Blocks.Pagination.SearchChannel:
		blocksearch:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case search = <-store.Blocks.Pagination.SearchChannel:
				default:
					break blocksearch
				}
			}
			log.Println("search: " + search)
			dispatcher.Dispatch(&actions.BlocksIndexChange{Index: 0})
			list, ok := api.GetBlockSearch(search)
			if ok {
				reverseBlockList(&list)
				dispatcher.Dispatch(&actions.SetBlockList{BlockList: list})
			} else {
				dispatcher.Dispatch(&actions.SetBlockList{BlockList: [config.ListSize]*proto.StoreBlock{nil}})
			}
		}
	}
}

func updateBlocksDashboard(d *BlocksDashboardView) {
	dispatcher.Dispatch(&actions.GatewayConnected{Connected: store.GatewayClient.Ping()})
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})

	actions.UpdateCounts()
	update.BlockchainStatus(store.TendermintClient)
	if !store.Blocks.Pagination.DisableUpdate {
		updateBlocks(d, util.Max(store.Blocks.Count-store.Blocks.Pagination.Index, 1))
	}
}

func updateBlocks(d *BlocksDashboardView, index int) {
	fmt.Printf("Getting Blocks from index %d\n", index)
	list, ok := api.GetBlockList(index)
	if ok {
		reverseBlockList(&list)
		dispatcher.Dispatch(&actions.SetBlockList{BlockList: list})
	}
}

func reverseBlockList(list *[config.ListSize]*proto.StoreBlock) {
	for i := len(list)/2 - 1; i >= 0; i-- {
		opp := len(list) - 1 - i
		list[i], list[opp] = list[opp], list[i]
	}
}
