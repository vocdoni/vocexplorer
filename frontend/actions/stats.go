package actions

import sctypes "go.vocdoni.io/dvote/vochain/scrutinizer/types"

//SetStats is the action to set the blockchain statistics
type SetStats struct {
	Stats *sctypes.VochainStats
}
