package update

import (
	"strings"

	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// DashboardInfo calls gateway apis, updates info needed for dashboard page
func DashboardInfo(c *client.Client) {
	GatewayInfo(c)
	BlockStatus(c)
	Counts(c)
}

// Counts calls gateway apis, updates total number of processes and entities
func Counts(c *client.Client) {
	procs, err := c.GetProcessCount()
	util.ErrPrint(err)
	entities, err := c.GetEntityCount()
	util.ErrPrint(err)
	store.Processes.ProcessCount = int(procs)
	store.Entities.EntityCount = int(entities)

}

// GatewayInfo calls gateway api, updates gateway health info
func GatewayInfo(c *client.Client) {
	apiList, health, ok, err := c.GetGatewayInfo()
	util.ErrPrint(err)
	store.Stats.APIList = apiList
	store.Stats.Ok = ok
	store.Stats.Health = health
}

// BlockStatus calls gateway api, updates blockchain statistics
func BlockStatus(c *client.Client) {
	blockTime, blockTimeStamp, height, err := c.GetBlockStatus()
	util.ErrPrint(err)
	store.Stats.BlockTime = blockTime
	store.Stats.BlockTimeStamp = blockTimeStamp
	store.Stats.Height = height
}

// GetIDs gets ids
func GetIDs(IDList *[]string, c *client.Client, getList func() ([]string, error)) {
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
		dispatcher.Dispatch(&actions.SetProcessContents{
			ID: store.Processes.CurrentProcessID,
			Process: storeutil.Process{
				ProcessType: t,
				State:       st,
				Results:     res},
		})
	}
}

// EntityProcessResults ensures the given entity's processes' results are all stored
func EntityProcessResults(c *client.Client, e *storeutil.Entity) {
	for _, ID := range e.ProcessIDs {
		if ID != "" {
			if _, ok := store.Processes.ProcessResults[ID]; !ok {
				t, st, res, err := c.GetProcessResults(strings.ToLower(ID))
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
