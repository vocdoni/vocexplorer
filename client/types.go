package client

import (
	"time"

	"go.vocdoni.io/dvote/types"
	"go.vocdoni.io/dvote/vochain/scrutinizer/indexertypes"
	"go.vocdoni.io/proto/build/go/models"
)

// APIrequest contains all of the possible request fields.
// Fields must be in alphabetical order
type APIrequest struct {
	CensusID     string             `json:"censusId,omitempty"`
	CensusURI    string             `json:"censusUri,omitempty"`
	CensusKey    []byte             `json:"censusKey,omitempty"`
	CensusKeys   [][]byte           `json:"censusKeys,omitempty"`
	CensusValue  []byte             `json:"censusValue,omitempty"`
	CensusDump   []byte             `json:"censusDump,omitempty"`
	CensusType   models.Census_Type `json:"censusType,omitempty"`
	Content      []byte             `json:"content,omitempty"`
	Digested     bool               `json:"digested,omitempty"`
	EntityId     types.HexBytes     `json:"entityId,omitempty"`
	Hash         types.HexBytes     `json:"hash,omitempty"`
	Height       uint32             `json:"height,omitempty"`
	ID           uint32             `json:"id,omitempty"`
	From         int                `json:"from,omitempty"`
	ListSize     int                `json:"listSize,omitempty"`
	Method       string             `json:"method"`
	Name         string             `json:"name,omitempty"`
	Namespace    uint32             `json:"namespace,omitempty"`
	Nullifier    types.HexBytes     `json:"nullifier,omitempty"`
	Payload      []byte             `json:"payload,omitempty"`
	ProcessID    types.HexBytes     `json:"processId,omitempty"`
	ProofData    types.HexBytes     `json:"proofData,omitempty"`
	PubKeys      []string           `json:"pubKeys,omitempty"`
	RootHash     types.HexBytes     `json:"rootHash,omitempty"`
	SearchTerm   string             `json:"searchTerm,omitempty"`
	Signature    types.HexBytes     `json:"signature,omitempty"`
	SrcNetId     string             `json:"sourceNetworkId,omitempty"`
	Status       string             `json:"status,omitempty"`
	Timestamp    int32              `json:"timestamp"`
	TxIndex      int32              `json:"txIndex,omitempty"`
	Type         string             `json:"type,omitempty"`
	URI          string             `json:"uri,omitempty"`
	Weight       *types.BigInt      `json:"weight,omitempty"`
	Weights      []*types.BigInt    `json:"weights,omitempty"`
	WithResults  bool               `json:"withResults,omitempty"`
	VoterAddress types.HexBytes     `json:"voterAddress,omitempty"`
}

// APIresponse contains all of the possible request fields.
// Fields must be in alphabetical order
// Those fields with valid zero-values (such as bool) must be pointers
type APIresponse struct {
	APIList              []string                         `json:"apiList,omitempty"`
	Balance              *uint64                          `json:"balance,omitempty"`
	Block                *indexertypes.BlockMetadata      `json:"block,omitempty"`
	BlockList            []*indexertypes.BlockMetadata    `json:"blockList,omitempty"`
	BlockTime            *[5]int32                        `json:"blockTime,omitempty"`
	BlockTimestamp       int32                            `json:"blockTimestamp,omitempty"`
	CensusID             string                           `json:"censusId,omitempty"`
	CensusList           []string                         `json:"censusList,omitempty"`
	CensusKeys           [][]byte                         `json:"censusKeys,omitempty"`
	CensusDump           []byte                           `json:"censusDump,omitempty"`
	CensusValue          []byte                           `json:"censusValue,omitempty"`
	ChainID              string                           `json:"chainId,omitempty"`
	Content              []byte                           `json:"content,omitempty"`
	Amount               *uint64                          `json:"amount,omitempty"`
	CreationTime         int64                            `json:"creationTime,omitempty"`
	Delegates            []string                         `json:"delegates,omitempty"`
	EncryptionPrivKeys   []Key                            `json:"encryptionPrivKeys,omitempty"`
	EncryptionPublicKeys []Key                            `json:"encryptionPubKeys,omitempty"`
	EntityID             string                           `json:"entityId,omitempty"`
	EntityIDs            []string                         `json:"entityIds,omitempty"`
	Envelope             *indexertypes.EnvelopePackage    `json:"envelope,omitempty"`
	Envelopes            []*indexertypes.EnvelopeMetadata `json:"envelopes,omitempty"`
	Files                []byte                           `json:"files,omitempty"`
	Final                *bool                            `json:"final,omitempty"`
	Finished             *bool                            `json:"finished,omitempty"`
	Health               int32                            `json:"health,omitempty"`
	Height               *uint32                          `json:"height,omitempty"`
	InvalidClaims        []int                            `json:"invalidClaims,omitempty"`
	InfoURI              string                           `json:"infoURI,omitempty"`
	Message              string                           `json:"message,omitempty"`
	Nonce                *uint32                          `json:"nonce,omitempty"`
	Nullifier            string                           `json:"nullifier,omitempty"`
	Nullifiers           *[]string                        `json:"nullifiers,omitempty"`
	Ok                   bool                             `json:"ok"`
	Paused               *bool                            `json:"paused,omitempty"`
	Payload              string                           `json:"payload,omitempty"`
	ProcessSummary       *ProcessSummary                  `json:"processSummary,omitempty"`
	ProcessID            types.HexBytes                   `json:"processId,omitempty"`
	ProcessIDs           []string                         `json:"processIds,omitempty"`
	Process              *indexertypes.Process            `json:"process,omitempty"`
	ProcessList          []string                         `json:"processList,omitempty"`
	Registered           *bool                            `json:"registered,omitempty"`
	Request              string                           `json:"request"`
	Results              [][]string                       `json:"results,omitempty"`
	Root                 types.HexBytes                   `json:"root,omitempty"`
	Siblings             types.HexBytes                   `json:"siblings,omitempty"`
	Size                 *int64                           `json:"size,omitempty"`
	State                string                           `json:"state,omitempty"`
	Stats                *VochainStats                    `json:"stats,omitempty"`
	Timestamp            int32                            `json:"timestamp"`
	Type                 string                           `json:"type,omitempty"`
	Tx                   *indexertypes.TxPackage          `json:"tx,omitempty"`
	TxList               []*indexertypes.TxMetadata       `json:"txList,omitempty"`
	URI                  string                           `json:"uri,omitempty"`
	ValidatorList        []*models.Validator              `json:"validatorlist,omitempty"`
	ValidProof           *bool                            `json:"validProof,omitempty"`
	Weight               *types.BigInt                    `json:"weight,omitempty"`
}

// VochainStats contains information about the current Vochain statistics and state
type VochainStats struct {
	BlockHeight      uint32    `json:"block_height"`
	EntityCount      int64     `json:"entity_count"`
	EnvelopeCount    uint64    `json:"envelope_count"`
	ProcessCount     int64     `json:"process_count"`
	TransactionCount uint64    `json:"transaction_count"`
	ValidatorCount   int       `json:"validator_count"`
	BlockTime        [5]int32  `json:"block_time"`
	BlockTimeStamp   int32     `json:"block_time_stamp"`
	ChainID          string    `json:"chain_id"`
	GenesisTimeStamp time.Time `json:"genesis_time_stamp"`
	Syncing          bool      `json:"syncing"`
}

type ProcessSummary struct {
	BlockCount      uint32               `json:"blockCount,omitempty"`
	EntityID        string               `json:"entityId,omitempty"`
	EntityIndex     uint32               `json:"entityIndex,omitempty"`
	EnvelopeHeight  *uint32              `json:"envelopeHeight,omitempty"`
	Metadata        string               `json:"metadata,omitempty"`
	SourceNetworkID string               `json:"sourceNetworkID,omitempty"`
	StartBlock      uint32               `json:"startBlock,omitempty"`
	State           string               `json:"state,omitempty"`
	EnvelopeType    *models.EnvelopeType `json:"envelopeType,omitempty"`
}

// Key associates a key string with an index, so clients can check
// the index of each process key.
type Key struct {
	Idx int    `json:"idx"`
	Key string `json:"key"`
}
