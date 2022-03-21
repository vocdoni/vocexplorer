package actions

import "gitlab.com/vocdoni/vocexplorer/client"

//SetStats is the action to set the blockchain statistics
type SetStats struct {
	Stats *client.VochainStats
}
