package pages

import (
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/dbapi"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
	router "marwan.io/vecty-router"
)

// TxsView renders the Txs page
type TxsView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the TxsView component
func (home *TxsView) Render() vecty.ComponentOrHTML {
	height, err := strconv.ParseInt(router.GetNamedVar(home)["id"], 0, 64)
	util.ErrPrint(err)
	tx, ok := dbapi.GetTx(height)
	// Get block which houses tx
	if tx == nil || !ok {
		log.Errorf("Tx unavailable")
		return elem.Div(
			elem.Main(vecty.Text("Tx not available")),
		)
	}
	block, ok := dbapi.GetBlock(tx.Store.Height)
	var tm time.Time
	if ok {
		tm, err = ptypes.Timestamp(block.GetTime())
	}
	util.ErrPrint(err)
	return components.Container(
		elem.Section(
			&components.TxContents{
				Tx:       tx,
				Time:     tm,
				HasBlock: !types.BlockIsEmpty(block),
			},
		),
	)
}
