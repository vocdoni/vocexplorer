package client

import (
	"strings"

	"gitlab.com/vocdoni/vocexplorer/util"
)

// UpdateDashboardInfo calls gateway apis, updates info needed for dashboard page
func UpdateDashboardInfo(c *Client, vc *VochainInfo) {
	UpdateGatewayInfo(c, vc)
	UpdateBlockStatus(c, vc)
	UpdateCounts(c, vc)
	// UpdateVochainProcessList(c, vc)
	// UpdateEntityList(c, vc)
}

// UpdateCounts calls gateway apis, updates total number of processes and entities
func UpdateCounts(c *Client, vc *VochainInfo) {
	procs, err := c.GetProcessCount()
	util.ErrPrint(err)
	entities, err := c.GetEntityCount()
	util.ErrPrint(err)
	vc.ProcessCount = int(procs)
	vc.EntityCount = int(entities)

}

// // UpdateVocDashDashboardInfo calls gateway apis, updates info needed for processes page
// func UpdateVocDashDashboardInfo(c *Client, vc *VochainInfo, index int) {
// 	UpdateVochainProcessList(c, vc, index)
// 	UpdateEntityList(c, vc, index)
// }

// UpdateGatewayInfo calls gateway api, updates vc
func UpdateGatewayInfo(c *Client, vc *VochainInfo) {
	apiList, health, ok, timestamp, err := c.GetGatewayInfo()
	util.ErrPrint(err)
	vc.APIList = apiList
	vc.Ok = ok
	vc.Health = health
	vc.Timestamp = timestamp
}

// UpdateBlockStatus calls gateway api, updates vc
func UpdateBlockStatus(c *Client, vc *VochainInfo) {
	blockTime, blockTimeStamp, height, ok, err := c.GetBlockStatus()
	util.ErrPrint(err)
	vc.BlockTime = blockTime
	vc.BlockTimeStamp = blockTimeStamp
	vc.Height = height
	vc.Ok = ok
}

// // UpdateVochainProcessList calls gateway api, updates vs
// func UpdateVochainProcessList(c *Client, vc *VochainInfo, index int) {
// 	GetIDs(&vc.ProcessIDs, c, func() ([]string, error) {
// 		finals, err := c.GetFinalProcessList(int64(index))
// 		if err != nil {
// 			return finals, err
// 		}
// 		lives, err := c.GetLiveProcessList(int64(index))
// 		return append(finals, lives...), err
// 	})
// }

// // UpdateEntityList calls gateway api, updates vs
// func UpdateEntityList(c *Client, vc *VochainInfo, index int) {
// 	GetIDs(&vc.EntityIDs, c, func() ([]string, error) {
// 		return c.GetScrutinizerEntities(int64(index))
// 	})
// }

// GetIDs gets ids
func GetIDs(IDList *[]string, c *Client, getList func() ([]string, error)) {
	var err error
	*IDList, err = getList()
	util.ErrPrint(err)
}

// UpdateProcessResults updates auxilary info for all currently displayed process id's
func UpdateProcessResults(c *Client, vc *VochainInfo) {
	if vc.ProcessResults == nil {
		vc.ProcessResults = make(map[string]ProcessInfo)
	}
	if vc.EnvelopeHeights == nil {
		vc.EnvelopeHeights = make(map[string]int64)
	}
	for _, ID := range vc.ProcessIDs {
		if ID != "" {
			if _, ok := vc.ProcessResults[ID]; !ok {
				t, st, _, err := c.GetProcessResults(strings.ToLower(ID))
				if !util.ErrPrint(err) {
					vc.ProcessResults[ID] = ProcessInfo{
						ProcessType: t,
						State:       st}
				}
			}
		}
	}
}

// UpdateProcessesDashboardInfo updates process info to include status and recent envelopes
func UpdateProcessesDashboardInfo(c *Client, process *FullProcessInfo, processID string) {
	if process == nil {
		process = new(FullProcessInfo)
	}
	t, st, res, err := c.GetProcessResults(processID)
	if !util.ErrPrint(err) {
		process.ProcessType = t
		process.Results = res
		process.State = st
	}
}

// UpdateAuxEntityInfo updates process info map to include all currently displayed process IDs
func UpdateAuxEntityInfo(c *Client, e *EntityInfo) {
	if e.ProcessTypes == nil {
		e.ProcessTypes = make(map[string]ProcessInfo)
	}
	if e.EnvelopeHeights == nil {
		e.EnvelopeHeights = make(map[string]int64)
	}
	for _, ID := range e.ProcessIDs {
		if ID != "" {
			if _, ok := e.ProcessTypes[ID]; !ok {
				t, st, _, err := c.GetProcessResults(strings.ToLower(ID))
				if !util.ErrPrint(err) {
					e.ProcessTypes[ID] = ProcessInfo{
						ProcessType: t,
						State:       st}
				}
			}
		}
	}
}
