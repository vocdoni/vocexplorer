package vochain

import (
	"encoding/hex"
	"errors"

	"gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/go-dvote/util"
)

// GetEntityCount gets number of entities
func (vs *VochainService) GetEntityCount() int64 {
	return vs.scrut.EntityCount()
}

// GetProcessCount gets number of processes
func (vs *VochainService) GetProcessCount() int64 {
	return vs.app.State.CountProcesses(true)
}

// GetEnvelopeCount gets number of envelopes in a given process
func (vs *VochainService) GetEnvelopeCount(processID string) (int64, error) {
	// check pid
	processID = util.TrimHex(processID)
	if !util.IsHexEncodedStringWithLength(processID, types.ProcessIDsize) {
		return 0, errors.New("cannot get envelope height: (malformed processId)")
	}
	pid, err := hex.DecodeString(processID)
	if err != nil {
		return 0, err
	}
	votes := vs.app.State.CountVotes(pid, true)
	return votes, nil
}

// GetTotalEnvelopeCount gets number of envelopes
func (vs *VochainService) GetTotalEnvelopeCount() (int64, error) {
	from := ""
	votes := int64(0)
	listSize := int64(100)
	for {
		// Get all live processes, sum envelopes
		newPIDs, err := vs.GetProcListLiveResults(from, listSize)
		if err != nil {
			return 0, err
		}
		for _, pid := range newPIDs {
			rawPid, err := hex.DecodeString(pid)
			if err != nil {
				return 0, err
			}
			votes += vs.app.State.CountVotes(rawPid, true)
		}
		if len(newPIDs) < int(listSize) {
			break
		}
		from = newPIDs[len(newPIDs)]
	}
	for {
		// Do the same for ended processes
		newPIDs, err := vs.GetProcListResults(from, listSize)
		if err != nil {
			return 0, err
		}
		for _, pid := range newPIDs {
			rawPid, err := hex.DecodeString(pid)
			if err != nil {
				return 0, err
			}
			votes += vs.app.State.CountVotes(rawPid, true)
		}
		if len(newPIDs) < int(listSize) {
			break
		}
		from = newPIDs[len(newPIDs)]
	}
	return votes, nil
}

// GetEntityProcessCount gets number of processes for a given entity
func (vs *VochainService) GetEntityProcessCount(eid string) (int64, error) {
	from := ""
	processes := int64(0)
	listSize := int64(100)
	for {
		// Get all live processes, sum envelopes
		newPIDs, err := vs.GetProcessList(eid, from, listSize)
		if err != nil {
			return 0, err
		}
		processes += int64(len(newPIDs))
		if len(newPIDs) < int(listSize) {
			break
		}
		from = newPIDs[len(newPIDs)]
	}
	return processes, nil
}

// GetBlockHeight gets the number of blocks
func (vs *VochainService) GetBlockHeight() int64 {
	return vs.GetStatus().LatestBlockHeight
}

// TODO add tx height
