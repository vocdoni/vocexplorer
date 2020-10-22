package api

import "gitlab.com/vocdoni/vocexplorer/config"

//GetProcessEnvelopeCount returns the height of envelopes belonging to given process stored by the database
func GetProcessEnvelopeCount(process string) (int64, bool) {
	return getHeight("/api/procenvheight/?process=" + process)
}

//GetProcessEnvelopeCountMap returns the entire map of process envelope heights
func GetProcessEnvelopeCountMap() (map[string]int64, bool) {
	return getHeightMap("/api/heightmap/?key=" + config.ProcessEnvelopeCountMapKey)
}

//GetEnvelopeCount returns the latest envelope height stored by the database
func GetEnvelopeCount() (int64, bool) {
	return getHeight("/api/height/?key=" + config.LatestEnvelopeCountKey)
}

//GetProcessCount returns the latest process height stored by the database
func GetProcessCount() (int64, bool) {
	return getHeight("/api/height/?key=" + config.LatestProcessCountKey)
}

//GetEntityCount returns the latest envelope height stored by the database
func GetEntityCount() (int64, bool) {
	return getHeight("/api/height/?key=" + config.LatestEntityCountKey)
}

//GetEntityProcessCount returns the number of processes belonging to a
func GetEntityProcessCount(entity string) (int64, bool) {
	return getHeight("/api/entityprocheight/?entity=" + entity)
}

//GetEntityProcessCountMap returns the entire map of entity process heights
func GetEntityProcessCountMap() (map[string]int64, bool) {
	return getHeightMap("/api/heightmap/?key=" + config.EntityProcessCountMapKey)
}

//GetBlockHeight returns the latest block height stored by the database
func GetBlockHeight() (int64, bool) {
	return getHeight("/api/height/?key=" + config.LatestBlockHeightKey)
}

//GetTxHeight returns the latest tx height stored by the database
func GetTxHeight() (int64, bool) {
	return getHeight("/api/height/?key=" + config.LatestTxHeightKey)
}

//GetValidatorBlockHeight returns the height of blocks belonging to given validator stored by the database
func GetValidatorBlockHeight(proposer string) (int64, bool) {
	return getHeight("/api/numblocksvalidator/?proposer=" + proposer)
}

//GetValidatorBlockHeightMap returns the entire map of validator block heights
func GetValidatorBlockHeightMap() (map[string]int64, bool) {
	return getHeightMap("/api/heightmap/?key=" + config.ValidatorHeightMapKey)
}

//GetValidatorCount returns the latest validator count stored by the database
func GetValidatorCount() (int64, bool) {
	return getHeight("/api/height/?key=" + config.LatestValidatorCountKey)
}
