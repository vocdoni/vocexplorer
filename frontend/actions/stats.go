package actions

import "go.vocdoni.io/proto/build/go/models"

//SetStats is the action to set the blockchain statistics
type SetStats struct {
	Stats *models.VochainStats
}
