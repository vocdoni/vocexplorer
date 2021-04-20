
func (c *Client) GetEnvelopeHeight(pid []byte) (uint32, error) {
	var req types.MetaRequest
	req.Method = "getEnvelopeHeight"
	req.ProcessID = pid
	resp, err := c.Request(req, nil)
	if err != nil {
		return 0, err
	}
	if !resp.Ok {
		return 0, fmt.Errorf(resp.Message)
	}
	return *resp.Height, nil
}
func (c *Client) GetBlockStatus() (*[5]int32, *uint32, int32, error) {
	var req types.MetaRequest
	req.Method = "getBlockStatus"
	req.ProcessID = pid
	resp, err := c.Request(req, nil)
	if err != nil {
		return 0, err
	}
	if !resp.Ok {
		return 0, fmt.Errorf(resp.Message)
	}
	return *resp.BlockTime, *resp.Height, *resp.BlockTimeStamp
}

// r.registerPublic("getProcessList", r.getProcessList)
// r.registerPublic("getProcessInfo", r.getProcessInfo)
// r.registerPublic("getProcessCount", r.getProcessCount)
// r.registerPublic("getResults", r.getResults)
// r.registerPublic("getResultsWeight", r.getResultsWeight)
// r.registerPublic("getEntityList", r.getEntityList)
// r.registerPublic("getEntityCount", r.getEntityCount)
// r.registerPublic("getValidatorList", r.getValidatorList)
// r.registerPublic("getEnvelope", r.getEnvelope)
// r.registerPublic("getEnvelopeList", r.getEnvelopeList)
// r.registerPublic("getBlock", r.getBlock)
// r.registerPublic("getBlockByHash", r.getBlockByHash)
// r.registerPublic("getBlockList", r.getBlockList)
// r.registerPublic("getTx", r.getTx)
// r.registerPublic("getTxListForBlock", r.getTxListForBlock)

func (c *Client) GetProcessList(pid []byte) ([][]string, string, bool, error) {
	var req types.MetaRequest
	req.Method = "getResults"
	req.ProcessID = pid
	resp, err := c.Request(req, nil)
	if err != nil {
		return nil, "", false, err
	}
	if !resp.Ok {
		return nil, "", false, fmt.Errorf("cannot get results: (%s)", resp.Message)
	}
	if resp.Message == "no results yet" {
		return nil, resp.State, false, nil
	}
	return resp.Results, resp.State, *resp.Final, nil
}
func (c *Client) GetResults(pid []byte) ([][]string, string, bool, error) {
	var req types.MetaRequest
	req.Method = "getResults"
	req.ProcessID = pid
	resp, err := c.Request(req, nil)
	if err != nil {
		return nil, "", false, err
	}
	if !resp.Ok {
		return nil, "", false, fmt.Errorf("cannot get results: (%s)", resp.Message)
	}
	if resp.Message == "no results yet" {
		return nil, resp.State, false, nil
	}
	return resp.Results, resp.State, *resp.Final, nil
}
func (c *Client) GetResults(pid []byte) ([][]string, string, bool, error) {
	var req types.MetaRequest
	req.Method = "getResults"
	req.ProcessID = pid
	resp, err := c.Request(req, nil)
	if err != nil {
		return nil, "", false, err
	}
	if !resp.Ok {
		return nil, "", false, fmt.Errorf("cannot get results: (%s)", resp.Message)
	}
	if resp.Message == "no results yet" {
		return nil, resp.State, false, nil
	}
	return resp.Results, resp.State, *resp.Final, nil
}