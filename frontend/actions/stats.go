package actions

import "go.vocdoni.io/dvote/types"

//SetStats is the action to set the blockchain statistics
type SetStats struct {
	Stats *types.VochainStats
}
