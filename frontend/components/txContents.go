package components

import (
	"encoding/json"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
	"github.com/xeonx/timeago"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// TxContents renders tx contents
type TxContents struct {
	vecty.Core
	Tx   *types.SendTx
	Time time.Time
}

// Render renders the TxContents component
func (contents *TxContents) Render() vecty.ComponentOrHTML {
	return elem.Main(
		renderTxHeader(contents.Tx, contents.Time),
		renderTxContents(contents.Tx),
	)
}

func renderTxHeader(tx *types.SendTx, tm time.Time) vecty.ComponentOrHTML {
	return elem.Div(vecty.Markup(vecty.Class("card-deck-col")),
		elem.Div(vecty.Markup(vecty.Class("card")),
			elem.Div(
				vecty.Markup(vecty.Class("card-header")),
				elem.Anchor(
					vecty.Markup(
						vecty.Class("nav-link"),
						vecty.Attribute("href", "/txs/"+util.IntToString((tx.Height))),
					),
					vecty.Text("Transaction "+util.IntToString(tx.Height)),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("card-body")),
				elem.Div(
					vecty.Markup(vecty.Class("dt")),
					vecty.Text(humanize.Ordinal(int(tx.Store.Index+1))+" transaction on block "),
					elem.Anchor(
						vecty.Markup(
							vecty.Attribute("href", "/blocks/"+util.IntToString(tx.Store.Height)),
						),
						vecty.Text(util.IntToString(tx.Store.Height)),
					),
				),
				elem.Div(
					elem.Div(
						vecty.Markup(vecty.Class("dt")),
						vecty.Text("Hash"),
					),
					elem.Div(
						vecty.Markup(vecty.Class("dd")),
						vecty.Text(tx.Hash.String()),
					),
					elem.Div(
						vecty.Text(timeago.English.Format(tm)),
					),
				),
			),
		),
	)
}

func renderTxContents(tx *types.SendTx) vecty.ComponentOrHTML {
	result, err := json.MarshalIndent(tx.Store.TxResult, "", "    ")
	util.ErrPrint(err)
	var rawTx dvotetypes.Tx
	err = json.Unmarshal(tx.Store.Tx, &rawTx)
	util.ErrPrint(err)
	var txContents []byte
	switch rawTx.Type {
	case "newProcess":
		var newProcess dvotetypes.NewProcessTx
		err = json.Unmarshal(tx.Store.Tx, &newProcess)
		txContents, err = json.MarshalIndent(newProcess, "", "    ")
		util.ErrPrint(err)
	case "cancelProcess":
		var cancelProcess dvotetypes.CancelProcessTx
		err = json.Unmarshal(tx.Store.Tx, &cancelProcess)
		txContents, err = json.MarshalIndent(cancelProcess, "", "    ")
		util.ErrPrint(err)
	case "admin":
		var admin dvotetypes.AdminTx
		err = json.Unmarshal(tx.Store.Tx, &admin)
		txContents, err = json.MarshalIndent(admin, "", "    ")
		util.ErrPrint(err)
	}

	// txContents := base64.StdEncoding.EncodeToString(tx.Store.Tx)
	accordionName := "accordionTx"
	return elem.Div(
		vecty.Markup(vecty.Class("accordion"), prop.ID(accordionName)),
		renderCollapsible("Transaction Contents", string(txContents), accordionName, "One"),
		renderCollapsible("Transaction Results", string(result), accordionName, "Two"),
	)
}
