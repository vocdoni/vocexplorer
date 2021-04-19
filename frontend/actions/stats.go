package actions

import (
	"gitlab.com/vocdoni/vocexplorer/api"
)

//SetStats is the action to set the blockchain statistics
type SetStats struct {
	Stats *api.VochainStats
}
