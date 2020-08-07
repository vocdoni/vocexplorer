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
	Tx       *types.SendTx
	Time     time.Time
	HasBlock bool
}

// Render renders the TxContents component
func (contents *TxContents) Render() vecty.ComponentOrHTML {
	return elem.Main(
		tenderFullTx(contents.Tx, contents.Time, contents.HasBlock),
	)
}

func tenderFullTx(tx *types.SendTx, tm time.Time, hasBlock bool) vecty.ComponentOrHTML {
	result, err := json.MarshalIndent(tx.Store.TxResult, "", "    ")
	util.ErrPrint(err)
	var rawTx dvotetypes.Tx
	err = json.Unmarshal(tx.Store.Tx, &rawTx)
	util.ErrPrint(err)
	var txContents []byte
	var processID string
	var nullifier string
	var entityID string

	switch rawTx.Type {
	case "vote":
		var typedTx dvotetypes.VoteTx
		err = json.Unmarshal(tx.Store.Tx, &typedTx)
		txContents, err = json.MarshalIndent(typedTx, "", "    ")
		util.ErrPrint(err)
		processID = typedTx.ProcessID
		nullifier = typedTx.Nullifier
	case "newProcess":
		var typedTx dvotetypes.NewProcessTx
		err = json.Unmarshal(tx.Store.Tx, &typedTx)
		txContents, err = json.MarshalIndent(typedTx, "", "    ")
		util.ErrPrint(err)
		processID = typedTx.ProcessID
		entityID = typedTx.EntityID
	case "cancelProcess":
		var typedTx dvotetypes.CancelProcessTx
		err = json.Unmarshal(tx.Store.Tx, &typedTx)
		txContents, err = json.MarshalIndent(typedTx, "", "    ")
		util.ErrPrint(err)
		processID = typedTx.ProcessID
	case "admin", "addValidator", "removeValidator", "addOracle", "removeOracle", "addProcessKeys", "revealProcessKeys":
		var typedTx dvotetypes.AdminTx
		err = json.Unmarshal(tx.Store.Tx, &typedTx)
		txContents, err = json.MarshalIndent(typedTx, "", "    ")
		util.ErrPrint(err)
		processID = typedTx.ProcessID
	}

	util.StripHexString(&entityID)
	util.StripHexString(&processID)
	util.StripHexString(&nullifier)

	// txContents := base64.StdEncoding.EncodeToString(tx.Store.Tx)
	accordionName := "accordionTx"
	return elem.Div(elem.Div(vecty.Markup(vecty.Class("card-deck-col")),
		elem.Div(vecty.Markup(vecty.Class("card")),
			elem.Div(
				vecty.Markup(vecty.Class("card-header")),
				elem.Anchor(
					vecty.Markup(
						vecty.Class("nav-link"),
						vecty.Attribute("href", "/txs/"+util.IntToString((tx.Store.TxHeight))),
					),
					vecty.Text("Transaction "+util.IntToString(tx.Store.TxHeight)),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("card-body")),
				elem.Div(
					vecty.Markup(vecty.Class("dt")),
					vecty.Text(humanize.Ordinal(int(tx.Store.Index+1))+" transaction on block "),
					vecty.If(
						hasBlock,
						elem.Anchor(
							vecty.Markup(
								vecty.Attribute("href", "/blocks/"+util.IntToString(tx.Store.Height)),
							),
							vecty.Text(util.IntToString(tx.Store.Height)),
						),
					),
					vecty.If(
						!hasBlock,
						vecty.Text(util.IntToString(tx.Store.Height)+" (block not yet available)"),
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
					vecty.If(
						tm.IsZero(),
						elem.Div(
							vecty.Text(timeago.English.Format(tm)),
						),
					),
				),
				elem.Div(
					elem.Div(
						vecty.Markup(vecty.Class("dt")),
						vecty.Text("Transaction Type"),
					),
					elem.Div(
						vecty.Markup(vecty.Class("dd")),
						vecty.Text(rawTx.Type),
					),
				),
				vecty.If(
					entityID != "",
					elem.Div(
						vecty.Text("Belongs to entity "),
						elem.Anchor(
							vecty.Markup(
								vecty.Attribute("href", "/entities/"+entityID),
							),
							vecty.Text(entityID),
						),
					),
				),
				vecty.If(
					processID != "",
					elem.Div(
						vecty.Text("Belongs to process "),
						elem.Anchor(
							vecty.Markup(
								vecty.Attribute("href", "/processes/"+processID),
							),
							vecty.Text(processID),
						),
					),
				),
				vecty.If(
					nullifier != "",
					elem.Div(
						vecty.Text("Belongs to envelope "),
						elem.Anchor(
							vecty.Markup(
								vecty.Attribute("href", "/envelopes/"+nullifier),
							),
							vecty.Text(nullifier),
						),
					),
				),
			),
		),
	),
		elem.Div(
			vecty.Markup(vecty.Class("accordion"), prop.ID(accordionName)),
			renderCollapsible("Transaction Contents", accordionName, "One", elem.Preformatted(vecty.Text(string(txContents)))),
			renderCollapsible("Transaction Results", accordionName, "Two", elem.Preformatted(vecty.Text(string(result)))),
		),
	)
}
