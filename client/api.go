package client

import (
	"fmt"

	"go.vocdoni.io/dvote/types"
	"go.vocdoni.io/proto/build/go/models"
	"google.golang.org/protobuf/proto"
)

func (c *Client) GetStats() (*models.VochainStats, error) {
	var req types.MetaRequest
	req.Method = "getStats"
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf(resp.Message)
	}
	stats := new(models.VochainStats)
	if err := proto.Unmarshal(resp.Content, stats); err != nil {
		return nil, err
	}
	return stats, nil
}

func (c *Client) GetEnvelopeHeight(pid []byte) (uint32, error) {
	var req types.MetaRequest
	req.Method = "getEnvelopeHeight"
	req.ProcessID = pid
	resp, err := c.Request(req)
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
	resp, err := c.Request(req)
	if err != nil {
		return nil, nil, 0, err
	}
	if !resp.Ok {
		return nil, nil, 0, fmt.Errorf(resp.Message)
	}
	return resp.BlockTime, resp.Height, resp.BlockTimestamp, nil
}

func (c *Client) GetProcessList(entityId []byte, searchTerm string, namespace uint32, status string, withResults bool, from, listSize int) ([]string, int64, error) {
	var req types.MetaRequest
	req.Method = "getProcessList"
	req.EntityId = entityId
	req.SearchTerm = searchTerm
	req.Namespace = namespace
	req.Status = status
	req.WithResults = withResults
	req.From = from
	req.ListSize = listSize
	resp, err := c.Request(req)
	if err != nil {
		return nil, 0, err
	}
	if !resp.Ok {
		return nil, 0, fmt.Errorf("cannot get process list: (%s)", resp.Message)
	}
	if resp.Message == "no results yet" {
		return nil, 0, nil
	}
	return resp.ProcessList, *resp.Size, nil
}

func (c *Client) GetProcess(pid []byte) (*models.Process, error) {
	var req types.MetaRequest
	req.Method = "getProcessInfo"
	req.ProcessID = pid
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf(resp.Message)
	}
	return resp.ProcessInfo.(*models.Process), nil
}

func (c *Client) GetProcessCount() (int64, error) {
	var req types.MetaRequest
	req.Method = "getProcessInfo"
	resp, err := c.Request(req)
	if err != nil {
		return 0, err
	}
	if !resp.Ok {
		return 0, fmt.Errorf(resp.Message)
	}
	return *resp.Size, nil
}

func (c *Client) GetResults(pid []byte) ([][]string, string, bool, error) {
	var req types.MetaRequest
	req.Method = "getResults"
	req.ProcessID = pid
	resp, err := c.Request(req)
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

// r.registerPublic("getResultsWeight", r.getResultsWeight)

func (c *Client) GetEntityList(searchTerm string, listSize, from int) ([]string, error) {
	var req types.MetaRequest
	req.Method = "getEntityList"
	req.SearchTerm = searchTerm
	req.ListSize = listSize
	req.From = from
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get entity list: (%s)", resp.Message)
	}
	return resp.EntityIDs, nil
}

func (c *Client) GetEntityCount() (int64, error) {
	var req types.MetaRequest
	req.Method = "getEntityCount"
	resp, err := c.Request(req)
	if err != nil {
		return 0, err
	}
	if !resp.Ok {
		return 0, fmt.Errorf("cannot get entity count: (%s)", resp.Message)
	}
	return *resp.Size, nil
}

func (c *Client) GetValidatorList() (*models.ValidatorList, error) {
	var req types.MetaRequest
	req.Method = "getValidatorList"
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get validator list: (%s)", resp.Message)
	}
	list := new(models.ValidatorList)
	if err := proto.Unmarshal(resp.ValidatorList, list); err != nil {
		return nil, err
	}
	return list, nil
}

func (c *Client) GetEnvelope(nullifier []byte) (*models.EnvelopePackage, error) {
	var req types.MetaRequest
	req.Method = "getEnvelope"
	req.Nullifier = nullifier
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get envelope: (%s)", resp.Message)
	}
	envelope := new(models.EnvelopePackage)
	if err := proto.Unmarshal(resp.Content, envelope); err != nil {
		return nil, err
	}
	return envelope, nil
}

func (c *Client) GetEnvelopeList(pid []byte, listSize int) (*models.EnvelopePackageList, error) {
	var req types.MetaRequest
	req.Method = "getEnvelopeList"
	req.ProcessID = pid
	req.ListSize = listSize
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get envelope list: (%s)", resp.Message)
	}
	list := new(models.EnvelopePackageList)
	if err := proto.Unmarshal(resp.Content, list); err != nil {
		return nil, err
	}
	return list, nil
}

func (c *Client) GetBlock(height uint32) (*models.TendermintHeader, error) {
	var req types.MetaRequest
	req.Method = "getBlock"
	req.BlockHeight = height
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get block: (%s)", resp.Message)
	}
	block := new(models.TendermintHeader)
	if err := proto.Unmarshal(resp.Content, block); err != nil {
		return nil, err
	}
	return block, nil
}

func (c *Client) GetBlockByHash(hash []byte) (*models.TendermintHeader, error) {
	var req types.MetaRequest
	req.Method = "getBlockByHash"
	req.Payload = hash
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get block: (%s)", resp.Message)
	}
	block := new(models.TendermintHeader)
	if err := proto.Unmarshal(resp.Content, block); err != nil {
		return nil, err
	}
	return block, nil
}

func (c *Client) GetBlockList(from, listSize int) (*models.TendermintHeaderList, error) {
	var req types.MetaRequest
	req.Method = "getBlockList"
	req.From = from
	req.ListSize = listSize
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get block: (%s)", resp.Message)
	}
	list := new(models.TendermintHeaderList)
	if err := proto.Unmarshal(resp.Content, list); err != nil {
		return nil, err
	}
	return list, nil
}

func (c *Client) GetTx(blockHeight uint32, txIndex int32) (*models.SignedTx, error) {
	var req types.MetaRequest
	req.Method = "getTx"
	req.BlockHeight = blockHeight
	req.TxIndex = txIndex
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get tx: (%s)", resp.Message)
	}
	tx := new(models.SignedTx)
	if err := proto.Unmarshal(resp.Content, tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (c *Client) GetTxListForBlock(blockHeight uint32) (*models.SignedTxList, error) {
	var req types.MetaRequest
	req.Method = "getTxListForBlock"
	req.BlockHeight = blockHeight
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get tx list for block %d: (%s)", blockHeight, resp.Message)
	}
	txList := new(models.SignedTxList)
	if err := proto.Unmarshal(resp.Content, txList); err != nil {
		return nil, err
	}
	return txList, nil
}
