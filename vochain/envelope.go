package vochain

import (
	"encoding/hex"
	"errors"

	"github.com/vocdoni/dvote-protobuf/build/go/models"
	"go.vocdoni.io/dvote/types"
	"go.vocdoni.io/dvote/util"
)

// GetEnvelope gets contents of given envelope
func (vs *VochainService) GetEnvelope(processID, nullifier string) (*models.Vote, error) {
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
