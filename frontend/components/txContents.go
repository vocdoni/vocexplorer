package components

import (
	"encoding/json"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// TxContents renders tx contents
type TxContents struct {
	vecty.Core
	HasBlock bool
	Time     time.Time
	Tx       *types.SendTx
}

// Render renders the TxContents component
func (contents *TxContents) Render() vecty.ComponentOrHTML {
	return elem.Main(
		renderFullTx(contents.Tx, contents.Time, contents.HasBlock),
	)
}

//TODO: link to envelope. Possibly store envelope nullifier/height in tx

func renderFullTx(tx *types.SendTx, tm time.Time, hasBlock bool) vecty.ComponentOrHTML {
	var txResult coretypes.ResultTx
	err := json.Unmarshal(tx.GetStore().GetTxResult(), &txResult)
	result, err := json.MarshalIndent(txResult, "", "\t")
	util.ErrPrint(err)
	var rawTx dvotetypes.Tx
	err = json.Unmarshal(tx.Store.Tx, &rawTx)
	util.ErrPrint(err)
	var txContents []byte
	var processID string
	// var nullifier string
	var entityID string

	switch rawTx.Type {
	case "vote":
		var typedTx dvotetypes.VoteTx
		err = json.Unmarshal(tx.Store.Tx, &typedTx)
		util.ErrPrint(err)

		// TODO decrypt votes
		// Decode vote package if vote is unencrypted
		// if len(typedTx.EncryptionKeyIndexes) == 0 {
		// 	var vote dvotetypes.VotePackage
		// 	rawVote, err := base64.StdEncoding.DecodeString(typedTx.VotePackage)
		// 	if util.ErrPrint(err) {
		// 		txContents, err = json.MarshalIndent(typedTx, "", "\t")
		// 		break
		// 	}
		// 	err = json.Unmarshal(rawVote, &vote)
		// 	if util.ErrPrint(err) {
		// 		txContents, err = json.MarshalIndent(typedTx, "", "\t")
		// 		break
		// 	}
		// 	voteIndent, err := json.MarshalIndent(vote, "", "\t")
		// 	typedTx.VotePackage = string(voteIndent)
		// }

		txContents, err = json.MarshalIndent(typedTx, "", "\t")
		util.ErrPrint(err)
		processID = typedTx.ProcessID
		// nullifier = typedTx.Nullifier
	case "newProcess":
		var typedTx dvotetypes.NewProcessTx
		err = json.Unmarshal(tx.Store.Tx, &typedTx)
		util.ErrPrint(err)
		txContents, err = json.MarshalIndent(typedTx, "", "\t")
		util.ErrPrint(err)
		processID = typedTx.ProcessID
		entityID = typedTx.EntityID
	case "cancelProcess":
		var typedTx dvotetypes.CancelProcessTx
		err = json.Unmarshal(tx.Store.Tx, &typedTx)
		util.ErrPrint(err)
		txContents, err = json.MarshalIndent(typedTx, "", "\t")
		util.ErrPrint(err)
		processID = typedTx.ProcessID
	case "admin", "addValidator", "removeValidator", "addOracle", "removeOracle", "addProcessKeys", "revealProcessKeys":
		var typedTx dvotetypes.AdminTx
		err = json.Unmarshal(tx.Store.Tx, &typedTx)
		util.ErrPrint(err)
		txContents, err = json.MarshalIndent(typedTx, "", "\t")
		util.ErrPrint(err)
		processID = typedTx.ProcessID
	}

	entityID = util.StripHexString(entityID)
	processID = util.StripHexString(processID)
	// nullifier = util.StripHexString(nullifier)

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
						vecty.Text(util.HexToString(tx.GetHash())),
					),
					vecty.If(
						!tm.IsZero(),
						elem.Div(
							vecty.Text(humanize.Time(tm)),
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
				// vecty.If(
				// 	nullifier != "" && rawTx.Type == "vote",
				// 	elem.Div(
				// 		vecty.Text("Contains vote envelope "),
				// 		elem.Anchor(
				// 			vecty.Markup(
				// 				vecty.Attribute("href", "/envelopes/"+nullifier),
				// 			),
				// 			vecty.Text(nullifier),
				// 		),
				// 	),
				// ),
			),
		),
	),
		elem.Div(
			vecty.Markup(vecty.Class("accordion"), prop.ID(accordionName)),
			renderCollapsible("Transaction Contents", accordionName, "One", elem.Preformatted(vecty.Text(string(txContents)))),
			renderCollapsible("Transaction MetaData", accordionName, "Two", elem.Preformatted(vecty.Text(string(result)))),
		),
	)
}
