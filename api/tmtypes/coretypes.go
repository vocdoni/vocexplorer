package tmtypes

import (
	"time"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/bytes"
)

// Genesis file
type ResultGenesis struct {
	Genesis *GenesisDoc `json:"genesis"`
}

// Single block (with meta)
type ResultBlock struct {
	BlockID BlockID `json:"block_id"`
	Block   *Block  `json:"block"`
}

// Info about the node's syncing state
type SyncInfo struct {
	LatestBlockHash   bytes.HexBytes `json:"latest_block_hash"`
	LatestAppHash     bytes.HexBytes `json:"latest_app_hash"`
	LatestBlockHeight int64          `json:"latest_block_height"`
	LatestBlockTime   time.Time      `json:"latest_block_time"`

	EarliestBlockHash   bytes.HexBytes `json:"earliest_block_hash"`
	EarliestAppHash     bytes.HexBytes `json:"earliest_app_hash"`
	EarliestBlockHeight int64          `json:"earliest_block_height"`
	EarliestBlockTime   time.Time      `json:"earliest_block_time"`

	CatchingUp bool `json:"catching_up"`
}

// Node Status
type ResultStatus struct {
	NodeInfo      DefaultNodeInfo `json:"node_info"`
	SyncInfo      SyncInfo        `json:"sync_info"`
	ValidatorInfo ValidatorInfo   `json:"validator_info"`
}

// Info about the node's validator
type ValidatorInfo struct {
	Address     bytes.HexBytes `json:"address"`
	PubKey      crypto.PubKey  `json:"pub_key"`
	VotingPower int64          `json:"voting_power"`
}

// Validators for a height.
type ResultValidators struct {
	BlockHeight int64        `json:"block_height"`
	Validators  []*Validator `json:"validators"`
	// Count of actual validators in this result
	Count int `json:"count"`
	// Total number of validators
	Total int `json:"total"`
}

// Result of querying for a tx
type ResultTx struct {
	Hash     bytes.HexBytes    `json:"hash"`
	Height   int64             `json:"height"`
	Index    uint32            `json:"index"`
	TxResult ResponseDeliverTx `json:"tx_result"`
	Tx       Tx                `json:"tx"`
	Proof    TxProof           `json:"proof,omitempty"`
}

// TxProof represents a Merkle proof of the presence of a transaction in the Merkle tree.
type TxProof struct {
	RootHash bytes.HexBytes `json:"root_hash"`
	Data     Tx             `json:"data"`
	Proof    SimpleProof    `json:"proof"`
}

type SimpleProof struct {
	Total    int      `json:"total"`     // Total number of items.
	Index    int      `json:"index"`     // Index of item to prove.
	LeafHash []byte   `json:"leaf_hash"` // Hash of item value.
	Aunts    [][]byte `json:"aunts"`     // Hashes from leaf's sibling to a root's child.
}

// Volatile state for each Validator
// NOTE: The ProposerPriority is not included in Validator.Hash();
// make sure to update that method if changes are made here
type Validator struct {
	Address     Address       `json:"address"`
	PubKey      crypto.PubKey `json:"pub_key"`
	VotingPower int64         `json:"voting_power"`

	ProposerPriority int64 `json:"proposer_priority"`
}

type ResponseDeliverTx struct {
	Code                 uint32   `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Data                 []byte   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Log                  string   `protobuf:"bytes,3,opt,name=log,proto3" json:"log,omitempty"`
	Info                 string   `protobuf:"bytes,4,opt,name=info,proto3" json:"info,omitempty"`
	GasWanted            int64    `protobuf:"varint,5,opt,name=gas_wanted,json=gasWanted,proto3" json:"gas_wanted,omitempty"`
	GasUsed              int64    `protobuf:"varint,6,opt,name=gas_used,json=gasUsed,proto3" json:"gas_used,omitempty"`
	Events               []Event  `protobuf:"bytes,7,rep,name=events,proto3" json:"events,omitempty"`
	Codespace            string   `protobuf:"bytes,8,opt,name=codespace,proto3" json:"codespace,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

type Event struct {
	Type                 string   `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Attributes           []Pair   `protobuf:"bytes,2,rep,name=attributes,proto3" json:"attributes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

type Pair struct {
	Key                  []byte   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}
