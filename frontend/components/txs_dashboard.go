package components

import (
	"fmt"
	"time"

	"github.com/hexops/vecty"
	"github.com/vocdoni/vocexplorer/api"
	"github.com/vocdoni/vocexplorer/api/dbtypes"
	"github.com/vocdoni/vocexplorer/config"
	"github.com/vocdoni/vocexplorer/frontend/actions"
	"github.com/vocdoni/vocexplorer/frontend/dispatcher"
	"github.com/vocdoni/vocexplorer/frontend/store"
	"github.com/vocdoni/vocexplorer/frontend/update"
	"github.com/vocdoni/vocexplorer/logger"
	"github.com/vocdoni/vocexplorer/util"
)

// TxsDashboardView renders the dashboard landing page
type TxsDashboardView struct {
	vecty.Core
	vecty.Mounter
	Rendered bool
}

// Mount is called after the component renders to signal that it can be rerendered safely
func (dash *TxsDashboardView) Mount() {
	if !dash.Rendered {
		dash.Rendered = true
		vecty.Rerender(dash)
	}
}

// Render renders the TxsDashboardView component
func (dash *TxsDashboardView) Render() vecty.ComponentOrHTML {
	if !dash.Rendered {
		return LoadingBar()
	}
	return Container(
		vecty.Markup(vecty.Attribute("id", "main")),
		renderServerConnectionBanner(),
		&TxList{},
	)
}

// UpdateTxsDashboard keeps the transactions dashboard updated
func UpdateTxsDashboard(d *TxsDashboardView) {
	dispatcher.Dispatch(&actions.EnableAllUpdates{})
	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	if !update.CheckCurrentPage("txs", ticker) {
		return
	}
	updateTxsDashboard(d)
	for {
		select {
		case <-store.RedirectChan:
			if !update.CheckCurrentPage("txs", ticker) {
				return
			}
		case <-ticker.C:
			if !update.CheckCurrentPage("txs", ticker) {
				return
			}
			updateTxsDashboard(d)
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
			if !update.CheckCurrentPage("txs", ticker) {
				return
			}
			dispatcher.Dispatch(&actions.TransactionsIndexChange{Index: i})
			if i < 1 {
				newHeight, _ := api.GetTxHeight()
				dispatcher.Dispatch(&actions.SetTransactionCount{Count: int(newHeight) - 1})
			}
			updateTxs(d, util.Max(store.Transactions.Count-store.Transactions.Pagination.Index, 1))

		case search := <-store.Transactions.Pagination.SearchChannel:
			if !update.CheckCurrentPage("txs", ticker) {
				return
			}
		txsearch:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case search = <-store.Transactions.Pagination.SearchChannel:
				default:
					break txsearch
				}
			}
			logger.Info("search: " + search)
			dispatcher.Dispatch(&actions.TransactionsIndexChange{Index: 0})
			list, ok := api.GetTransactionSearch(search)
			if ok {
				reverseTxList(&list)
				dispatcher.Dispatch(&actions.SetTransactionList{TransactionList: list})
			} else {
				dispatcher.Dispatch(&actions.SetTransactionList{TransactionList: [config.ListSize]*dbtypes.Transaction{nil}})
			}
		}
	}
}

func updateTxsDashboard(d *TxsDashboardView) {
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})

	if !store.Transactions.Pagination.DisableUpdate {
		actions.UpdateCounts()
		updateTxs(d, util.Max(store.Transactions.Count-store.Transactions.Pagination.Index, 1))
	}
}

func updateTxs(d *TxsDashboardView, index int) {
	logger.Info(fmt.Sprintf("Getting transactions from index %d\n", index))
	list, ok := api.GetTxList(index)
	if ok {
		reverseTxList(&list)
		dispatcher.Dispatch(&actions.SetTransactionList{TransactionList: list})
	}
}

func reverseTxList(list *[config.ListSize]*dbtypes.Transaction) {
	for i := len(list)/2 - 1; i >= 0; i-- {
		opp := len(list) - 1 - i
		list[i], list[opp] = list[opp], list[i]
	}
}
