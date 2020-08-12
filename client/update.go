package client

import (
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
	vc.ProcessCount = procs
	vc.EntityCount = entities

}

// UpdateVocDashDashboardInfo calls gateway apis, updates info needed for processes page
func UpdateVocDashDashboardInfo(c *Client, vc *VochainInfo, index int) {
	UpdateVochainProcessList(c, vc, index)
	UpdateEntityList(c, vc, index)
}

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

// UpdateVochainProcessList calls gateway api, updates vs
func UpdateVochainProcessList(c *Client, vc *VochainInfo, index int) {
	GetIDs(&vc.ProcessIDs, c, func() ([]string, error) {
		finals, err := c.GetFinalProcessList(int64(index))
		if err != nil {
			return finals, err
		}
		lives, err := c.GetLiveProcessList(int64(index))
		return append(finals, lives...), err
	})
}

// UpdateEntityList calls gateway api, updates vs
func UpdateEntityList(c *Client, vc *VochainInfo, index int) {
	GetIDs(&vc.EntityIDs, c, func() ([]string, error) {
		return c.GetScrutinizerEntities(int64(index))
	})
}

// GetIDs gets ids
func GetIDs(IDList *[]string, c *Client, getList func() ([]string, error)) {
	var err error
	*IDList, err = getList()
	util.ErrPrint(err)
}

// UpdateAuxProcessInfo updates auxilary info for all currently displayed process id's
func UpdateAuxProcessInfo(c *Client, vc *VochainInfo) {
	if vc.ProcessSearchList == nil {
		vc.ProcessSearchList = make(map[string]ProcessInfo)
	}
	if vc.EnvelopeHeights == nil {
		vc.EnvelopeHeights = make(map[string]int64)
	}
	// If all processes are populated, send no requests. Process results are not updated without page refresh.
	if len(vc.ProcessSearchList) >= len(vc.ProcessIDs) && len(vc.EnvelopeHeights) >= len(vc.ProcessIDs) {
		return
	}
	numReq := 0
	for _, ID := range vc.ProcessSearchIDs {
		if _, ok := vc.ProcessSearchList[ID]; !ok {
			t, st, _, err := c.GetProcessResults(ID)
			if !util.ErrPrint(err) {
				vc.ProcessSearchList[ID] = ProcessInfo{
					ProcessType: t,
					State:       st}
			}
			numReq++
		}
		if _, ok := vc.EnvelopeHeights[ID]; !ok {
			height, err := c.GetEnvelopeHeight(ID)
			if !util.ErrPrint(err) {
				vc.EnvelopeHeights[ID] = height
			}
			numReq++
		}
	}
	// If currently-displayed processes are populated, start to populate ones which could be displayed
	// This reduces load time & allows for type/state search.
	for _, ID := range vc.ProcessIDs {
		if numReq >= 20 {
			break
		}
		if _, ok := vc.ProcessSearchList[ID]; !ok {
			t, st, _, err := c.GetProcessResults(ID)
			if !util.ErrPrint(err) {
				vc.ProcessSearchList[ID] = ProcessInfo{
					ProcessType: t,
					State:       st}
			}
			numReq++
		}
		if _, ok := vc.EnvelopeHeights[ID]; !ok {
			height, err := c.GetEnvelopeHeight(ID)
			if !util.ErrPrint(err) {
				vc.EnvelopeHeights[ID] = height
			}
			numReq++
		}

	}
}

// UpdateProcessesDashboardInfo updates process info to include status and recent envelopes
func UpdateProcessesDashboardInfo(c *Client, process *FullProcessInfo, processID string, index int) {
	if process == nil {
		process = new(FullProcessInfo)
	}
	t, st, res, err := c.GetProcessResults(processID)
	if !util.ErrPrint(err) {
		process.ProcessType = t
		process.Results = res
		process.State = st
	}
	GetIDs(&process.Nullifiers, c, func() ([]string, error) {
		return c.GetEnvelopeList(processID, int64(index))
	})
}

// UpdateEntitiesDashboardInfo updates entity info to include recent processes
func UpdateEntitiesDashboardInfo(c *Client, entity *EntityInfo, entityID string, index int) {
	if entity == nil {
		entity = new(EntityInfo)
	}
	GetIDs(&entity.ProcessIDs, c, func() ([]string, error) {
		return c.GetProcessList(entityID, int64(index))
	})
}

// UpdateAuxEntityInfo updates process info map to include all currently displayed process IDs
func UpdateAuxEntityInfo(c *Client, e *EntityInfo) {
	if e.ProcessList == nil {
		e.ProcessList = make(map[string]ProcessInfo)
	}
	if e.EnvelopeHeights == nil {
		e.EnvelopeHeights = make(map[string]int64)
	}
	// If all processes are populated, send no requests. Process results are not updated without page refresh.
	if len(e.ProcessList) >= len(e.ProcessIDs) && len(e.EnvelopeHeights) >= len(e.ProcessIDs) {
		return
	}
	numReq := 0
	for _, ID := range e.ProcessSearchIDs {
		if _, ok := e.ProcessList[ID]; !ok {
			t, st, _, err := c.GetProcessResults(ID)
			if !util.ErrPrint(err) {
				e.ProcessList[ID] = ProcessInfo{
					ProcessType: t,
					State:       st}
			}
			numReq++
		}
		if _, ok := e.EnvelopeHeights[ID]; !ok {
			height, err := c.GetEnvelopeHeight(ID)
			if !util.ErrPrint(err) {
				e.EnvelopeHeights[ID] = height
			}
			numReq++
		}
	}
	// If currently-displayed processes are populated, start to populate ones which could be displayed
	// This reduces load time & allows for type/state search.
	for _, ID := range e.ProcessIDs {
		if numReq >= 20 {
			break
		}
		if _, ok := e.ProcessList[ID]; !ok {
			t, st, _, err := c.GetProcessResults(ID)
			if !util.ErrPrint(err) {
				e.ProcessList[ID] = ProcessInfo{
					ProcessType: t,
					State:       st}
			}
			numReq++
		}
		if _, ok := e.EnvelopeHeights[ID]; !ok {
			height, err := c.GetEnvelopeHeight(ID)
			if !util.ErrPrint(err) {
				e.EnvelopeHeights[ID] = height
			}
			numReq++
		}

	}
}
