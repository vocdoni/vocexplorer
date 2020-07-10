package client

import (
	"gitlab.com/vocdoni/vocexplorer/util"
)

// UpdateVochainInfo calls gateway apis, updates vs
func UpdateVochainInfo(c *Client, vc *VochainInfo) {
	UpdateGatewayInfo(c, vc)
	UpdateBlockStatus(c, vc)
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
	for {
		tempList, err := getList(lastID)
		util.ErrPrint(err)
		*IDList = append(*IDList, tempList...)
		// Repeat if request was full, make sure never gets stuck if fromID is not working
		if len(tempList) < 64 || tempList[len(tempList)-1] == lastID {
			break
		}
		lastID = tempList[len(tempList)-1]
		// fmt.Println("last ID " + lastID)

	}
}
