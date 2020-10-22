package vochain

import (
	"encoding/hex"
	"errors"
	"fmt"

	"gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/go-dvote/util"
)

// GetEnvelopeList gets list of envelopes in a given process, starting at from
func (vs *VochainService) GetEnvelopeList(processID string, from int64, listSize int64) ([]string, error) {
	if listSize > MaxListIterations || listSize <= 0 {
		listSize = MaxListIterations
	}
	processID = util.TrimHex(processID)
	if !util.IsHexEncodedStringWithLength(processID, types.ProcessIDsize) {
		return nil, errors.New("cannot get envelope list: (malformed processId)")
	}
	pid, err := hex.DecodeString(processID)
	if err != nil {
		return nil, errors.New("cannot decode processID")
	}
	nullifiers := vs.app.State.EnvelopeList(pid, from, listSize, true)
	strnull := []string{}
	for _, n := range nullifiers {
		strnull = append(strnull, fmt.Sprintf("%x", n))
	}
	return strnull, nil
}

// GetEnvelope gets contents of given envelope
func (vs *VochainService) GetEnvelope(processID, nullifier string) (*types.Vote, error) {
	// check pid
	if !util.IsHexEncodedStringWithLength(processID, types.ProcessIDsize) {
		return nil, errors.New("cannot get envelope: (malformed processId)")
	}
	// check nullifier
	nullifier = util.TrimHex(nullifier)
	if !util.IsHexEncodedStringWithLength(nullifier, types.VoteNullifierSize) {
		return nil, errors.New("cannot get envelope: (malformed nullifier)")
	}
	pid, err := hex.DecodeString(util.TrimHex(processID))
	if err != nil {
		return nil, err
	}
	nullifierBytes, err := hex.DecodeString(util.TrimHex(nullifier))
	if err != nil {
		return nil, err
	}
	envelope, err := vs.app.State.Envelope(pid, nullifierBytes, true)
	if err != nil {
		return nil, err
	}
	return envelope, nil
}
