package components

import (
	"strconv"

	"github.com/golang/protobuf/ptypes"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/dbapi"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
	router "marwan.io/vecty-router"
)

// TxsView renders the Txs page
type TxsView struct {
	vecty.Core
	cfg *config.Cfg
}

// Render renders the TxsView component
func (home *TxsView) Render() vecty.ComponentOrHTML {
	height, err := strconv.ParseInt(router.GetNamedVar(home)["id"], 0, 64)
	util.ErrPrint(err)
	tx := dbapi.GetTx(height)
	// Get block which houses tx
	block := dbapi.GetBlock(tx.Store.Height)
	if tx == nil {
		log.Errorf("Tx unavailable")
		return elem.Div(
			&Header{},
			elem.Main(vecty.Text("Tx not available")),
		)
	}
	tm, err := ptypes.Timestamp(block.GetTime())
	util.ErrPrint(err)
	return elem.Div(
		&Header{},
		&TxContents{
			Tx:       tx,
			Time:     tm,
			HasBlock: !types.BlockIsEmpty(block),
		},
	)
}
