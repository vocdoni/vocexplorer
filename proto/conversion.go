package proto

import (
	"encoding/hex"
	"math/big"

	"github.com/vocdoni/vocexplorer/api/dbtypes"
)

// Mirror returns the mirrored type
func (t *BlockchainInfo) Mirror() *dbtypes.BlockchainInfo {
	return &dbtypes.BlockchainInfo{
		Network:           t.Network,
		Version:           t.Version,
		LatestBlockHeight: t.LatestBlockHeight,
		GenesisTimeStamp:  t.GenesisTimeStamp.AsTime(),
		ChainID:           t.ChainID,
		BlockTime:         t.BlockTime,
		BlockTimeStamp:    t.BlockTimeStamp,
		Height:            t.Height,
		MaxBytes:          t.MaxBytes,
		Syncing:           t.Syncing,
	}
}

// Mirror returns the mirrored type
func (t *Height) Mirror() *dbtypes.Height {
	return &dbtypes.Height{
		Height: t.Height,
	}
}

// Mirror returns the mirrored type
func (t *Envelope) Mirror() *dbtypes.Envelope {
	weight := new(big.Int)
	weight.SetBytes(t.Weight)
	return &dbtypes.Envelope{
		EncryptionKeyIndexes: t.EncryptionKeyIndexes,
		Nullifier:            hex.EncodeToString(t.Nullifier),
		ProcessID:            hex.EncodeToString(t.ProcessID),
		Package:              t.Package,
		ProcessHeight:        t.ProcessHeight,
		GlobalHeight:         t.GlobalHeight,
		TxHeight:             t.TxHeight,
		Weight:               weight.String(),
	}
}

// Mirror returns the mirrored type
func (t *StoreBlock) Mirror() *dbtypes.StoreBlock {
	return &dbtypes.StoreBlock{
		Hash:     t.Hash,
		Height:   t.Height,
		NumTxs:   t.NumTxs,
		Time:     t.Time.AsTime(),
		Proposer: t.Proposer,
	}
}

// Mirror returns the mirrored type
func (t *Transaction) Mirror() *dbtypes.Transaction {
	return &dbtypes.Transaction{
		Height:    t.Height,
		Index:     t.Index,
		Tx:        t.Tx,
		TxHeight:  t.TxHeight,
		Nullifier: hex.EncodeToString(t.Nullifier),
		Hash:      t.Hash,
	}
}

// Mirror returns the mirrored type
func (t *ItemList) Mirror() *dbtypes.ItemList {
	return &dbtypes.ItemList{
		Items: t.Items,
	}
}

// Mirror returns the mirrored type
func (t *Validator) Mirror() *dbtypes.Validator {
	return &dbtypes.Validator{
		Address:          t.Address,
		PubKey:           t.PubKey,
		VotingPower:      t.VotingPower,
		ProposerPriority: t.ProposerPriority,
		Height:           t.Height.Mirror(),
	}
}

// Mirror returns the mirrored type
func (t *Process) Mirror() *dbtypes.Process {
	return &dbtypes.Process{
		ID:          t.ID,
		EntityID:    t.EntityID,
		LocalHeight: t.LocalHeight.Mirror(),
	}
}

// Mirror returns the mirrored type
func (t *HeightMap) Mirror() *dbtypes.HeightMap {
	return &dbtypes.HeightMap{
		Heights: t.Heights,
	}
}
