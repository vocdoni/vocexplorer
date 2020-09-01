package update

import (
	"strings"

	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// DashboardInfo calls gateway apis, updates info needed for dashboard page
func DashboardInfo(c *api.GatewayClient) {
	GatewayInfo(c)
	BlockStatus(c)
	// Counts(c)
}

// Counts calls gateway apis, updates total number of processes and entities
func Counts(c *api.GatewayClient) {
	procs, err := c.GetProcessCount()
	util.ErrPrint(err)
	entities, err := c.GetEntityCount()
	util.ErrPrint(err)
	store.Processes.Count = int(procs)
	store.Entities.Count = int(entities)

}

// GatewayInfo calls gateway api, updates gateway health info
func GatewayInfo(c *api.GatewayClient) {
	apiList, health, ok, err := c.GetGatewayInfo()
	util.ErrPrint(err)
	dispatcher.Dispatch(&actions.SetGatewayInfo{
		APIList: apiList,
		Ok:      ok,
		Health:  health,
	})
}

// BlockStatus calls gateway api, updates blockchain statistics
func BlockStatus(c *api.GatewayClient) {
	blockTime, blockTimeStamp, height, err := c.GetBlockStatus()
	util.ErrPrint(err)
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
	util.ErrPrint(err)
}

// ProcessResults updates auxilary info for all currently displayed process id's
func ProcessResults() {
	for _, ID := range store.Processes.ProcessIDs {
		if ID != "" {
			if _, ok := store.Processes.ProcessResults[ID]; !ok {
				t, st, res, err := store.GatewayClient.GetProcessResults(strings.ToLower(ID))
				if !util.ErrPrint(err) {
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

// CurrentProcessResults updates current process information
func CurrentProcessResults() {
	t, st, res, err := store.GatewayClient.GetProcessResults(store.Processes.CurrentProcessID)
	if !util.ErrPrint(err) {
		dispatcher.Dispatch(&actions.SetCurrentProcess{
			Process: storeutil.Process{
				ProcessType: t,
				State:       st,
				Results:     res},
		})
	}
}

// EntityProcessResults ensures the given entity's processes' results are all stored
func EntityProcessResults(e storeutil.Entity) {
	for _, ID := range e.ProcessIDs {
		if ID != "" {
			if _, ok := store.Processes.ProcessResults[ID]; !ok {
				t, st, res, err := store.GatewayClient.GetProcessResults(strings.ToLower(ID))
				if !util.ErrPrint(err) {
					dispatcher.Dispatch(&actions.SetProcessContents{
						ID: store.Processes.CurrentProcessID,
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
