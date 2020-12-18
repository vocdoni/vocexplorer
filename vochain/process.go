package vochain

import (
	"encoding/hex"
	"errors"

	"gitlab.com/vocdoni/vocexplorer/api"
	"go.vocdoni.io/dvote/types"
	"go.vocdoni.io/dvote/util"
	"go.vocdoni.io/dvote/vochain/scrutinizer"
)

// GetProcessKeys gets process keys
func (vs *VochainService) GetProcessKeys(processID string) (*api.Pkeys, error) {
	pid, err := hex.DecodeString(util.TrimHex(processID))
	if err != nil {
		return nil, err
	}
	p, err := vs.app.State.Process(pid, true)
	if err != nil {
		return nil, err
	}
	pkeys := api.Pkeys{}
	for i, p := range p.EncryptionPublicKeys {
		if len(p) > 0 {
			k := types.Key{}
			k.Key = p
			k.Idx = i
			pkeys.Pub = append(pkeys.Pub, k)
		}
	}
	for i, p := range p.EncryptionPrivateKeys {
		if len(p) > 0 {
			k := types.Key{}
			k.Key = p
			k.Idx = i
			pkeys.Priv = append(pkeys.Priv, k)
		}
	}
	for i, p := range p.CommitmentKeys {
		if len(p) > 0 {
			k := types.Key{}
			k.Key = p
			k.Idx = i
			pkeys.Comm = append(pkeys.Comm, k)
		}
	}
	for i, p := range p.RevealKeys {
		if len(p) > 0 {
			k := types.Key{}
			k.Key = p
			k.Idx = i
			pkeys.Rev = append(pkeys.Rev, k)
		}
	}
	return &pkeys, nil
}

// GetProcListResults gets list of finished processes on the Vochain
func (vs *VochainService) GetProcListResults(listSize int64) ([]string, error) {
	return vs.scrut.ProcessListWithResults(listSize, "")
}

// GetProcListLiveResults gets list of live processes on the Vochain
func (vs *VochainService) GetProcListLiveResults(listSize int64) ([]string, error) {
	return vs.scrut.ProcessListWithLiveResults(listSize, "")
}

// GetProcessList gets list of processes for a given entity
func (vs *VochainService) GetProcessList(entityID string, listSize int64) ([][]byte, error) {
	if listSize > MaxListIterations || listSize <= 0 {
		listSize = MaxListIterations
	}
	// check/sanitize eid
	entityID = util.TrimHex(entityID)
	if !util.IsHexEncodedStringWithLength(entityID, types.EntityIDsize) &&
		!util.IsHexEncodedStringWithLength(entityID, types.EntityIDsizeV2) {
		return nil, errors.New("cannot get process list: (malformed entityId)")
	}
	eid, err := hex.DecodeString(entityID)
	if err != nil {
		return nil, errors.New("cannot decode entityID")
	}
	return vs.scrut.ProcessList(eid, []byte{}, listSize)
}

// GetProcessResults gets the results of a given process
func (vs *VochainService) GetProcessResults(processID string) (string, string, [][]uint32, error) {
	var err error
	processID = util.TrimHex(processID)
	if !util.IsHexEncodedStringWithLength(processID, types.ProcessIDsize) {
		return "", "", nil, errors.New("cannot get results: (malformed processId)")
	}

	// Get process info
	pid, err := hex.DecodeString(processID)
	if err != nil {
		return "", "", nil, errors.New("cannot decode processID")
	}
	procInfo, err := vs.scrut.ProcessInfo(pid)
	if err != nil {
		return "", "", nil, err
	}
	var procType string
	var state string

	if procInfo.EnvelopeType.Anonymous {
		procType = "anonymous"
	} else {
		procType = "poll"
	}
	if procInfo.EnvelopeType.EncryptedVotes {
		procType = procType + " encrypted"
	} else {
		procType = procType + " open"
	}
	if procInfo.EnvelopeType.Serial {
		procType = procType + " serial"
	} else {
		procType = procType + " single"
	}
	state = procInfo.Status.String()

	// Get results info
	vr, err := vs.scrut.VoteResult(pid)
	if err != nil && err != scrutinizer.ErrNoResultsYet {
		return procType, state, nil, err
	}
	if err == scrutinizer.ErrNoResultsYet {
		return procType, state, nil, errors.New(scrutinizer.ErrNoResultsYet.Error())
	}
	results := vs.scrut.GetFriendlyResults(vr)
	return procType, state, results, nil
}
