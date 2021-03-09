package vochain

import (
	"github.com/tendermint/tendermint/p2p"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/vocdoni/vocexplorer/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetStatus gets the basic blockchain status
func (vs *VochainService) GetStatus() *proto.BlockchainInfo {
	return &proto.BlockchainInfo{
		BlockTime:         vs.info.BlockTimes()[:],
		BlockTimeStamp:    int32(vs.app.State.Header(true).Timestamp),
		ChainID:           vs.app.Node.GenesisDoc().ChainID,
		GenesisTimeStamp:  timestamppb.New(vs.app.Node.GenesisDoc().GenesisTime),
		Height:            vs.app.State.Header(true).Height,
		LatestBlockHeight: vs.app.Node.BlockStore().Height(),
		MaxBytes:          vs.app.Node.GenesisDoc().ConsensusParams.Block.MaxBytes,
		Network:           vs.app.Node.NodeInfo().(p2p.DefaultNodeInfo).Network,
		Syncing:           !vs.info.Sync(),
		Version:           vs.app.Node.NodeInfo().(p2p.DefaultNodeInfo).Version,
	}
}

// GetValidators the list of validators
func (vs *VochainService) GetValidators() ([]*tmtypes.Validator, error) {
	validators := vs.app.Node.ConsensusState().Validators
	return validators.Validators, nil
}
