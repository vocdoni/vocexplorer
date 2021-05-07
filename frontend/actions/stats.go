package actions

import "go.vocdoni.io/dvote/api"

//SetStats is the action to set the blockchain statistics
type SetStats struct {
	Stats *api.VochainStats
}
