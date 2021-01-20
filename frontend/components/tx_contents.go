package components

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/vocdoni/vocexplorer/api"
	"github.com/vocdoni/vocexplorer/api/dbtypes"
	"github.com/vocdoni/vocexplorer/frontend/actions"
	"github.com/vocdoni/vocexplorer/frontend/bootstrap"
	"github.com/vocdoni/vocexplorer/frontend/dispatcher"
	"github.com/vocdoni/vocexplorer/frontend/store"
	"github.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"github.com/vocdoni/vocexplorer/logger"
	"github.com/vocdoni/vocexplorer/util"
	"go.vocdoni.io/dvote/types"
	"go.vocdoni.io/proto/build/go/models"
	"google.golang.org/protobuf/proto"
)

// TxContents renders tx contents
type TxContents struct {
	vecty.Core
	vecty.Mounter
	Rendered    bool
	Unavailable bool
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
	if t.Unavailable {
		return Unavailable("Transaction unavailable")
	}
	if store.Transactions.CurrentTransaction == nil {
		return Unavailable("Loading transaction...")
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
		vecty.Markup(vecty.Attribute("id", "main")),
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
				store.Transactions.CurrentTransaction.TxHeight,
			)),
		),
	}

	if store.Transactions.CurrentDecodedTransaction != nil {
		logger.Info("nullifier: " + store.Transactions.CurrentDecodedTransaction.Nullifier)
		logger.Info("type: " + util.GetTransactionType(&store.Transactions.CurrentDecodedTransaction.RawTx))
		contents = append(contents, vecty.List{
			elem.Div(
				vecty.Markup(vecty.Class("details")),
				elem.Span(
					vecty.Text(humanize.Ordinal(int(store.Transactions.CurrentTransaction.Index+1))+" transaction on "),
					vecty.If(
						!dbtypes.BlockIsEmpty(store.Transactions.CurrentBlock),
						Link(
							"/block/"+util.IntToString(store.Transactions.CurrentTransaction.Height),
							"block "+util.IntToString(store.Transactions.CurrentTransaction.Height),
							"bold-link",
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
					vecty.Text(util.GetTransactionName(util.GetTransactionType(&store.Transactions.CurrentDecodedTransaction.RawTx))),
				),
				elem.DefinitionTerm(
					vecty.Text("Hash"),
				),
				elem.Description(
					vecty.Text(util.HexToString(store.Transactions.CurrentTransaction.Hash)),
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
					store.Transactions.CurrentDecodedTransaction.Nullifier != "" && util.GetTransactionType(&store.Transactions.CurrentDecodedTransaction.RawTx) == types.TxVote,
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

	return vecty.List{
		elem.Navigation(
			vecty.Markup(vecty.Attribute("aria-label", "Tab navigation")),
			vecty.Markup(vecty.Class("tabs")),
			elem.UnorderedList(
				TabLink(t, contents),
			),
		),
		elem.Div(
			vecty.Markup(vecty.Class("tabs-content")),
			TabContents(contents, preformattedTransactionContents()),
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

// UpdateTxContents keeps the transaction contents up to date
func UpdateTxContents(d *TxContents) {
	dispatcher.Dispatch(&actions.SetCurrentTransaction{Transaction: nil})
	dispatcher.Dispatch(&actions.EnableAllUpdates{})
	// Fetch transaction contents
	tx, ok := api.GetTxByHeight(store.Transactions.CurrentTransactionHeight)
	if ok && tx != nil {
		d.Unavailable = false
		dispatcher.Dispatch(&actions.SetCurrentTransaction{Transaction: tx})
	} else {
		d.Unavailable = true
		dispatcher.Dispatch(&actions.SetCurrentTransaction{Transaction: nil})
		return
	}
	// Set block associated with transaction
	block, ok := api.GetStoreBlock(store.Transactions.CurrentTransaction.Height)
	if ok {
		dispatcher.Dispatch(&actions.SetTransactionBlock{Block: block})
	}

	var rawTx models.Tx
	err := proto.Unmarshal(tx.Tx, &rawTx)
	if err != nil {
		logger.Error(err)
	}
	var txContents []byte
	var processID string
	var nullifier string
	var entityID string

	switch rawTx.Payload.(type) {
	case *models.Tx_Vote:
		typedTx := rawTx.GetVote()
		// typedTx.Nullifier = tx.Nullifier
		txContents, err = json.MarshalIndent(typedTx, "", "\t")
		if err != nil {
			logger.Error(err)
		}
		processID = hex.EncodeToString(typedTx.GetProcessId())
		nullifier = tx.Nullifier
	case *models.Tx_NewProcess:
		typedTx := rawTx.GetNewProcess()
		txContents, err = json.MarshalIndent(typedTx, "", "\t")
		if err != nil {
			logger.Error(err)
		}
		processID = hex.EncodeToString(typedTx.Process.GetProcessId())
		entityID = hex.EncodeToString(typedTx.Process.GetEntityId())
	case *models.Tx_Admin:
		typedTx := rawTx.GetAdmin()
		txContents, err = json.MarshalIndent(typedTx, "", "\t")
		if err != nil {
			logger.Error(err)
		}
		processID = hex.EncodeToString(typedTx.GetProcessId())
	case *models.Tx_SetProcess:
		typedTx := rawTx.GetSetProcess()
		txContents, err = json.MarshalIndent(typedTx, "", "\t")
		if err != nil {
			logger.Error(err)
		}
		processID = hex.EncodeToString(typedTx.GetProcessId())
	}

	entityID = util.TrimHex(entityID)
	processID = util.TrimHex(processID)
	nullifier = util.TrimHex(nullifier)
	logger.Info(fmt.Sprintf("generated Nullifier: %s", nullifier))
	var envelopeHeight int64
	if nullifier != "" {
		envelopeHeight, ok = api.GetEnvelopeHeightFromNullifier(nullifier)
	}
	if !ok {
		logger.Error(fmt.Errorf("unable to retrieve envelope height, envelope may not exist"))
	}
	dispatcher.Dispatch(&actions.SetCurrentDecodedTransaction{
		Transaction: &storeutil.DecodedTransaction{
			RawTxContents:  txContents,
			RawTx:          rawTx,
			Time:           block.Time,
			EnvelopeHeight: envelopeHeight,
			ProcessID:      processID,
			EntityID:       entityID,
			Nullifier:      nullifier,
		},
	})
}
