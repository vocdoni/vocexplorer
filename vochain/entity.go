package vochain

// GetScrutinizerEntities gets list of entities indexed by the scrutinizer on the Vochain
func (vs *VochainService) GetScrutinizerEntities(listSize int64) ([]string, error) {
	return vs.scrut.EntityList(listSize, "")
}
