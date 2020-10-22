package api

import (
	"encoding/json"

	"gitlab.com/vocdoni/vocexplorer/logger"
)

//GetStats gets the latest statistics
func GetStats() (*VochainStats, bool) {
	body, ok := requestBody("/api/stats")
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return &VochainStats{}, false
	}
	stats := new(VochainStats)
	err := json.NewDecoder(body).Decode(&stats)
	if err != nil {
		logger.Error(err)
		return stats, false
	}
	return stats, true
}
