package api

import (
	"encoding/json"
)

// Key is a cryptographic key
type Key struct {
	Idx int    `json:"idx"`
	Key string `json:"key"`
}

// Pkeys is the set of cryptographic keys for a process
type Pkeys struct {
	Pub  []Key
	Priv []Key
	Comm []Key
	Rev  []Key
}

// RequestMessage holds a decoded request but does not decode the body. from go-dvote
type RequestMessage struct {
	ID          string          `json:"id"`
	MetaRequest json.RawMessage `json:"request"`
	Signature   string          `json:"signature"`
}

// MetaRequest holds a gateway api request, from go-dvote/types
type MetaRequest struct {
	CensusID   string   `json:"censusId,omitempty"`
	CensusURI  string   `json:"censusUri,omitempty"`
	ClaimData  string   `json:"claimData,omitempty"`
	ClaimsData []string `json:"claimsData,omitempty"`
	Content    string   `json:"content,omitempty"`
	Digested   bool     `json:"digested,omitempty"`
	EntityID   string   `json:"entityId,omitempty"`
	From       int64    `json:"from,omitempty"`
	FromID     string   `json:"fromId,omitempty"`
	ListSize   int64    `json:"listSize,omitempty"`
	Method     string   `json:"method"`
	Name       string   `json:"name,omitempty"`
	Nullifier  string   `json:"nullifier,omitempty"`
	ProcessID  string   `json:"processId,omitempty"`
	ProofData  string   `json:"proofData,omitempty"`
	PubKeys    []string `json:"pubKeys,omitempty"`
	RawTx      string   `json:"rawTx,omitempty"`
	RootHash   string   `json:"rootHash,omitempty"`
	Signature  string   `json:"signature,omitempty"`
	Timestamp  int32    `json:"timestamp"`
	Type       string   `json:"type,omitempty"`
	URI        string   `json:"uri,omitempty"`
}

// ResponseMessage wraps an api response, from go-dvote/types
type ResponseMessage struct {
	MetaResponse json.RawMessage `json:"response"`

	ID        string `json:"id"`
	Signature string `json:"signature"`
}

// MetaResponse holds a gateway api response, from go-dvote/types
type MetaResponse struct {
	APIList              []string   `json:"apiList,omitempty"`
	BlockTime            *[5]int32  `json:"blockTime,omitempty"`
	BlockTimestamp       int32      `json:"blockTimestamp,omitempty"`
	CensusID             string     `json:"censusId,omitempty"`
	CensusList           []string   `json:"censusList,omitempty"`
	ClaimsData           []string   `json:"claimsData,omitempty"`
	Content              string     `json:"content,omitempty"`
	EntityID             string     `json:"entityId,omitempty"`
	EntityIDs            []string   `json:"entityIds,omitempty"`
	Files                []byte     `json:"files,omitempty"`
	Finished             *bool      `json:"finished,omitempty"`
	Health               int32      `json:"health,omitempty"`
	Height               *int64     `json:"height,omitempty"`
	InvalidClaims        []int      `json:"invalidClaims,omitempty"`
	Message              string     `json:"message,omitempty"`
	Nullifier            string     `json:"nullifier,omitempty"`
	Nullifiers           *[]string  `json:"nullifiers,omitempty"`
	Ok                   bool       `json:"ok"`
	Paused               *bool      `json:"paused,omitempty"`
	Payload              string     `json:"payload,omitempty"`
	ProcessIDs           []string   `json:"processIds,omitempty"`
	ProcessList          []string   `json:"processList,omitempty"`
	Registered           *bool      `json:"registered,omitempty"`
	Request              string     `json:"request"`
	Results              [][]uint32 `json:"results,omitempty"`
	Root                 string     `json:"root,omitempty"`
	Siblings             string     `json:"siblings,omitempty"`
	Size                 *int64     `json:"size,omitempty"`
	State                string     `json:"state,omitempty"`
	Timestamp            int32      `json:"timestamp"`
	Type                 string     `json:"type,omitempty"`
	URI                  string     `json:"uri,omitempty"`
	ValidProof           *bool      `json:"validProof,omitempty"`
	CommitmentKeys       []Key      `json:"commitmentKeys,omitempty"`
	EncryptionPrivKeys   []Key      `json:"encryptionPrivKeys,omitempty"`
	RevealKeys           []Key      `json:"revealKeys,omitempty"`
	EncryptionPublicKeys []Key      `json:"encryptionPubKeys,omitempty"`
}