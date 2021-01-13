package vochain

import (
	"encoding/hex"

	"go.vocdoni.io/dvote/util"
	"go.vocdoni.io/proto/build/go/models"
)

// GetEnvelope gets contents of given envelope
func (vs *VochainService) GetEnvelope(processID, nullifier string) (*models.Vote, error) {
	nullifier = util.TrimHex(nullifier)
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
