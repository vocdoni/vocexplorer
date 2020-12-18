package vochain

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/tendermint/tendermint/p2p"
	tmtypes "github.com/tendermint/tendermint/types"
	"gitlab.com/vocdoni/vocexplorer/proto"
	"go.vocdoni.io/dvote/log"
)

// GetStatus gets the basic blockchain status
func (vs *VochainService) GetStatus() *proto.BlockchainInfo {
	genesisTime, err := ptypes.TimestampProto(vs.app.Node.GenesisDoc().GenesisTime)
	if err != nil {
		log.Warn(err)
	}
	status := &proto.BlockchainInfo{
		BlockTime:         vs.info.BlockTimes()[:],
		BlockTimeStamp:    int32(vs.app.State.Header(true).Timestamp),
		ChainID:           vs.app.Node.GenesisDoc().ChainID,
		GenesisTimeStamp:  genesisTime,
		Height:            vs.app.State.Header(true).Height,
		LatestBlockHeight: vs.app.Node.BlockStore().Height(),
		MaxBytes:          vs.app.Node.GenesisDoc().ConsensusParams.Block.MaxBytes,
		Network:           vs.app.Node.NodeInfo().(p2p.DefaultNodeInfo).Network,
		Syncing:           !vs.info.Sync(),
		Version:           vs.app.Node.NodeInfo().(p2p.DefaultNodeInfo).Version,
	}
	return status
}

// GetValidators the list of validators
func (vs *VochainService) GetValidators() ([]*tmtypes.Validator, error) {
	validators := vs.app.Node.ConsensusState().Validators
	return validators.Validators, nil
}
