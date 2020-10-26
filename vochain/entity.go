package vochain

import (
	"gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/go-dvote/util"
)

// GetScrutinizerEntities gets list of entities indexed by the scrutinizer on the Vochain
func (vs *VochainService) GetScrutinizerEntities(fromID string, listSize int64) []string {
	if listSize > MaxListIterations || listSize <= 0 {
		listSize = MaxListIterations
	}
	return vs.scrut.List(listSize, util.TrimHex(""), types.ScrutinizerEntityPrefix)
	// fromID does not work in this version of go-dvote. I think the problem is that entities are not ordered chronologically, so fromID seeks to the given ID and returns all ID's that are *lexicographically* after fromID
	// return vs.scrut.List(listSize, util.TrimHex(fromID), types.ScrutinizerEntityPrefix)
}
