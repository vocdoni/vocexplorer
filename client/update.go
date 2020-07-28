package client

import (
	"gitlab.com/vocdoni/vocexplorer/util"
)

// UpdateDashboardInfo calls gateway apis, updates info needed for dashboard page
func UpdateDashboardInfo(c *Client, vc *VochainInfo) {
	UpdateGatewayInfo(c, vc)
	UpdateBlockStatus(c, vc)
	// UpdateVochainProcessList(c, vc)
	// UpdateEntityList(c, vc)
}

// UpdateVocDashDashboardInfo calls gateway apis, updates info needed for processes page
func UpdateVocDashDashboardInfo(c *Client, vc *VochainInfo) {
	UpdateVochainProcessList(c, vc)
	UpdateEntityList(c, vc)
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
func UpdateVochainProcessList(c *Client, vc *VochainInfo) {
	GetAllIDs(&vc.ProcessIDs, c, func(fromID string) ([]string, error) {
		return c.GetFinalProcessList(fromID)
	})

	GetAllIDs(&vc.ProcessIDs, c, func(fromID string) ([]string, error) {
		return c.GetLiveProcessList(fromID)
	})
}

// UpdateEntityList calls gateway api, updates vs
func UpdateEntityList(c *Client, vc *VochainInfo) {
	GetAllIDs(&vc.EntityIDs, c, func(fromID string) ([]string, error) {
		return c.GetScrutinizerEntities(fromID)
	})
}

// GetAllIDs iteratively calls getList until all IDs have been collected and stored in IDList
func GetAllIDs(IDList *[]string, c *Client, getList func(string) ([]string, error)) {
	lastID := ""
	if len(*IDList) > 0 {
		lastID = (*IDList)[len(*IDList)-1]

		/*THIS RETURN BREAKS THE UPDATING OF ENTITY AND PROCESS IDS.
		 *statement is here because of bug(?) or confusion: fromID field seems not to be working for these calls.
		 *Each call returns IDs from the beginning of the ID list, regardless of the fromID field.
		 *This means 'lastID' does nothing, so list keeps updating with duplicate ids.
		 */

		return
	}
	for {
		tempList, err := getList(lastID)
		util.ErrPrint(err)
		if len(tempList) <= 0 {
			break
		}
		if tempList[len(tempList)-1] == lastID {
			break
		}
		*IDList = append(*IDList, tempList...)
		// Repeat if request was full, make sure never gets stuck if fromID is not working
		if len(tempList) < 64 || tempList[len(tempList)-1] == lastID {
			break
		}
		lastID = tempList[len(tempList)-1]
		// fmt.Println("last ID " + lastID)

	}
}

// UpdateProcessEnvelopeHeights updates envelope height map to include all current process IDs
func UpdateProcessEnvelopeHeights(c *Client, vc *VochainInfo) {
	if vc.EnvelopeHeights == nil {
		vc.EnvelopeHeights = make(map[string]int64)
	}
	for _, ID := range vc.ProcessSearchIDs {
		if _, ok := vc.EnvelopeHeights[ID]; !ok {
			height, err := c.GetEnvelopeHeight(ID)
			if !util.ErrPrint(err) {
				vc.EnvelopeHeights[ID] = height
			}
		}
	}
}

// UpdateProcessSearchInfo updates process search info map to include all currently displayed process IDs
func UpdateProcessSearchInfo(c *Client, vc *VochainInfo) {
	if vc.ProcessSearchList == nil {
		vc.ProcessSearchList = make(map[string]ProcessInfo)
	}
	// If all processes are populated, send no requests. Process results are not updated without page refresh.
	if len(vc.ProcessSearchList) >= len(vc.ProcessIDs) {
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
	}
	// If currently-displayed processes are populated, start to populate ones which could be displayed
	// This reduces load time & allows for type/state search.
	for _, ID := range vc.ProcessIDs {
		if numReq >= 10 {
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
	}
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
	GetAllIDs(&process.Nullifiers, c, func(fromID string) ([]string, error) {
		return c.GetEnvelopeList(processID, fromID)
	})
}

// UpdateEntitiesDashboardInfo updates entity info to include recent processes
func UpdateEntitiesDashboardInfo(c *Client, entity *EntityInfo, entityID string) {
	if entity == nil {
		entity = new(EntityInfo)
	}
	GetAllIDs(&entity.ProcessIDs, c, func(fromID string) ([]string, error) {
		return c.GetProcessList(entityID, fromID)
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
