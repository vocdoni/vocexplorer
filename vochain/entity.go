package vochain

// GetScrutinizerEntities gets list of entities indexed by the scrutinizer on the Vochain
func (vs *VochainService) GetScrutinizerEntities(fromID string, listSize int64) ([]string, error) {
	if listSize > MaxListIterations || listSize <= 0 {
		listSize = MaxListIterations
	}
	// vs.app.Node.BlockStore().LoadBlock(height).AppHash
	// vs.app.State.
	return vs.scrut.EntityList(listSize, fromID)
}
