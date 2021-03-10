package vochain

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/vocdoni/vocexplorer/api"
	"github.com/vocdoni/vocexplorer/config"
	"go.vocdoni.io/dvote/types"
	"go.vocdoni.io/dvote/util"
	"go.vocdoni.io/dvote/vochain/scrutinizer"
	"go.vocdoni.io/proto/build/go/models"
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
func (vs *VochainService) GetProcListResults(listSize int) []string {
	return vs.scrut.ProcessListWithResults(listSize, 0)
}

// GetProcListLiveResults gets list of live processes on the Vochain
func (vs *VochainService) GetProcListLiveResults(listSize int) []string {
	return vs.scrut.ProcessListWithLiveResults(listSize, 0)
}

// GetProcessList gets list of processes for a given entity
func (vs *VochainService) GetProcessList(entityID string, listSize int) ([][]byte, error) {
	// check/sanitize eid
	entityID = util.TrimHex(entityID)
	eid, err := hex.DecodeString(entityID)
	if err != nil {
		return nil, fmt.Errorf("cannot decode entityID")
	}
	return vs.scrut.ProcessList(eid, config.DefaultNamespace, 0, listSize)
}

// GetProcessResults gets the results of a given process
func (vs *VochainService) GetProcessResults(processID string) (api.ProcessResults, error) {
	var err error
	var process api.ProcessResults = api.ProcessResults{}
	processID = util.TrimHex(processID)
	// Get process info
	pid, err := hex.DecodeString(processID)
	if err != nil {
		return process, fmt.Errorf("cannot decode processID")
	}
	procInfo, err := vs.scrut.ProcessInfo(pid)
	if err != nil {
		return process, err
	}
	if procInfo == nil {
		return process, nil
	}

	// Set basic readable process type
	if procInfo.Envelope.EncryptedVotes {
		process.Type = "Encrypted"
	} else {
		process.Type = "Open"
	}
	if procInfo.Envelope.Anonymous {
		process.Type = process.Type + " anonymous process"
	} else {
		process.Type = process.Type + " poll"
	}
	// Set readable process state
	process.State = strings.Title(strings.ToLower(models.ProcessStatus_name[procInfo.Status]))

	// Set full-info EnvelopeType
	if procInfo.Envelope != nil {
		process.EnvelopeType = api.EnvelopeType{
			Serial:         procInfo.Envelope.Serial,
			Anonymous:      procInfo.Envelope.Anonymous,
			EncryptedVotes: procInfo.Envelope.EncryptedVotes,
			UniqueValues:   procInfo.Envelope.UniqueValues}
	}

	// Set full-info ProcessMode
	if procInfo.Mode != nil {
		process.Mode = api.ProcessMode{
			AutoStart:         procInfo.Mode.AutoStart,
			Interruptible:     procInfo.Mode.Interruptible,
			DynamicCensus:     procInfo.Mode.DynamicCensus,
			EncryptedMetaData: procInfo.Mode.EncryptedMetaData}
	}

	// Set VoteOptions
	if procInfo.VoteOpts != nil {
		process.VoteOptions = api.ProcessVoteOptions{
			MaxCount:          procInfo.VoteOpts.MaxCount,
			MaxValue:          procInfo.VoteOpts.MaxValue,
			MaxVoteOverwrites: procInfo.VoteOpts.MaxVoteOverwrites,
			MaxTotalCost:      procInfo.VoteOpts.MaxTotalCost,
			CostExponent:      procInfo.VoteOpts.CostExponent}
	}

	// Set census info
	process.CensusOrigin = models.CensusOrigin_name[procInfo.CensusOrigin]
	process.CensusRoot = procInfo.CensusRoot
	process.CensusURI = procInfo.CensusURI

	// Set start + end block
	process.StartBlock = procInfo.StartBlock
	process.EndBlock = procInfo.EndBlock

	// Get results info
	vr, err := vs.scrut.GetResults(pid)
	if err != nil && err != scrutinizer.ErrNoResultsYet {
		return process, err
	}
	if err == scrutinizer.ErrNoResultsYet {
		return process, fmt.Errorf(scrutinizer.ErrNoResultsYet.Error())
	}
	process.Results = vs.scrut.GetFriendlyResults(vr.Votes)
	return process, nil
}
