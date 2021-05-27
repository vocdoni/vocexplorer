package client

import (
	"fmt"
	"strings"

	"gitlab.com/vocdoni/vocexplorer/util"
	"go.vocdoni.io/dvote/api"
	"go.vocdoni.io/dvote/vochain/scrutinizer/indexertypes"
	"go.vocdoni.io/proto/build/go/models"
)

func (c *Client) GetGatewayInfo() error {
	var req api.MetaRequest
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
	apiList := strings.Join(resp.APIList, "")
	if !strings.Contains(apiList, "vote") {
		return fmt.Errorf("gateway %s does not enable vote api", c.Address)
	}
	if !strings.Contains(apiList, "indexer") {
		return fmt.Errorf("gateway %s does not enable indexer api", c.Address)
	}
	return nil
}

func (c *Client) GetStats() (*api.VochainStats, error) {
	var req api.MetaRequest
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
	var req api.MetaRequest
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
	var req api.MetaRequest
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
	var req api.MetaRequest
	req.Method = "getProcessList"
	req.EntityId = entityId
	req.SearchTerm = searchTerm
	req.Namespace = namespace
	req.Status = status
	req.WithResults = withResults
	req.From = util.Max(from, 0)
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

func (c *Client) GetProcess(pid []byte) (*indexertypes.Process, error) {
	var req api.MetaRequest
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

func (c *Client) GetProcessSummary(pid []byte) (*api.ProcessSummary, error) {
	var req api.MetaRequest
	req.Method = "getProcessSummary"
	req.ProcessID = pid
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf(resp.Message)
	}
	// Remove need to null-check envelope height
	if resp.ProcessSummary != nil {
		if resp.ProcessSummary.EnvelopeHeight == nil {
			resp.ProcessSummary.EnvelopeHeight = new(uint32)
		}
	}
	return resp.ProcessSummary, nil
}

func (c *Client) GetProcessKeys(pid []byte) ([]api.Key, []api.Key, []api.Key, []api.Key, error) {
	var req api.MetaRequest
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
	var req api.MetaRequest
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
	var req api.MetaRequest
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
		return nil, resp.State, resp.Type, false, nil
	}
	return resp.Results, resp.State, resp.Type, *resp.Final, nil
}

// r.registerPublic("getResultsWeight", r.getResultsWeight)

func (c *Client) GetEntityList(searchTerm string, listSize, from int) ([]string, error) {
	var req api.MetaRequest
	req.Method = "getEntityList"
	req.SearchTerm = searchTerm
	req.ListSize = listSize
	req.From = util.Max(from, 0)
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
	var req api.MetaRequest
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
	var req api.MetaRequest
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

func (c *Client) GetEnvelope(nullifier []byte) (*indexertypes.EnvelopePackage, error) {
	var req api.MetaRequest
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

func (c *Client) GetEnvelopeList(pid []byte, from, listSize int, searchTerm string) ([]*indexertypes.EnvelopeMetadata, error) {
	var req api.MetaRequest
	req.Method = "getEnvelopeList"
	req.ProcessID = pid
	req.SearchTerm = searchTerm
	req.ListSize = listSize
	req.From = util.Max(from, 0)
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get envelope list: (%s)", resp.Message)
	}
	return resp.Envelopes, nil
}

func (c *Client) GetBlock(height uint32) (*indexertypes.BlockMetadata, error) {
	var req api.MetaRequest
	req.Method = "getBlock"
	req.Height = height
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get block: (%s)", resp.Message)
	}
	if resp.Block.Height == 0 {
		resp.Block.Height = height
	}
	return resp.Block, nil
}

func (c *Client) GetBlockByHash(hash []byte) (*indexertypes.BlockMetadata, error) {
	var req api.MetaRequest
	req.Method = "getBlockByHash"
	req.Hash = hash
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get block: (%s)", resp.Message)
	}
	if len(resp.Block.Hash) == 0 {
		resp.Block.Hash = hash
	}
	return resp.Block, nil
}

func (c *Client) GetBlockList(from, listSize int) ([]*indexertypes.BlockMetadata, error) {
	var req api.MetaRequest
	req.Method = "getBlockList"
	req.From = util.Max(from, 0)
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

func (c *Client) GetTx(blockHeight uint32, txIndex int32) (*indexertypes.TxPackage, error) {
	var req api.MetaRequest
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
	if resp.Tx.BlockHeight == 0 {
		resp.Tx.BlockHeight = blockHeight
	}
	return resp.Tx, nil
}

func (c *Client) GetTxListForBlock(blockHeight uint32, from, listSize int) ([]*indexertypes.TxMetadata, error) {
	var req api.MetaRequest
	req.Method = "getTxListForBlock"
	req.Height = blockHeight
	req.From = util.Max(from, 0)
	req.ListSize = listSize
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get tx list for block %d: (%s)", blockHeight, resp.Message)
	}
	for _, tx := range resp.TxList {
		if tx.BlockHeight == 0 {
			tx.BlockHeight = blockHeight
		}
	}
	return resp.TxList, nil
}
