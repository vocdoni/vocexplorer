package components

import (
	"fmt"
	"math/big"
	"regexp"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"

	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/logger"
	"gitlab.com/vocdoni/vocexplorer/util"
	"go.vocdoni.io/dvote/crypto/ethereum"
	"go.vocdoni.io/proto/build/go/models"
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
		return Unavailable("Transaction unavailable", "")
	}
	if store.Transactions.CurrentTransaction == nil {
		return Unavailable("Loading transaction...", "")
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
				"%s transaction on ", humanize.Ordinal(int(store.Transactions.CurrentTransaction.Index+1)),
			)),
			Link(
				"/block/"+util.IntToString(store.Transactions.CurrentTransaction.BlockHeight),
				"block "+util.IntToString(store.Transactions.CurrentTransaction.BlockHeight),
				"bold-link",
			),
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
					vecty.Text(util.GetTransactionName(util.GetTransactionType(store.Transactions.CurrentDecodedTransaction.RawTx))),
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
					store.Transactions.CurrentDecodedTransaction.Nullifier != "",
					elem.DefinitionTerm(
						vecty.Text("Contains vote envelope"),
					),
					elem.Description(
						Link(
							"/envelope/"+store.Transactions.CurrentDecodedTransaction.Nullifier,
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
		vecty.Text(store.Transactions.CurrentDecodedTransaction.RawTxContents),
	))
}

// UpdateTxContents keeps the transaction contents up to date
func UpdateTxContents(d *TxContents, blockHeight uint32, index int32) {
	dispatcher.Dispatch(&actions.SetCurrentTransaction{Transaction: nil})
	dispatcher.Dispatch(&actions.EnableAllUpdates{})
	d.fetchTransaction(blockHeight, index)
	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	if !update.CheckCurrentPage("tx", ticker) {
		return
	}
	for {
		select {
		case <-store.RedirectChan:
			if !update.CheckCurrentPage("tx", ticker) {
				return
			}
		case <-ticker.C:
			if !update.CheckCurrentPage("tx", ticker) {
				return
			}
			// If transaction never loaded, load it
			if d.Unavailable {
				d.fetchTransaction(blockHeight, index)
			}
		}
	}
}

func (t *TxContents) fetchTransaction(blockHeight uint32, index int32) {
	// Fetch transaction contents
	tx, err := store.Client.GetTx(blockHeight, index)
	if err == nil {
		t.Unavailable = false
		dispatcher.Dispatch(&actions.SetCurrentTransaction{Transaction: tx})
	} else {
		logger.Error(err)
		t.Unavailable = true
		dispatcher.Dispatch(&actions.SetCurrentTransaction{Transaction: nil})
		return
	}
	// Set block associated with transaction
	block, err := store.Client.GetBlock(blockHeight)
	if err != nil {
		logger.Error(err)
	} else {
		dispatcher.Dispatch(&actions.SetTransactionBlock{Block: block})
	}

	// decoded the transaction
	decoded := decodeTransaction(tx.Tx)
	if decoded == nil {
		dispatcher.Dispatch(&actions.SetCurrentDecodedTransaction{Transaction: nil})
		return
	}
	// If vote type, generate the nullifier as well
	switch decoded.RawTx.Payload.(type) {
	case *models.Tx_Vote:
		generateNullifier(decoded, tx.Tx, tx.Signature)
	}
	decoded.Time = store.Transactions.CurrentBlock.Timestamp
	dispatcher.Dispatch(&actions.SetCurrentDecodedTransaction{Transaction: decoded})
}

func convertB64ToHex(source, key, hex string) string {
	censusRootRe, err := regexp.Compile("\"" + key + "\": \".*\"")
	if err != nil {
		logger.Warn(err.Error())
	}
	if len(hex) > 0 {
		source = censusRootRe.ReplaceAllString(source, "\""+key+"\": \""+hex+"\"")
	}
	return source
}

func formatQuestions(votes []*models.QuestionResult) string {
	votesString := "\"votes\": [\n"
	for i, vote := range votes {
		voteString := "\t\t\t{\n\t\t\t\t\"question\": [\n"
		for j, question := range vote.Question {
			val := new(big.Int).SetBytes(question)
			voteString += "\t\t\t\t\t"
			voteString += val.String()
			if j < (len(vote.Question) - 1) {
				voteString += ",\n"
			}
		}
		voteString += "\n\t\t\t\t]\n\t\t\t}"
		if i < (len(votes) - 1) {
			voteString += ",\n"
		}
		votesString += voteString
	}
	votesString += "\n\t\t],"
	return votesString
}

func generateNullifier(decoded *storeutil.DecodedTransaction, tx, signature []byte) {
	pubKey, err := ethereum.PubKeyFromSignature(tx, signature)
	if err != nil {
		logger.Error(fmt.Errorf("cannot extract public key from signature: (%w)", err))
		return
	}
	addr, err := ethereum.AddrFromPublicKey(pubKey)
	if err != nil {
		logger.Error(fmt.Errorf("cannot extract address from public key: (%w)", err))
		return
	}

	// assign a nullifier
	decoded.Nullifier = util.HexToString(ethereum.HashRaw([]byte(fmt.Sprintf("%s%s", addr.Bytes(), decoded.ProcessID))))
	decoded.RawTxContents = convertB64ToHex(decoded.RawTxContents, "nullifier", decoded.Nullifier)
}
