package vochain

import (
	"encoding/hex"
	"errors"

	"gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/go-dvote/util"
	"gitlab.com/vocdoni/go-dvote/vochain/scrutinizer"
	"gitlab.com/vocdoni/vocexplorer/api"
)

// GetProcessKeys gets process keys
func (vs *VochainService) GetProcessKeys(processID string) (*api.Pkeys, error) {
	pid, err := hex.DecodeString(processID)
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
func (vs *VochainService) GetProcListResults(fromID string, listSize int64) ([]string, error) {
	if listSize > MaxListIterations || listSize <= 0 {
		listSize = MaxListIterations
	}
	return vs.scrut.ProcessListWithResults(listSize, fromID)
}

// GetProcListLiveResults gets list of live processes on the Vochain
func (vs *VochainService) GetProcListLiveResults(fromID string, listSize int64) ([]string, error) {
	if listSize > MaxListIterations || listSize <= 0 {
		listSize = MaxListIterations
	}
	return vs.scrut.ProcessListWithLiveResults(listSize, fromID)
}

// GetProcessList gets list of processes for a given entity, starting at from
func (vs *VochainService) GetProcessList(entityID, fromID string, listSize int64) ([]string, error) {
	if listSize > MaxListIterations || listSize <= 0 {
		listSize = MaxListIterations
	}
	// check/sanitize eid and fromId
	entityID = util.TrimHex(entityID)
	if !util.IsHexEncodedStringWithLength(entityID, types.EntityIDsize) &&
		!util.IsHexEncodedStringWithLength(entityID, types.EntityIDsizeV2) {
		return nil, errors.New("cannot get process list: (malformed entityId)")
	}
	if len(fromID) > 0 {
		fromID = util.TrimHex(fromID)
		if !util.IsHexEncodedStringWithLength(fromID, types.ProcessIDsize) {
			return nil, errors.New("cannot get process list: (malformed entityId)")
		}
	}
	eid, err := hex.DecodeString(entityID)
	if err != nil {
		return nil, errors.New("cannot decode entityID")
	}
	fromIDBytes := []byte{}
	if fromID != "" {
		fromIDBytes, err = hex.DecodeString(fromID)
		if err != nil {
			return nil, errors.New("cannot decode fromID")
		}
	}
	return vs.scrut.ProcessList(eid, fromIDBytes, listSize)
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
	var state string
	if procInfo.Canceled {
		state = "canceled"
	} else {
		state = "active"
	}

	// Get results info
	vr, err := vs.scrut.VoteResult(pid)
	if err != nil && err != scrutinizer.ErrNoResultsYet {
		return "", "", nil, err
	}
	if err == scrutinizer.ErrNoResultsYet {
		return "", "", nil, errors.New(scrutinizer.ErrNoResultsYet.Error())
	}
	return procInfo.Type, state, vr, nil
}
