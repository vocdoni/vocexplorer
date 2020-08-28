package client

import (
	"strings"

	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// UpdateDashboardInfo calls gateway apis, updates info needed for dashboard page
func UpdateDashboardInfo(c *client.Client) {
	UpdateGatewayInfo(c)
	UpdateBlockStatus(c)
	UpdateCounts(c)
}

// UpdateCounts calls gateway apis, updates total number of processes and entities
func UpdateCounts(c *client.Client) {
	procs, err := c.GetProcessCount()
	util.ErrPrint(err)
	entities, err := c.GetEntityCount()
	util.ErrPrint(err)
	store.Processes.ProcessCount = int(procs)
	store.Entities.EntityCount = int(entities)

}

// UpdateGatewayInfo calls gateway api, updates gateway health info
func UpdateGatewayInfo(c *client.Client) {
	apiList, health, ok, err := c.GetGatewayInfo()
	util.ErrPrint(err)
	store.Stats.APIList = apiList
	store.Stats.Ok = ok
	store.Stats.Health = health
}

// UpdateBlockStatus calls gateway api, updates blockchain statistics
func UpdateBlockStatus(c *client.Client) {
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

// UpdateProcessResults updates auxilary info for all currently displayed process id's
func UpdateProcessResults(c *client.Client) {
	if store.Processes.ProcessResults == nil {
		store.Processes.ProcessResults = make(map[string]storeutil.Process)
	}
	if store.Processes.EnvelopeHeights == nil {
		store.Processes.EnvelopeHeights = make(map[string]int64)
	}
	for _, ID := range store.Processes.ProcessIDs {
		if ID != "" {
			if _, ok := store.Processes.ProcessResults[ID]; !ok {
				t, st, res, err := c.GetProcessResults(strings.ToLower(ID))
				if !util.ErrPrint(err) {
					store.Processes.ProcessResults[ID] = storeutil.Process{
						ProcessType: t,
						State:       st,
						Results:     res}
				}
			}
		}
	}
}

// UpdateProcessesDashboardInfo updates process info to include status and recent envelopes
func UpdateProcessesDashboardInfo(c *client.Client, process *storeutil.Process, processID string) {
	if process == nil {
		process = new(storeutil.Process)
	}
	t, st, res, err := c.GetProcessResults(processID)
	if !util.ErrPrint(err) {
		process.ProcessType = t
		process.Results = res
		process.State = st
	}
}

// UpdateAuxEntityInfo updates process info map to include all currently displayed process IDs
func UpdateAuxEntityInfo(c *client.Client, e *storeutil.Entity) {
	if e.Processes == nil {
		e.Processes = make(map[string]storeutil.Process)
	}
	if e.EnvelopeHeights == nil {
		e.EnvelopeHeights = make(map[string]int64)
	}
	for _, ID := range e.ProcessIDs {
		if ID != "" {
			if _, ok := e.Processes[ID]; !ok {
				t, st, res, err := c.GetProcessResults(strings.ToLower(ID))
				if !util.ErrPrint(err) {
					e.Processes[ID] = storeutil.Process{
						ProcessType: t,
						State:       st,
						Results:     res}
				}
			}
		}
	}
}
