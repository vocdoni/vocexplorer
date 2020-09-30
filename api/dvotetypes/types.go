package dvotetypes

const (
	ProcessIDsize = 32

	// List of transation names
	TxVote              = "vote"
	TxNewProcess        = "newProcess"
	TxCancelProcess     = "cancelProcess"
	TxAddValidator      = "addValidator"
	TxRemoveValidator   = "removeValidator"
	TxAddOracle         = "addOracle"
	TxRemoveOracle      = "removeOracle"
	TxAddProcessKeys    = "addProcessKeys"
	TxRevealProcessKeys = "revealProcessKeys"
)

// ValidTypes represents an allowed specific tx type
var ValidTypes = map[string]string{
	TxVote:              "VoteTx",
	TxNewProcess:        "NewProcessTx",
	TxCancelProcess:     "CancelProcessTx",
	TxAddValidator:      "AdminTx",
	TxRemoveValidator:   "AdminTx",
	TxAddOracle:         "AdminTx",
	TxRemoveOracle:      "AdminTx",
	TxAddProcessKeys:    "AdminTx",
	TxRevealProcessKeys: "AdminTx",
}

// VotePackage represents the payload of a vote (usually base64 encoded)
type VotePackage struct {
	Nonce string `json:"nonce,omitempty"`
	Votes []int  `json:"votes"`
}

// Tx is an abstraction for any specific tx which is primarly defined by its type
// For now we have 3 tx types {voteTx, newProcessTx, adminTx}
type Tx struct {
	Type string `json:"type"`
}

// VoteTx represents the info required for submmiting a vote
type VoteTx struct {
	EncryptionKeyIndexes []int  `json:"encryptionKeyIndexes,omitempty"`
	Nonce                string `json:"nonce,omitempty"`
	Nullifier            string `json:"nullifier,omitempty"`
	ProcessID            string `json:"processId"`
	Proof                string `json:"proof,omitempty"`
	Signature            string `json:"signature,omitempty"`
	Type                 string `json:"type,omitempty"`
	VotePackage          string `json:"votePackage,omitempty"`
}

func (tx *VoteTx) TxType() string {
	return "VoteTx"
}

// NewProcessTx represents the info required for starting a new process
type NewProcessTx struct {
	// EntityID the process belongs to
	EntityID string `json:"entityId"`
	// MkRoot merkle root of all the census in the process
	MkRoot string `json:"mkRoot,omitempty"`
	// MkURI merkle tree URI
	MkURI string `json:"mkURI,omitempty"`
	// NumberOfBlocks represents the tendermint block where the process goes from active to finished
	NumberOfBlocks int64  `json:"numberOfBlocks"`
	ProcessID      string `json:"processId"`
	ProcessType    string `json:"processType"`
	Signature      string `json:"signature,omitempty"`
	// StartBlock represents the tendermint block where the process goes from scheduled to active
	StartBlock int64  `json:"startBlock"`
	Type       string `json:"type,omitempty"`
}

func (tx *NewProcessTx) TxType() string {
	return "NewProcessTx"
}

// CancelProcessTx represents a tx for canceling a valid process
type CancelProcessTx struct {
	// EntityID the process belongs to
	ProcessID string `json:"processId"`
	Signature string `json:"signature,omitempty"`
	Type      string `json:"type,omitempty"`
}

func (tx *CancelProcessTx) TxType() string {
	return "CancelProcessTx"
}

// AdminTx represents a Tx that can be only executed by some authorized addresses
type AdminTx struct {
	Address              string `json:"address"`
	CommitmentKey        string `json:"commitmentKey,omitempty"`
	EncryptionPrivateKey string `json:"encryptionPrivateKey,omitempty"`
	EncryptionPublicKey  string `json:"encryptionPublicKey,omitempty"`
	KeyIndex             int    `json:"keyIndex,omitempty"`
	Nonce                string `json:"nonce"`
	Power                int64  `json:"power,omitempty"`
	ProcessID            string `json:"processId,omitempty"`
	PubKey               string `json:"publicKey,omitempty"`
	RevealKey            string `json:"revealKey,omitempty"`
	Signature            string `json:"signature,omitempty"`
	Type                 string `json:"type"` // addValidator, removeValidator, addOracle, removeOracle
}

func (tx *AdminTx) TxType() string {
	return "AdminTx"
}

// ValidateType a valid Tx type specified in ValidTypes. Returns empty string if invalid type.
func ValidateType(t string) string {
	val, ok := ValidTypes[t]
	if !ok {
		return ""
	}
	return val
}
