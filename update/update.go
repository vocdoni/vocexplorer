package update

import (
	"strings"

	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
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
	if err != nil {
		log.Error(err)
	}
	entities, err := c.GetEntityCount()
	if err != nil {
		log.Error(err)
	}
	store.Processes.Count = int(procs)
	store.Entities.Count = int(entities)

}

// GatewayInfo calls gateway api, updates gateway health info
func GatewayInfo(c *api.GatewayClient) {
	apiList, health, ok, err := c.GetGatewayInfo()
	if err != nil {
		log.Error(err)
	}
	dispatcher.Dispatch(&actions.SetGatewayInfo{
		APIList: apiList,
		Ok:      ok,
		Health:  health,
	})
}

// BlockStatus calls gateway api, updates blockchain statistics
func BlockStatus(c *api.GatewayClient) {
	blockTime, blockTimeStamp, height, err := c.GetBlockStatus()
	if err != nil {
		log.Error(err)
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
		log.Error(err)
	}
}

// ProcessResults updates auxilary info for all currently displayed process id's
func ProcessResults() {
	for _, ID := range store.Processes.ProcessIDs {
		if ID != "" {
			if _, ok := store.Processes.ProcessResults[ID]; !ok {
				t, st, res, err := store.GatewayClient.GetProcessResults(strings.ToLower(ID))
				if err != nil {
					log.Error(err)
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

// CurrentProcessResults updates current process information
func CurrentProcessResults() {
	t, st, res, err := store.GatewayClient.GetProcessResults(store.Processes.CurrentProcessID)
	if err != nil {
		log.Error(err)
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
	for _, ID := range store.Entities.CurrentEntity.ProcessIDs {
		if ID != "" {
			if _, ok := store.Processes.ProcessResults[ID]; !ok {
				t, st, res, err := store.GatewayClient.GetProcessResults(strings.ToLower(ID))
				if err != nil {
					log.Error(err)
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
