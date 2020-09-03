package components

import (
	"encoding/json"
	"fmt"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/golang/protobuf/ptypes"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"gitlab.com/vocdoni/go-dvote/log"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// TxContents renders tx contents
type TxContents struct {
	vecty.Core
	vecty.Mounter
	Rendered bool
}

// Mount triggers when TxContents renders
func (t *TxContents) Mount() {
	if !t.Rendered {
		t.Rendered = true
		vecty.Rerender(t)
	}
}

// Render renders the TxContents component
func (t *TxContents) Render() vecty.ComponentOrHTML {
	if !t.Rendered {
		return LoadingBar()
	}
	if store.Transactions.CurrentTransaction == nil {
		return Container(
			elem.Section(
				bootstrap.Card(bootstrap.CardParams{
					Body: vecty.List{
						elem.Heading3(
							vecty.Text("Transaction does not exist"),
						),
					},
				}),
			),
		)
	}
	contents := vecty.List{
		elem.Section(
			vecty.Markup(vecty.Class("details-view", "no-column")),
			elem.Div(
				vecty.Markup(vecty.Class("row")),
				elem.Div(
					vecty.Markup(vecty.Class("main-column")),
					bootstrap.Card(bootstrap.CardParams{
						Body: TransactionView(),
					}),
				),
			),
		),
	}

	if store.Transactions.CurrentDecodedTransaction != nil {
		contents = append(contents, elem.Section(
			vecty.Markup(vecty.Class("row")),
			elem.Div(
				vecty.Markup(vecty.Class("col-12")),
				bootstrap.Card(bootstrap.CardParams{
					Body: t.TransactionDetails(),
				}),
			),
		))
	}
	return Container(contents)
}

//TransactionView renders a single transaction card with (most of) the tx information
func TransactionView() vecty.List {
	contents := vecty.List{
		elem.Heading1(
			vecty.Markup(vecty.Class("card-title")),
			vecty.Text("Transaction details"),
		),
		elem.Heading2(
			vecty.Text(fmt.Sprintf(
				"Transaction height: %d",
				store.Transactions.CurrentTransaction.Store.TxHeight,
			)),
		),
	}

	if store.Transactions.CurrentDecodedTransaction != nil {
		contents = append(contents, vecty.List{
			elem.Div(
				vecty.Markup(vecty.Class("details")),
				elem.Span(vecty.Text(fmt.Sprintf(
					"%s (%s)",
					humanize.Time(store.Transactions.CurrentDecodedTransaction.Time),
					store.Transactions.CurrentDecodedTransaction.Time.Local().String(),
				))),
			),
			elem.HorizontalRule(),
			elem.DescriptionList(
				elem.DefinitionTerm(
					vecty.Text("Transaction Type"),
				),
				elem.Description(
					vecty.Text(store.Transactions.CurrentDecodedTransaction.RawTx.Type),
				),
				vecty.If(
					store.Transactions.CurrentDecodedTransaction.EntityID != "",
					vecty.List{
						elem.DefinitionTerm(
							vecty.Text("Belongs to entity"),
						),
						elem.Description(
							Link(
								"/entity/"+store.Transactions.CurrentDecodedTransaction.EntityID,
								store.Transactions.CurrentDecodedTransaction.EntityID,
								"",
							),
						),
					},
				),
				vecty.If(
					store.Transactions.CurrentDecodedTransaction.ProcessID != "",
					vecty.List{
						elem.DefinitionTerm(
							vecty.Text("Belongs to process"),
						),
						elem.Description(
							Link(
								"/process/"+store.Transactions.CurrentDecodedTransaction.ProcessID,
								store.Transactions.CurrentDecodedTransaction.ProcessID,
								"",
							),
						),
					},
				),
				vecty.If(
					store.Transactions.CurrentDecodedTransaction.Nullifier != "" && store.Transactions.CurrentDecodedTransaction.RawTx.Type == "vote",
					elem.DefinitionTerm(
						vecty.Text("Contains vote envelope"),
					),
					elem.Description(
						Link(
							"/envelope/"+util.IntToString(store.Transactions.CurrentDecodedTransaction.EnvelopeHeight),
							store.Transactions.CurrentDecodedTransaction.Nullifier,
							"",
						),
					),
				),
			),
		}...)
	}

	return contents
}

// TransactionTab records the current active tab for the transaction page
type TransactionTab struct {
	*Tab
}

func (t *TransactionTab) store() string {
	return store.Transactions.Pagination.Tab
}
func (t *TransactionTab) dispatch() interface{} {
	return &actions.TransactionTabChange{
		Tab: t.alias(),
	}
}

// TransactionDetails displays the transaction details pane for a single transaction
func (t *TxContents) TransactionDetails() vecty.ComponentOrHTML {
	contents := &TransactionTab{&Tab{
		Text:  "Contents",
		Alias: "contents",
	}}
	metadata := &TransactionTab{&Tab{
		Text:  "Metadata",
		Alias: "metadata",
	}}

	return vecty.List{
		elem.Navigation(
			vecty.Markup(vecty.Class("tabs")),
			elem.UnorderedList(
				TabLink(t, contents),
				TabLink(t, metadata),
			),
		),
		elem.Div(
			vecty.Markup(vecty.Class("tabs-content")),
			TabContents(contents, preformattedTransactionContents()),
			TabContents(metadata, preformattedTransactionMetadata()),
		),
	}
}

func preformattedTransactionContents() vecty.ComponentOrHTML {
	if len(store.Transactions.CurrentDecodedTransaction.RawTxContents) <= 0 {
		return elem.Preformatted(
			vecty.Markup(vecty.Class("empty")),
			vecty.Text("Empty contents"),
		)
	}
	return elem.Preformatted(elem.Code(
		vecty.Text(string(store.Transactions.CurrentDecodedTransaction.RawTxContents)),
	))
}

func preformattedTransactionMetadata() vecty.ComponentOrHTML {
	if len(store.Transactions.CurrentDecodedTransaction.Metadata) <= 0 {
		return elem.Preformatted(
			vecty.Markup(vecty.Class("empty")),
			vecty.Text("Empty metadata"),
		)
	}
	return elem.Preformatted(elem.Code(
		vecty.Text(string(store.Transactions.CurrentDecodedTransaction.Metadata)),
	))
}

// UpdateAndRenderTxContents keeps the transaction contents up to date
func UpdateAndRenderTxContents(d *TxContents) {
	actions.EnableUpdates()
	// Fetch transaction contents
	tx, ok := api.GetTx(store.Transactions.CurrentTransactionHeight)
	if ok {
		dispatcher.Dispatch(&actions.SetCurrentTransaction{Transaction: tx})
	}
	// Set block associated with transaction
	block, ok := api.GetBlock(store.Transactions.CurrentTransaction.Store.Height)
	if ok {
		dispatcher.Dispatch(&actions.SetTransactionBlock{Block: block})
	}

	var txResult coretypes.ResultTx
	err := json.Unmarshal(tx.GetStore().GetTxResult(), &txResult)
	if err != nil {
		log.Error(err)
	}

	var rawTx dvotetypes.Tx
	err = json.Unmarshal(tx.Store.Tx, &rawTx)
	if err != nil {
		log.Error(err)
	}
	var txContents []byte
	var processID string
	var nullifier string
	var entityID string
	var tm time.Time
	if !types.BlockIsEmpty(store.Transactions.CurrentBlock) {
		tm, err = ptypes.Timestamp(store.Transactions.CurrentBlock.GetTime())
		if err != nil {
			log.Error(err)
		}
	}

	switch rawTx.Type {
	case "vote":
		var typedTx dvotetypes.VoteTx
		err = json.Unmarshal(tx.Store.Tx, &typedTx)
		if err != nil {
			log.Error(err)
		}
		typedTx.Nullifier = tx.Store.Nullifier

		// TODO decrypt votes
		// Decode vote package if vote is unencrypted
		// if len(typedTx.EncryptionKeyIndexes) == 0 {
		// 	var vote dvotetypes.VotePackage
		// 	rawVote, err := base64.StdEncoding.DecodeString(typedTx.VotePackage)
		// 	if err != nil {
		// log.Error(err)
		// 		txContents, err = json.MarshalIndent(typedTx, "", "\t")
		// 		break
		// 	}
		// 	err = json.Unmarshal(rawVote, &vote)
		// 	if err != nil {
		// log.Error(err)
		// 		txContents, err = json.MarshalIndent(typedTx, "", "\t")
		// 		break
		// 	}
		// 	voteIndent, err := json.MarshalIndent(vote, "", "\t")
		// 	typedTx.VotePackage = string(voteIndent)
		// }

		txContents, err = json.MarshalIndent(typedTx, "", "\t")
		if err != nil {
			log.Error(err)
		}
		processID = typedTx.ProcessID
		nullifier = typedTx.Nullifier
	case "newProcess":
		var typedTx dvotetypes.NewProcessTx
		err = json.Unmarshal(tx.Store.Tx, &typedTx)
		if err != nil {
			log.Error(err)
		}
		txContents, err = json.MarshalIndent(typedTx, "", "\t")
		if err != nil {
			log.Error(err)
		}
		processID = typedTx.ProcessID
		entityID = typedTx.EntityID
	case "cancelProcess":
		var typedTx dvotetypes.CancelProcessTx
		err = json.Unmarshal(tx.Store.Tx, &typedTx)
		if err != nil {
			log.Error(err)
		}
		txContents, err = json.MarshalIndent(typedTx, "", "\t")
		if err != nil {
			log.Error(err)
		}
		processID = typedTx.ProcessID
	case "admin", "addValidator", "removeValidator", "addOracle", "removeOracle", "addProcessKeys", "revealProcessKeys":
		var typedTx dvotetypes.AdminTx
		err = json.Unmarshal(tx.Store.Tx, &typedTx)
		if err != nil {
			log.Error(err)
		}
		txContents, err = json.MarshalIndent(typedTx, "", "\t")
		if err != nil {
			log.Error(err)
		}
		processID = typedTx.ProcessID
	}

	entityID = util.StripHexString(entityID)
	processID = util.StripHexString(processID)
	nullifier = util.StripHexString(nullifier)
	var envelopeHeight int64
	if nullifier != "" {
		envelopeHeight, ok = api.GetEnvelopeHeightFromNullifier(nullifier)
	}
	var metadata []byte
	if len(txResult.Hash.Bytes()) > 0 && txResult.Height > 0 && len(txResult.Tx.Hash()) > 0 {
		metadata, err = json.MarshalIndent(txResult, "", "\t")
		if err != nil {
			log.Error(err)
		}
	}
	dispatcher.Dispatch(&actions.SetCurrentDecodedTransaction{
		Transaction: &storeutil.DecodedTransaction{
			Metadata:       metadata,
			RawTxContents:  txContents,
			RawTx:          rawTx,
			Time:           tm,
			EnvelopeHeight: envelopeHeight,
			ProcessID:      processID,
			EntityID:       entityID,
			Nullifier:      nullifier,
		},
	})
}

//TODO: link to envelope. Possibly store envelope nullifier/height in tx

func (t *TxContents) renderFullTx() vecty.ComponentOrHTML {
	accordionName := "accordionTx"

	if store.Transactions.CurrentDecodedTransaction == nil {
		return bootstrap.Card(bootstrap.CardParams{
			Header: vecty.Text("Transaction " + util.IntToString(store.Transactions.CurrentTransaction.Store.TxHeight)),
			Body: elem.Div(
				elem.Div(
					vecty.Markup(vecty.Class("dt")),
					vecty.Text(humanize.Ordinal(int(store.Transactions.CurrentTransaction.Store.Index+1))+" transaction on block "),
					vecty.If(
						!types.BlockIsEmpty(store.Transactions.CurrentBlock),
						Link(
							"/block/"+util.IntToString(store.Transactions.CurrentTransaction.Store.Height),
							util.IntToString(store.Transactions.CurrentTransaction.Store.Height),
							"",
						),
					),
					vecty.If(
						types.BlockIsEmpty(store.Transactions.CurrentBlock),
						vecty.Text(util.IntToString(store.Transactions.CurrentTransaction.Store.Height)+" (block not yet available)"),
					),
				),
			),
		},
		)
	}

	return bootstrap.Card(bootstrap.CardParams{
		Header: vecty.Text("Transaction " + util.IntToString(store.Transactions.CurrentTransaction.Store.TxHeight)),
		Body: elem.Div(
			elem.Div(
				vecty.Markup(vecty.Class("dt")),
				vecty.Text(humanize.Ordinal(int(store.Transactions.CurrentTransaction.Store.Index+1))+" transaction on block "),
				vecty.If(
					!types.BlockIsEmpty(store.Transactions.CurrentBlock),
					Link(
						"/block/"+util.IntToString(store.Transactions.CurrentTransaction.Store.Height),
						util.IntToString(store.Transactions.CurrentTransaction.Store.Height),
						"",
					),
				),
				vecty.If(
					types.BlockIsEmpty(store.Transactions.CurrentBlock),
					vecty.Text(util.IntToString(store.Transactions.CurrentTransaction.Store.Height)+" (block not yet available)"),
				),
			),
			elem.Div(
				elem.Div(
					vecty.Markup(vecty.Class("dt")),
					vecty.Text("Hash"),
				),
				elem.Div(
					vecty.Markup(vecty.Class("dd")),
					vecty.Text(util.HexToString(store.Transactions.CurrentTransaction.GetHash())),
				),
				vecty.If(
					!store.Transactions.CurrentDecodedTransaction.Time.IsZero(),
					elem.Div(
						vecty.Text(humanize.Time(store.Transactions.CurrentDecodedTransaction.Time)),
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
					vecty.Text(store.Transactions.CurrentDecodedTransaction.RawTx.Type),
				),
			),
			vecty.If(
				store.Transactions.CurrentDecodedTransaction.EntityID != "",
				elem.Div(
					vecty.Text("Belongs to entity "),
					Link(
						"/entity/"+store.Transactions.CurrentDecodedTransaction.EntityID,
						store.Transactions.CurrentDecodedTransaction.EntityID,
						"",
					),
				),
			),
			vecty.If(
				store.Transactions.CurrentDecodedTransaction.ProcessID != "",
				elem.Div(
					vecty.Text("Belongs to process "),
					Link(
						"/process/"+store.Transactions.CurrentDecodedTransaction.ProcessID,
						store.Transactions.CurrentDecodedTransaction.ProcessID,
						"",
					),
				),
			),
			vecty.If(
				store.Transactions.CurrentDecodedTransaction.Nullifier != "" && store.Transactions.CurrentDecodedTransaction.RawTx.Type == "vote",
				elem.Div(
					vecty.Text("Contains vote envelope "),
					Link(
						"/envelope/"+util.IntToString(store.Transactions.CurrentDecodedTransaction.EnvelopeHeight),
						store.Transactions.CurrentDecodedTransaction.Nullifier,
						"",
					),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("accordion"), prop.ID(accordionName)),
				vecty.If(
					len(store.Transactions.CurrentDecodedTransaction.RawTxContents) > 0,
					renderCollapsible("Transaction Contents", accordionName, "One", elem.Preformatted(vecty.Text(string(store.Transactions.CurrentDecodedTransaction.RawTxContents)))),
				),
				vecty.If(
					len(store.Transactions.CurrentDecodedTransaction.Metadata) > 0,
					renderCollapsible("Transaction MetaData", accordionName, "Two", elem.Preformatted(vecty.Text(string(store.Transactions.CurrentDecodedTransaction.Metadata)))),
				),
			),
		),
	})
}
