package vochain

// GetScrutinizerEntities gets list of entities indexed by the scrutinizer on the Vochain
func (vs *VochainService) GetScrutinizerEntities(listSize int) []string {
	return vs.scrut.EntityList(listSize, 0)
}
