// package tmtypes

// import (
// 	"errors"
// 	"math/bits"
// 	"time"

// 	"github.com/tendermint/tendermint/crypto"
// 	"github.com/tendermint/tendermint/crypto/tmhash"
// 	tmbytes "github.com/tendermint/tendermint/libs/bytes"
// 	"github.com/tendermint/tendermint/proto/tendermint/version"
// )

// var (
// 	leafPrefix  = []byte{0}
// 	innerPrefix = []byte{1}
// )

// // Address is hex bytes.
// type Address = crypto.Address

// // ConsensusParams contains consensus critical parameters that determine the
// // validity of blocks.
// type ConsensusParams struct {
// 	Block     BlockParams     `json:"block"`
// 	Evidence  EvidenceParams  `json:"evidence"`
// 	Validator ValidatorParams `json:"validator"`
// }

// // BlockParams define limits on the block size and gas plus minimum time
// // between blocks.
// type BlockParams struct {
// 	MaxBytes int64 `json:"max_bytes"`
// 	MaxGas   int64 `json:"max_gas"`
// 	// Minimum time increment between consecutive blocks (in milliseconds)
// 	// Not exposed to the application.
// 	TimeIotaMs int64 `json:"time_iota_ms"`
// }

// // EvidenceParams determine how we handle evidence of malfeasance.
// type EvidenceParams struct {
// 	MaxAgeNumBlocks int64         `json:"max_age_num_blocks"` // only accept new evidence more recent than this
// 	MaxAgeDuration  time.Duration `json:"max_age_duration"`
// }

// // ValidatorParams restrict the public key types validators can use.
// // NOTE: uses ABCI pubkey naming, not Amino names.
// type ValidatorParams struct {
// 	PubKeyTypes []string `json:"pub_key_types"`
// }

// // Txs is a slice of Tx.
// type Txs []Tx

// // Tx is an arbitrary byte array.
// // NOTE: Tx has no types at this level, so when wire encoded it's just length-prefixed.
// // Might we want types here ?
// type Tx []byte

// // Hash computes the TMHASH hash of the wire encoded transaction.
// func (tx Tx) Hash() []byte {
// 	return tmhash.Sum(tx)
// }

// // ResultBlock is a Single block (with meta)
// type ResultBlock struct {
// 	BlockID BlockID `json:"block_id"`
// 	Block   *Block  `json:"block"`
// }

// // Block defines the atomic unit of a Tendermint blockchain.
// type Block struct {
// 	Header     `json:"header"`
// 	Data       `json:"data"`
// 	Evidence   EvidenceData `json:"evidence"`
// 	LastCommit *Commit      `json:"last_commit"`
// }

// // Size returns size of the block in bytes.
// func (b *Block) Size() int {
// 	pbb, err := b.ToProto()
// 	if err != nil {
// 		return 0
// 	}

// 	return pbb.Size()
// }

// // ToProto converts Block to protobuf
// func (b *Block) ToProto() (*tmproto.Block, error) {
// 	if b == nil {
// 		return nil, errors.New("nil Block")
// 	}

// 	pb := new(tmproto.Block)

// 	pb.Header = *b.Header.ToProto()
// 	pb.LastCommit = b.LastCommit.ToProto()
// 	pb.Data = b.Data.ToProto()

// 	protoEvidence, err := b.Evidence.ToProto()
// 	if err != nil {
// 		return nil, err
// 	}
// 	pb.Evidence = *protoEvidence

// 	return pb, nil
// }

// // Header defines the structure of a Tendermint block header.
// // NOTE: changes to the Header should be duplicated in:
// // - header.Hash()
// // - abci.Header
// // - https://github.com/tendermint/spec/blob/master/spec/blockchain/blockchain.md
// type Header struct {
// 	// basic block info
// 	Version version.Consensus `json:"version"`
// 	ChainID string            `json:"chain_id"`
// 	Height  int64             `json:"height"`
// 	Time    time.Time         `json:"time"`

// 	// prev block info
// 	LastBlockID BlockID `json:"last_block_id"`

// 	// hashes of block data
// 	LastCommitHash tmbytes.HexBytes `json:"last_commit_hash"` // commit from validators from the last block
// 	DataHash       tmbytes.HexBytes `json:"data_hash"`        // transactions

// 	// hashes from the app output from the prev block
// 	ValidatorsHash     tmbytes.HexBytes `json:"validators_hash"`      // validators for the current block
// 	NextValidatorsHash tmbytes.HexBytes `json:"next_validators_hash"` // validators for the next block
// 	ConsensusHash      tmbytes.HexBytes `json:"consensus_hash"`       // consensus params for current block
// 	AppHash            tmbytes.HexBytes `json:"app_hash"`             // state after txs from the previous block
// 	// root hash of all results from the txs from the previous block
// 	LastResultsHash tmbytes.HexBytes `json:"last_results_hash"`

// 	// consensus info
// 	EvidenceHash    tmbytes.HexBytes `json:"evidence_hash"`    // evidence included in the block
// 	ProposerAddress Address          `json:"proposer_address"` // original proposer of the block
// }

// // Hash returns the hash of the header.
// // It computes a Merkle tree from the header fields
// // ordered as they appear in the Header.
// // Returns nil if ValidatorHash is missing,
// // since a Header is not valid unless there is
// // a ValidatorsHash (corresponding to the validator set).
// func (h *Header) Hash() tmbytes.HexBytes {
// 	if h == nil || len(h.ValidatorsHash) == 0 {
// 		return nil
// 	}
// 	return SimpleHashFromByteSlices([][]byte{
// 		cdcEncode(h.Version),
// 		cdcEncode(h.ChainID),
// 		cdcEncode(h.Height),
// 		cdcEncode(h.Time),
// 		cdcEncode(h.LastBlockID),
// 		cdcEncode(h.LastCommitHash),
// 		cdcEncode(h.DataHash),
// 		cdcEncode(h.ValidatorsHash),
// 		cdcEncode(h.NextValidatorsHash),
// 		cdcEncode(h.ConsensusHash),
// 		cdcEncode(h.AppHash),
// 		cdcEncode(h.LastResultsHash),
// 		cdcEncode(h.EvidenceHash),
// 		cdcEncode(h.ProposerAddress),
// 	})
// }

// // Data contains the set of transactions included in the block
// type Data struct {

// 	// Txs that will be applied by state @ block.Height+1.
// 	// NOTE: not all txs here are valid.  We're just agreeing on the order first.
// 	// This means that block.AppHash does not include these txs.
// 	Txs Txs `json:"txs"`

// 	// Volatile
// 	// hash tmbytes.HexBytes
// }

// // EvidenceData contains any evidence of malicious wrong-doing by validators
// type EvidenceData struct {
// 	Evidence EvidenceList `json:"evidence"`

// 	// Volatile
// 	// hash tmbytes.HexBytes
// }

// // EvidenceList is a list of Evidence. Evidences is not a word.
// type EvidenceList []Evidence

// // Evidence represents any provable malicious activity by a validator
// type Evidence interface {
// 	Height() int64                                     // height of the equivocation
// 	Time() time.Time                                   // time of the equivocation
// 	Address() []byte                                   // address of the equivocating validator
// 	Bytes() []byte                                     // bytes which comprise the evidence
// 	Hash() []byte                                      // hash of the evidence
// 	Verify(chainID string, pubKey crypto.PubKey) error // verify the evidence
// 	Equal(Evidence) bool                               // check equality of evidence

// 	ValidateBasic() error
// 	String() string
// }

// // Commit contains the evidence that a block was committed by a set of validators.
// // NOTE: Commit is empty for height 1, but never nil.
// type Commit struct {
// 	// NOTE: The signatures are in order of address to preserve the bonded
// 	// ValidatorSet order.
// 	// Any peer with a block can gossip signatures by index with a peer without
// 	// recalculating the active ValidatorSet.
// 	Height     int64       `json:"height"`
// 	Round      int         `json:"round"`
// 	BlockID    BlockID     `json:"block_id"`
// 	Signatures []CommitSig `json:"signatures"`
// }

// // BitArray is a thread-safe implementation of a bit array.
// type BitArray struct {
// 	Bits  int      `json:"bits"`  // NOTE: persisted via reflect, must be exported
// 	Elems []uint64 `json:"elems"` // NOTE: persisted via reflect, must be exported
// }

// // CommitSig is a part of the Vote included in a Commit.
// type CommitSig struct {
// 	BlockIDFlag      BlockIDFlag `json:"block_id_flag"`
// 	ValidatorAddress Address     `json:"validator_address"`
// 	Timestamp        time.Time   `json:"timestamp"`
// 	Signature        []byte      `json:"signature"`
// }

// // BlockIDFlag indicates which BlockID the signature is for.
// type BlockIDFlag byte

// // BlockID defines the unique ID of a block as its Hash and its PartSetHeader
// type BlockID struct {
// 	Hash        tmbytes.HexBytes `json:"hash"`
// 	PartsHeader PartSetHeader    `json:"parts"`
// }

// type PartSetHeader struct {
// 	Total int              `json:"total"`
// 	Hash  tmbytes.HexBytes `json:"hash"`
// }

// // SimpleHashFromByteSlices computes a Merkle tree where the leaves are the byte slice,
// // in the provided order.
// func SimpleHashFromByteSlices(items [][]byte) []byte {
// 	switch len(items) {
// 	case 0:
// 		return nil
// 	case 1:
// 		return leafHash(items[0])
// 	default:
// 		k := getSplitPoint(len(items))
// 		left := SimpleHashFromByteSlices(items[:k])
// 		right := SimpleHashFromByteSlices(items[k:])
// 		return innerHash(left, right)
// 	}
// }

// // returns tmhash(0x00 || leaf)
// func leafHash(leaf []byte) []byte {
// 	return tmhash.Sum(append(leafPrefix, leaf...))
// }

// // returns tmhash(0x01 || left || right)
// func innerHash(left []byte, right []byte) []byte {
// 	return tmhash.Sum(append(innerPrefix, append(left, right...)...))
// }

// // getSplitPoint returns the largest power of 2 less than length
// func getSplitPoint(length int) int {
// 	if length < 1 {
// 		panic("Trying to split a tree with size < 1")
// 	}
// 	uLength := uint(length)
// 	bitlen := bits.Len(uLength)
// 	k := 1 << uint(bitlen-1)
// 	if k == length {
// 		k >>= 1
// 	}
// 	return k
// }
