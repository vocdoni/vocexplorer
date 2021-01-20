package vochain

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/vocdoni/vocexplorer/api"
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
func (vs *VochainService) GetProcListResults(listSize int64) []string {
	return vs.scrut.ProcessListWithResults(listSize, []byte{})
}

// GetProcListLiveResults gets list of live processes on the Vochain
func (vs *VochainService) GetProcListLiveResults(listSize int64) []string {
	return vs.scrut.ProcessListWithLiveResults(listSize, []byte{})
}

// GetProcessList gets list of processes for a given entity
func (vs *VochainService) GetProcessList(entityID string, listSize int64) ([][]byte, error) {
	if listSize > MaxListIterations || listSize <= 0 {
		listSize = MaxListIterations
	}
	// check/sanitize eid
	entityID = util.TrimHex(entityID)
	eid, err := hex.DecodeString(entityID)
	if err != nil {
		return nil, fmt.Errorf("cannot decode entityID")
	}
	return vs.scrut.ProcessList(eid, []byte{}, listSize)
}

// GetProcessResults gets the results of a given process
func (vs *VochainService) GetProcessResults(processID string) (string, string, [][]uint64, error) {
	var err error
	processID = util.TrimHex(processID)
	// Get process info
	pid, err := hex.DecodeString(processID)
	if err != nil {
		return "", "", nil, fmt.Errorf("cannot decode processID")
	}
	procInfo, err := vs.scrut.ProcessInfo(pid)
	if err != nil {
		return "", "", nil, err
	}
	var procType string
	var state string

	if procInfo.EnvelopeType.EncryptedVotes {
		procType = "Encrypted"
	} else {
		procType = "Open"
	}
	if procInfo.EnvelopeType.Anonymous {
		procType = procType + " anonymous process"
	} else {
		procType = procType + " poll"
	}
	state = strings.Title(strings.ToLower(procInfo.Status.String()))

	// Get results info
	vr, err := vs.scrut.VoteResult(pid)
	if err != nil && err != scrutinizer.ErrNoResultsYet {
		return procType, state, nil, err
	}
	if err == scrutinizer.ErrNoResultsYet {
		return procType, state, nil, fmt.Errorf(scrutinizer.ErrNoResultsYet.Error())
	}
	results := vs.scrut.GetFriendlyResults(vr)
	return procType, state, results, nil
}
