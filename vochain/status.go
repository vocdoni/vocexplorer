package vochain

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/tendermint/tendermint/p2p"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/vocexplorer/proto"
)

// GetStatus gets the basic blockchain status
func (vs *VochainService) GetStatus() *proto.BlockchainInfo {
	genesisTime, err := ptypes.TimestampProto(vs.app.Node.GenesisDoc().GenesisTime)
	if err != nil {
		log.Warn(err)
	}
	status := &proto.BlockchainInfo{
		BlockTime:         vs.info.BlockTimes()[:],
		BlockTimeStamp:    int32(vs.app.State.Header(true).Time.Unix()),
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
func (vs *VochainService) GetValidators() ([]types.GenesisValidator, error) {
	return vs.scrut.VochainState.Validators(false)
}