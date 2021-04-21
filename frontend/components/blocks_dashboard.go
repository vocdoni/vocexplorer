package components

import (
	"fmt"
	"time"

	"github.com/hexops/vecty"
	"github.com/vocdoni/vocexplorer/api"
	"go.vocdoni.io/proto/build/go/models"

	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/logger"
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
	return Container(
		vecty.Markup(vecty.Attribute("id", "main")),
		renderServerConnectionBanner(),
		&BlockList{},
	)
}

// UpdateBlocksDashboard keeps the blocks dashboard updated
func UpdateBlocksDashboard(d *BlocksDashboardView) {
	dispatcher.Dispatch(&actions.EnableAllUpdates{})

	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	if !update.CheckCurrentPage("blocks", ticker) {
		return
	}
	updateBlocksDashboard(d)
	for {
		select {
		case <-store.RedirectChan:
			if !update.CheckCurrentPage("blocks", ticker) {
				return
			}
		case <-ticker.C:
			if !update.CheckCurrentPage("blocks", ticker) {
				return
			}
			updateBlocksDashboard(d)
		case i := <-store.Blocks.Pagination.PagChannel:
			if !update.CheckCurrentPage("blocks", ticker) {
				return
			}
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
			if i < 1 { // If on first page, update counts
				newHeight, _ := api.GetBlockHeight()
				dispatcher.Dispatch(&actions.BlocksHeightUpdate{Height: int(newHeight) - 1})
			}
			logger.Info(fmt.Sprintf("update blocks to index %d\n", i))
			updateBlocks(d, util.Max(store.Blocks.Count-store.Blocks.Pagination.Index, 1))

			// case search := <-store.Blocks.Pagination.SearchChannel:
			// 	if !update.CheckCurrentPage("blocks", ticker) {
			// 		return
			// 	}
			// blocksearch:
			// 	for {
			// 		// If many indices waiting in buffer, scan to last one.
			// 		select {
			// 		case search = <-store.Blocks.Pagination.SearchChannel:
			// 		default:
			// 			break blocksearch
			// 		}
			// 	}
			// 	logger.Info("search: " + search)
			// 	dispatcher.Dispatch(&actions.BlocksIndexChange{Index: 0})
			// 	list, ok := store.Client.GetBlockSearch(search)
			// 	if ok {
			// 		reverseBlockList(&list)
			// 		dispatcher.Dispatch(&actions.SetBlockList{BlockList: list})
			// 	} else {
			// 		dispatcher.Dispatch(&actions.SetBlockList{BlockList: [config.ListSize]*dbtypes.StoreBlock{nil}})
			// 	}
		}
	}
}

func updateBlocksDashboard(d *BlocksDashboardView) {
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})
	if !store.Blocks.Pagination.DisableUpdate {
		stats, err := store.Client.GetStats()
		if err != nil {
			logger.Error(err)
			return
		}
		actions.UpdateCounts(stats)
		updateBlocks(d, util.Max(store.Blocks.Count-store.Blocks.Pagination.Index, 1))
	}
}

func updateBlocks(d *BlocksDashboardView, index int) {
	logger.Info(fmt.Sprintf("Getting Blocks from index %d\n", index))
	// list, err := store.Client.GetBlockList(index)
	list, err := store.Client.GetBlockList(index, config.ListSize)
	if err != nil {
		logger.Error(err)
		return
	}
	reverseBlockList(list)
	dispatcher.Dispatch(&actions.SetBlockList{BlockList: list})
}

func reverseBlockList(list *models.BlockHeaderList) {
	for i := len(list.BlockHeaders)/2 - 1; i >= 0; i-- {
		opp := len(list.BlockHeaders) - 1 - i
		list.BlockHeaders[i], list.BlockHeaders[opp] = list.BlockHeaders[opp], list.BlockHeaders[i]
	}
}
