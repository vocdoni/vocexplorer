package components

import (
	"fmt"
	"time"

	"github.com/hexops/vecty"

	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/logger"
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
				_, newHeight, _, err := store.Client.GetBlockStatus()
				if err != nil {
					logger.Error(err)
				}
				dispatcher.Dispatch(&actions.BlocksHeightUpdate{Height: int(*newHeight) - 1})
			}
			logger.Info(fmt.Sprintf("update blocks to index %d\n", i))
			updateBlocks(d, store.Blocks.Count-store.Blocks.Pagination.Index-config.ListSize+1)
		}
	}
}

func updateBlocksDashboard(d *BlocksDashboardView) {
	dispatcher.Dispatch(&actions.GatewayConnected{GatewayErr: store.Client.GetGatewayInfo()})
	if !store.Blocks.Pagination.DisableUpdate {
		stats, err := store.Client.GetStats()
		if err != nil {
			logger.Error(err)
			return
		}
		actions.UpdateCounts(stats)
		updateBlocks(d, store.Blocks.Count-store.Blocks.Pagination.Index-config.ListSize+1)
	}
}

func updateBlocks(d *BlocksDashboardView, index int) {
	listSize := config.ListSize
	if index < 0 {
		listSize += index
		index = 0
	}
	logger.Info(fmt.Sprintf("Getting %d blocks from index %d\n", listSize, index))
	list, err := store.Client.GetBlockList(index, listSize)
	if err != nil {
		logger.Error(err)
		return
	}
	dispatcher.Dispatch(&actions.SetBlockList{BlockList: list})
}
