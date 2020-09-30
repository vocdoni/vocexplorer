package update

import (
	"fmt"
	"strings"
	"time"

	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/api/rpc"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// Gateway API updates

// DashboardInfo calls gateway apis, updates info needed for dashboard page
func DashboardInfo(c *api.GatewayClient) {
	GatewayInfo(c)
	BlockStatus(c)
}

// GatewayInfo calls gateway api, updates gateway health info
func GatewayInfo(c *api.GatewayClient) {
	apiList, health, err := c.GetGatewayInfo()
	if err != nil {
		log.Warn(err)
	}
	dispatcher.Dispatch(&actions.SetGatewayInfo{
		APIList: apiList,
		Health:  health,
	})
}

// BlockStatus calls gateway api, updates blockchain statistics
func BlockStatus(c *api.GatewayClient) {
	blockTime, blockTimeStamp, height, err := c.GetBlockStatus()
	if err != nil {
		log.Warn(err)
	}
	dispatcher.Dispatch(&actions.SetBlockStatus{
		BlockTime:      blockTime,
		BlockTimeStamp: blockTimeStamp,
		Height:         height,
	})

}

// GetIDs gets ids
func GetIDs(IDList *[]string, c *api.GatewayClient, getList func() ([]string, error)) {
	var err error
	*IDList, err = getList()
	if err != nil {
		log.Warn(err)
	}
}

// ProcessResults updates auxilary info for all currently displayed process id's
func ProcessResults() {
	for _, process := range store.Processes.Processes {
		if process != nil {
			ID := process.ID
			if ID != "" {
				if _, ok := store.Processes.ProcessResults[ID]; !ok {
					t, st, res, err := store.GatewayClient.GetProcessResults(strings.ToLower(ID))
					if err != nil {
						log.Warn(err)
					} else {
						dispatcher.Dispatch(&actions.SetProcessContents{
							ID: ID,
							Process: storeutil.Process{
								ProcessType: t,
								State:       st,
								Results:     res},
						})
					}
				}
			}
		}
	}
}

// EnvelopeProcessResults updates auxilary info for all process id's belonging to currently displayed envelopes
func EnvelopeProcessResults() {
	for _, envelope := range store.Envelopes.Envelopes {
		if envelope != nil {
			ID := strings.ToLower(util.TrimHex(envelope.ProcessID))
			if ID != "" {
				if _, ok := store.Processes.ProcessResults[ID]; !ok {
					t, st, res, err := store.GatewayClient.GetProcessResults(ID)
					if err != nil {
						log.Warn(err)
					} else {
						dispatcher.Dispatch(&actions.SetProcessContents{
							ID: ID,
							Process: storeutil.Process{
								ProcessType: t,
								State:       st,
								Results:     res},
						})
					}
				}
			}
		}
	}
}

// CurrentProcessResults updates current process information
func CurrentProcessResults() {
	t, st, res, err := store.GatewayClient.GetProcessResults(store.Processes.CurrentProcess.ID)
	if err != nil {
		log.Warn(err)
	} else {
		dispatcher.Dispatch(&actions.SetCurrentProcess{
			Process: storeutil.Process{
				ProcessType: t,
				State:       st,
				Results:     res},
		})
	}
}

// EntityProcessResults ensures the given entity's processes' results are all stored
func EntityProcessResults() {
	for _, process := range store.Entities.CurrentEntity.Processes {
		if process != nil {
			ID := process.ID
			if ID != "" {
				if _, ok := store.Processes.ProcessResults[ID]; !ok {
					t, st, res, err := store.GatewayClient.GetProcessResults(strings.ToLower(ID))
					if err != nil {
						log.Warn(err)
					} else {
						dispatcher.Dispatch(&actions.SetProcessContents{
							ID: ID,
							Process: storeutil.Process{
								ProcessType: t,
								State:       st,
								Results:     res},
						})
					}
				}
			}
		}
	}
}

// Tendermint API updates

//BlockchainStatus updates the blockchain statistics
func BlockchainStatus(t *rpc.TendermintRPC) {
	status := api.GetHealth(t)
	genesis := api.GetGenesis(t)
	dispatcher.Dispatch(&actions.SetResultStatus{Status: status})
	dispatcher.Dispatch(&actions.SetGenesis{Genesis: genesis})
}

// CheckCurrentPage returns true and stops ticker if the current page is title
func CheckCurrentPage(title string, ticker *time.Ticker) bool {
	if store.CurrentPage != title {
		fmt.Println("redirecting")
		ticker.Stop()
		return false
	}
	return true
}
