package vochain

import (
	"fmt"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

// GetBlock retrieves a single block from the vochain storage
func (vs *VochainService) GetBlock(height int64) (*ctypes.ResultBlock, error) {
	height, err := getHeight(vs.app.Node.BlockStore().Height(), height)
	if err != nil {
		return nil, err
	}
	block := vs.app.Node.BlockStore().LoadBlock(height)
	blockMeta := vs.app.Node.BlockStore().LoadBlockMeta(height)
	if blockMeta == nil {
		return &ctypes.ResultBlock{BlockID: types.BlockID{}, Block: block}, nil
	}
	return &ctypes.ResultBlock{BlockID: blockMeta.BlockID, Block: block}, nil
}

// latestHeight can be either latest committed or uncommitted (+1) height.
func getHeight(latestHeight int64, height int64) (int64, error) {
	if height <= 0 {
		return 0, fmt.Errorf("height must be greater than 0, but got %d", height)
	}
	if height > latestHeight {
		return 0, fmt.Errorf("height %d must be less than or equal to the current blockchain height %d",
			height, latestHeight)
	}
	return height, nil
}
