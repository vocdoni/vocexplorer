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
	if dash != nil {
		return Container(
			renderGatewayConnectionBanner(),
			renderServerConnectionBanner(),
			&TxList{},
		)
	}
	return &bootstrap.Alert{
		Contents: "Connecting to blockchain clients",
		Type:     "warning",
	}
}

// UpdateTxsDashboard keeps the transactions dashboard updated
func UpdateTxsDashboard(d *TxsDashboardView) {
	dispatcher.Dispatch(&actions.EnableAllUpdates{})

	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	updateTxsDashboard(d)
	for {
		select {
		case <-store.RedirectChan:
			fmt.Println("Redirecting...")
			ticker.Stop()
			return
		case <-ticker.C:
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
			dispatcher.Dispatch(&actions.TransactionsIndexChange{Index: i})
			oldTxs := store.Transactions.Count
			newHeight, _ := api.GetTxHeight()
			dispatcher.Dispatch(&actions.SetTransactionCount{Count: int(newHeight) - 1})
			if i < 1 {
				oldTxs = store.Transactions.Count
			}
			updateTxs(d, util.Max(oldTxs-store.Transactions.Pagination.Index, 1))

		case search := <-store.Transactions.Pagination.SearchChannel:
		txsearch:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case search = <-store.Transactions.Pagination.SearchChannel:
				default:
					break txsearch
				}
			}
			log.Println("search: " + search)
			dispatcher.Dispatch(&actions.TransactionsIndexChange{Index: 0})
			list, ok := api.GetTransactionSearch(search)
			if ok {
				reverseTxList(&list)
				dispatcher.Dispatch(&actions.SetTransactionList{TransactionList: list})
			} else {
				dispatcher.Dispatch(&actions.SetTransactionList{TransactionList: [config.ListSize]*proto.SendTx{nil}})
			}
		}
	}
}

func updateTxsDashboard(d *TxsDashboardView) {
	dispatcher.Dispatch(&actions.GatewayConnected{Connected: store.GatewayClient.Ping()})
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})

	actions.UpdateCounts()
	update.BlockchainStatus(store.TendermintClient)
	if !store.Transactions.Pagination.DisableUpdate {
		updateTxs(d, util.Max(store.Transactions.Count-store.Transactions.Pagination.Index, 1))
	}
}

func updateTxs(d *TxsDashboardView, index int) {
	fmt.Printf("Getting transactions from index %d\n", index)
	list, ok := api.GetTxList(index)
	if ok {
		reverseTxList(&list)
		dispatcher.Dispatch(&actions.SetTransactionList{TransactionList: list})
	}
}

func reverseTxList(list *[config.ListSize]*proto.SendTx) {
	for i := len(list)/2 - 1; i >= 0; i-- {
		opp := len(list) - 1 - i
		list[i], list[opp] = list[opp], list[i]
	}
}
