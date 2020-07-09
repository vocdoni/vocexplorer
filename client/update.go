package client

import (
	"gitlab.com/vocdoni/vocexplorer/util"
)

// UpdateVochainInfo calls gateway apis, updates vs
func UpdateVochainInfo(c *Client, vc *VochainInfo) {
	UpdateGatewayInfo(c, vc)
}

// UpdateGatewayInfo calls gateway api, updates vc
func UpdateGatewayInfo(c *Client, vc *VochainInfo) {
	resp, err := c.GetGatewayInfo()
	util.ErrPrint(err)
	vc.APIList = resp.APIList
	vc.Ok = resp.Ok
	vc.Health = resp.Health
	vc.Timestamp = resp.Timestamp
}
