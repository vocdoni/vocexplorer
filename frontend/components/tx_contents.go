package components

import (
	"encoding/json"
	"fmt"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/golang/protobuf/ptypes"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	abci "github.com/tendermint/tendermint/abci/types"
	"gitlab.com/vocdoni/go-dvote/log"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/proto"
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
			renderGatewayConnectionBanner(),
			renderServerConnectionBanner(),
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
	return Container(
		renderGatewayConnectionBanner(),
		renderServerConnectionBanner(),
		contents)
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
				elem.Span(
					vecty.Text(humanize.Ordinal(int(store.Transactions.CurrentTransaction.Store.Index+1))+" transaction on "),
					vecty.If(
						!proto.BlockIsEmpty(store.Transactions.CurrentBlock),
						Link(
							"/block/"+util.IntToString(store.Transactions.CurrentTransaction.Store.Height),
							"block "+util.IntToString(store.Transactions.CurrentTransaction.Store.Height),
							"",
						),
					),
				),
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
				elem.DefinitionTerm(
					vecty.Text("Hash"),
				),
				elem.Description(
					vecty.Text(util.HexToString(store.Transactions.CurrentTransaction.GetHash())),
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
		Text:  "Transaction Log",
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
			vecty.Text("Empty transaction log"),
		)
	}
	return elem.Preformatted(elem.Code(
		vecty.Text(string(store.Transactions.CurrentDecodedTransaction.Metadata)),
	))
}

// UpdateTxContents keeps the transaction contents up to date
func UpdateTxContents(d *TxContents) {
	dispatcher.Dispatch(&actions.EnableAllUpdates{})
	// Fetch transaction contents
	tx, ok := api.GetTx(store.Transactions.CurrentTransactionHeight)
	if ok && tx != nil {
		dispatcher.Dispatch(&actions.SetCurrentTransaction{Transaction: tx})
	} else {
		return
	}
	// Set block associated with transaction
	block, ok := api.GetStoreBlock(store.Transactions.CurrentTransaction.Store.Height)
	if ok {
		dispatcher.Dispatch(&actions.SetTransactionBlock{Block: block})
	}

	var txResult abci.ResponseDeliverTx
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
	if !proto.BlockIsEmpty(store.Transactions.CurrentBlock) {
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

	entityID = util.TrimHex(entityID)
	processID = util.TrimHex(processID)
	nullifier = util.TrimHex(nullifier)
	var envelopeHeight int64
	if nullifier != "" {
		envelopeHeight, ok = api.GetEnvelopeHeightFromNullifier(nullifier)
	}
	if !ok {
		log.Error("unable to retrieve envelope height, envelope may not exist")
	}
	var metadata []byte
	if !txResult.Equal(abci.ResponseDeliverTx{}) {
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
