package proto

import (
	"github.com/golang/protobuf/ptypes"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/api/dbtypes"
)

// Mirror returns the mirrored type
func (t *BlockchainInfo) Mirror() *dbtypes.BlockchainInfo {
	time, err := ptypes.Timestamp(t.GenesisTimeStamp)
	if err != nil {
		log.Error(err)
	}
	m := &dbtypes.BlockchainInfo{
		Network:           t.Network,
		Version:           t.Version,
		LatestBlockHeight: t.LatestBlockHeight,
		GenesisTimeStamp:  time,
		ChainID:           t.ChainID,
		BlockTime:         t.BlockTime,
		BlockTimeStamp:    t.BlockTimeStamp,
		Height:            t.Height,
		MaxBytes:          t.MaxBytes,
		Syncing:           t.Syncing,
	}
	return m
}

// Mirror returns the mirrored type
func (t *Height) Mirror() *dbtypes.Height {
	m := &dbtypes.Height{
		Height: t.Height,
	}
	return m
}

// Mirror returns the mirrored type
func (t *Envelope) Mirror() *dbtypes.Envelope {
	m := &dbtypes.Envelope{
		EncryptionKeyIndexes: t.EncryptionKeyIndexes,
		Nullifier:            t.Nullifier,
		ProcessID:            t.ProcessID,
		Package:              t.Package,
		ProcessHeight:        t.ProcessHeight,
		GlobalHeight:         t.GlobalHeight,
		TxHeight:             t.TxHeight,
	}
	return m
}

// Mirror returns the mirrored type
func (t *StoreBlock) Mirror() *dbtypes.StoreBlock {
	time, err := ptypes.Timestamp(t.Time)
	if err != nil {
		log.Error(err)
	}
	m := &dbtypes.StoreBlock{
		Hash:     t.Hash,
		Height:   t.Height,
		NumTxs:   t.NumTxs,
		Time:     time,
		Proposer: t.Proposer,
	}
	return m
}

// Mirror returns the mirrored type
func (t *Transaction) Mirror() *dbtypes.Transaction {
	m := &dbtypes.Transaction{
		Height:    t.Height,
		Index:     t.Index,
		Tx:        t.Tx,
		TxHeight:  t.TxHeight,
		Nullifier: t.Nullifier,
		Hash:      t.Hash,
	}
	return m
}

// Mirror returns the mirrored type
func (t *ItemList) Mirror() *dbtypes.ItemList {
	m := &dbtypes.ItemList{
		Items: t.Items,
	}
	return m
}

// Mirror returns the mirrored type
func (t *Validator) Mirror() *dbtypes.Validator {
	m := &dbtypes.Validator{
		Address:          t.Address,
		PubKey:           t.PubKey,
		VotingPower:      t.VotingPower,
		ProposerPriority: t.ProposerPriority,
		Height:           t.Height.Mirror(),
	}
	return m
}

// Mirror returns the mirrored type
func (t *Process) Mirror() *dbtypes.Process {
	m := &dbtypes.Process{
		ID:          t.ID,
		EntityID:    t.EntityID,
		LocalHeight: t.LocalHeight.Mirror(),
	}
	return m
}

// Mirror returns the mirrored type
func (t *HeightMap) Mirror() *dbtypes.HeightMap {
	m := &dbtypes.HeightMap{
		Heights: t.Heights,
	}
	return m
}
