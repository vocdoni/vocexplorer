package client

import (
	"fmt"
	"strings"

	"go.vocdoni.io/dvote/types"
	"go.vocdoni.io/proto/build/go/models"
)

func (c *Client) GetGatewayInfo() error {
	var req types.MetaRequest
	req.Method = "getInfo"
	resp, err := c.Request(req)
	if err != nil {
		return err
	}
	if !resp.Ok {
		return fmt.Errorf(resp.Message)
	}
	if resp.Health <= 0 {
		return fmt.Errorf("gateway %s health is %d", c.Address, resp.Health)
	}
	if !strings.Contains(strings.Join(resp.APIList, ""), "vote") {
		return fmt.Errorf("gateway %s does not enable vote api", c.Address)
	}
	return nil
}

func (c *Client) GetStats() (*types.VochainStats, error) {
	var req types.MetaRequest
	req.Method = "getStats"
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf(resp.Message)
	}
	return resp.Stats, nil
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

func (c *Client) GetProcessList(entityId []byte, searchTerm string, namespace uint32, status string, withResults bool, from, listSize int) ([]string, error) {
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
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get process list: (%s)", resp.Message)
	}
	if resp.Message == "no results yet" {
		return nil, nil
	}
	return resp.ProcessList, nil
}

func (c *Client) GetProcess(pid []byte) (*types.Process, error) {
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
	return resp.Process, nil
}

func (c *Client) GetProcessKeys(pid []byte) ([]types.Key, []types.Key, []types.Key, []types.Key, error) {
	var req types.MetaRequest
	req.Method = "getProcessKeys"
	req.ProcessID = pid
	resp, err := c.Request(req)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if !resp.Ok {
		return nil, nil, nil, nil, fmt.Errorf(resp.Message)
	}
	return resp.EncryptionPublicKeys, resp.EncryptionPrivKeys, resp.CommitmentKeys, resp.RevealKeys, nil
}

func (c *Client) GetProcessCount(entityId []byte) (int64, error) {
	var req types.MetaRequest
	req.Method = "getProcessCount"
	req.EntityId = entityId
	resp, err := c.Request(req)
	if err != nil {
		return 0, err
	}
	if !resp.Ok {
		return 0, fmt.Errorf(resp.Message)
	}
	return *resp.Size, nil
}

func (c *Client) GetResults(pid []byte) ([][]string, string, string, bool, error) {
	var req types.MetaRequest
	req.Method = "getResults"
	req.ProcessID = pid
	resp, err := c.Request(req)
	if err != nil {
		return nil, "", "", false, err
	}
	if !resp.Ok {
		return nil, "", "", false, fmt.Errorf("cannot get results: (%s)", resp.Message)
	}
	if resp.Message == "no results yet" {
		return nil, resp.State, "", false, nil
	}
	return resp.Results, resp.State, resp.Type, *resp.Final, nil
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

func (c *Client) GetValidatorList() ([]*models.Validator, error) {
	var req types.MetaRequest
	req.Method = "getValidatorList"
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get validator list: (%s)", resp.Message)
	}
	return resp.ValidatorList, nil
}

func (c *Client) GetEnvelope(nullifier []byte) (*types.EnvelopePackage, error) {
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
	resp.Envelope.Meta.Nullifier = nullifier
	return resp.Envelope, nil
}

func (c *Client) GetEnvelopeList(pid []byte, from, listSize int) ([]*types.EnvelopeMetadata, error) {
	var req types.MetaRequest
	req.Method = "getEnvelopeList"
	req.ProcessID = pid
	req.ListSize = listSize
	req.From = from
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get envelope list: (%s)", resp.Message)
	}
	return resp.Envelopes, nil
}

func (c *Client) GetBlock(height uint32) (*types.BlockMetadata, error) {
	var req types.MetaRequest
	req.Method = "getBlock"
	req.Height = height
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get block: (%s)", resp.Message)
	}
	return resp.Block, nil
}

func (c *Client) GetBlockByHash(hash []byte) (*types.BlockMetadata, error) {
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
	return resp.Block, nil
}

func (c *Client) GetBlockList(from, listSize int) ([]*types.BlockMetadata, error) {
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
	return resp.BlockList, nil
}

func (c *Client) GetTx(blockHeight uint32, txIndex int32) (*types.TxPackage, error) {
	var req types.MetaRequest
	req.Method = "getTx"
	req.Height = blockHeight
	req.TxIndex = txIndex
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get tx: (%s)", resp.Message)
	}
	return resp.Tx, nil
}

func (c *Client) GetTxListForBlock(blockHeight uint32, from, listSize int) ([]*types.TxMetadata, error) {
	var req types.MetaRequest
	req.Method = "getTxListForBlock"
	req.Height = blockHeight
	req.From = from
	req.ListSize = listSize
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get tx list for block %d: (%s)", blockHeight, resp.Message)
	}
	return resp.TxList, nil
}
