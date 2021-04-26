package components

import (
	"encoding/hex"
	"encoding/json"
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
	"gitlab.com/vocdoni/vocexplorer/logger"
	"gitlab.com/vocdoni/vocexplorer/util"
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
				"Transaction index: %d on block: ", store.Transactions.CurrentTransactionRef.TxIndex,
			)),
			Link(
				"/block/"+util.IntToString(store.Transactions.CurrentTransactionRef.BlockHeight),
				"block "+util.IntToString(store.Transactions.CurrentTransactionRef.BlockHeight),
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
					vecty.Text(util.HexToString(store.Transactions.CurrentDecodedTransaction.Hash)),
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
					store.Transactions.CurrentDecodedTransaction.Nullifier != "" && util.GetTransactionType(store.Transactions.CurrentDecodedTransaction.RawTx) == types.TxVote,
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
func UpdateTxContents(d *TxContents) {
	dispatcher.Dispatch(&actions.SetCurrentTransaction{Transaction: nil})
	dispatcher.Dispatch(&actions.EnableAllUpdates{})
	// Fetch transaction contents
	tx, err := store.Client.GetTx(store.Transactions.CurrentTransactionRef.BlockHeight, store.Transactions.CurrentTransactionRef.TxIndex)
	if err == nil {
		d.Unavailable = false
		dispatcher.Dispatch(&actions.SetCurrentTransaction{Transaction: tx})
	} else {
		logger.Error(err)
		d.Unavailable = true
		dispatcher.Dispatch(&actions.SetCurrentTransaction{Transaction: nil})
		return
	}
	// Set block associated with transaction
	block, err := store.Client.GetBlock(store.Transactions.CurrentTransactionRef.BlockHeight)
	if err != nil {
		logger.Error(err)
	} else {
		dispatcher.Dispatch(&actions.SetTransactionBlock{Block: block})
	}

	var rawTx models.Tx
	err = proto.Unmarshal(tx.Tx, &rawTx)
	if err != nil {
		logger.Error(err)
	}
	var txContents string
	var processID string
	var nullifier string
	var entityID string

	switch rawTx.Payload.(type) {
	case *models.Tx_Vote:
		typedTx := rawTx.GetVote()
		// typedTx.Nullifier = tx.Nullifier

		txBytes, err := json.MarshalIndent(typedTx, "", "\t")
		if err != nil {
			logger.Error(err)
		}
		txContents = string(txBytes)
		processID = hex.EncodeToString(typedTx.GetProcessId())
		nullifier = util.HexToString(typedTx.Nullifier)
		txContents = convertB64ToHex(txContents, "nonce", hex.EncodeToString(typedTx.Nonce))
		txContents = convertB64ToHex(txContents, "processId", processID)
		txContents = convertB64ToHex(txContents, "siblings", hex.EncodeToString(typedTx.GetProof().GetGraviton().Siblings))
		txContents = convertB64ToHex(txContents, "votePackage", hex.EncodeToString(typedTx.VotePackage))
		txContents = convertB64ToHex(txContents, "nullifier", nullifier)
	case *models.Tx_NewProcess:
		typedTx := rawTx.GetNewProcess()
		txBytes, err := json.MarshalIndent(typedTx, "", "\t")
		if err != nil {

			logger.Error(err)
		}
		txContents = string(txBytes)
		processID = hex.EncodeToString(typedTx.Process.GetProcessId())
		entityID = hex.EncodeToString(typedTx.Process.GetEntityId())
		txContents = convertB64ToHex(txContents, "nonce", hex.EncodeToString(typedTx.Nonce))
		txContents = convertB64ToHex(txContents, "processId", processID)
		txContents = convertB64ToHex(txContents, "entityId", entityID)
		txContents = convertB64ToHex(txContents, "censusRoot", hex.EncodeToString(typedTx.Process.CensusRoot))
		txContents = convertB64ToHex(txContents, "paramsSignature", hex.EncodeToString(typedTx.Process.ParamsSignature))
	case *models.Tx_Admin:
		typedTx := rawTx.GetAdmin()
		txBytes, err := json.MarshalIndent(typedTx, "", "\t")
		if err != nil {
			logger.Error(err)
		}
		txContents = string(txBytes)
		processID = hex.EncodeToString(typedTx.GetProcessId())
		txContents = convertB64ToHex(txContents, "processId", processID)
		txContents = convertB64ToHex(txContents, "address", hex.EncodeToString(typedTx.Address))
		txContents = convertB64ToHex(txContents, "commitmentKey", hex.EncodeToString(typedTx.CommitmentKey))
		txContents = convertB64ToHex(txContents, "encryptionPrivateKey", hex.EncodeToString(typedTx.EncryptionPrivateKey))
		txContents = convertB64ToHex(txContents, "encryptionPublicKey", hex.EncodeToString(typedTx.EncryptionPublicKey))
		txContents = convertB64ToHex(txContents, "publicKey", hex.EncodeToString(typedTx.PublicKey))
		txContents = convertB64ToHex(txContents, "revealKey", hex.EncodeToString(typedTx.RevealKey))
		txContents = convertB64ToHex(txContents, "nonce", hex.EncodeToString(typedTx.Nonce))
	case *models.Tx_SetProcess:
		typedTx := rawTx.GetSetProcess()
		txBytes, err := json.MarshalIndent(typedTx, "", "\t")
		if err != nil {
			logger.Error(err)
		}
		txContents = string(txBytes)
		processID = hex.EncodeToString(typedTx.GetProcessId())
		txContents = convertB64ToHex(txContents, "nonce", hex.EncodeToString(typedTx.Nonce))
		txContents = convertB64ToHex(txContents, "processId", processID)
		txContents = convertB64ToHex(txContents, "censusRoot", hex.EncodeToString(typedTx.CensusRoot))
		if typedTx.GetResults() != nil {
			if len(typedTx.GetResults().EntityId) > 0 {
				entityID = hex.EncodeToString(typedTx.GetResults().EntityId)
				txContents = convertB64ToHex(txContents, "entityId", entityID)
			}
			votesRe, err := regexp.Compile("\"votes\":[^_]*\\],")
			if err != nil {
				logger.Warn(err.Error())
			}
			txContents = votesRe.ReplaceAllString(txContents, formatQuestions(typedTx.Results.Votes))
		}
	}

	entityID = util.TrimHex(entityID)
	processID = util.TrimHex(processID)
	nullifier = util.TrimHex(nullifier)
	dispatcher.Dispatch(&actions.SetCurrentDecodedTransaction{
		Transaction: &storeutil.DecodedTransaction{
			RawTxContents: txContents,
			RawTx:         &rawTx,
			Time:          time.Unix(store.Transactions.CurrentBlock.Timestamp, 0),
			ProcessID:     processID,
			EntityID:      entityID,
			Nullifier:     nullifier,
		},
	})
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
