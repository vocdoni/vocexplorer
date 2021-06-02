package components

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hexops/vecty"
	"go.vocdoni.io/proto/build/go/models"
	"google.golang.org/protobuf/proto"

	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/logger"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// TransactionsDashboardView renders the dashboard landing page
type TransactionsDashboardView struct {
	vecty.Core
	vecty.Mounter
	Rendered bool
}

// Mount is called after the component renders to signal that it can be rerendered safely
func (dash *TransactionsDashboardView) Mount() {
	if !dash.Rendered {
		dash.Rendered = true
		vecty.Rerender(dash)
	}
}

// Render renders the TransactionsDashboardView component
func (dash *TransactionsDashboardView) Render() vecty.ComponentOrHTML {
	if !dash.Rendered {
		return LoadingBar()
	}
	return Container(
		vecty.Markup(vecty.Attribute("id", "main")),
		renderServerConnectionBanner(),
		&TransactionList{},
	)
}

// UpdateTransactionsDashboard keeps the Transactions dashboard updated
func UpdateTransactionsDashboard(d *TransactionsDashboardView) {
	dispatcher.Dispatch(&actions.EnableAllUpdates{})

	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	if !update.CheckCurrentPage("transactions", ticker) {
		return
	}
	updateTransactionsDashboard(d)
	for {
		select {
		case <-store.RedirectChan:
			if !update.CheckCurrentPage("transactions", ticker) {
				return
			}
		case <-ticker.C:
			if !update.CheckCurrentPage("transactions", ticker) {
				return
			}
			updateTransactionsDashboard(d)
		case i := <-store.Transactions.Pagination.PagChannel:
			if !update.CheckCurrentPage("transactions", ticker) {
				return
			}
		txloop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-store.Transactions.Pagination.PagChannel:
				default:
					break txloop
				}
			}
			dispatcher.Dispatch(&actions.TransactionsIndexChange{Index: i})
			logger.Info(fmt.Sprintf("update Transactions to index %d\n", i))
			updateTransactions(d, int(store.Transactions.Count)-store.Transactions.Pagination.Index-config.ListSize+1)
		}
	}
}

func updateTransactionsDashboard(d *TransactionsDashboardView) {
	if !store.Transactions.Pagination.DisableUpdate {
		stats, err := store.Client.GetStats()
		if err != nil {
			logger.Error(err)
			return
		}
		actions.UpdateCounts(stats)
		updateTransactions(d, int(store.Transactions.Count)-store.Transactions.Pagination.Index-config.ListSize+1)
	}
	dispatcher.Dispatch(&actions.GatewayConnected{GatewayErr: store.Client.GetGatewayInfo()})
}

func updateTransactions(d *TransactionsDashboardView, index int) {
	listSize := config.ListSize
	if index < 0 {
		listSize += index
		index = 0
	}
	logger.Info(fmt.Sprintf("Getting %d Transactions from index %d\n", listSize, index))
	list := []*storeutil.FullTransaction{}
	for i := 0; i < listSize; i++ {
		tx, err := store.Client.GetTxByHeight(uint32(index + i))
		if err != nil {
			logger.Error(err)
			return
		}
		if tx == nil {
			continue
		}
		list = append(list, &storeutil.FullTransaction{
			Decoded: decodeTransaction(tx.Tx),
			Package: tx,
		})
	}
	dispatcher.Dispatch(&actions.SetTransactionList{TransactionList: list})
}

func decodeTransaction(tx []byte) *storeutil.DecodedTransaction {
	var rawTx models.Tx
	if err := proto.Unmarshal(tx, &rawTx); err != nil {
		logger.Error(err)
	}
	var txContents string
	var processID string
	var nullifier string
	var entityID string

	switch rawTx.Payload.(type) {
	case *models.Tx_Vote:
		typedTx := rawTx.GetVote()
		if typedTx == nil {
			logger.Error(fmt.Errorf("vote transaction empty"))
			break
		}
		txBytes, err := json.MarshalIndent(typedTx, "", "\t")
		if err != nil {
			logger.Error(err)
		}
		txContents = string(txBytes)
		processID = hex.EncodeToString(typedTx.GetProcessId())
		txContents = convertB64ToHex(txContents, "nonce", hex.EncodeToString(typedTx.Nonce))
		txContents = convertB64ToHex(txContents, "processId", processID)
		// TODO decode more proof types
		if typedTx.GetProof() != nil {
			switch typedTx.GetProof().Payload.(type) {
			case *models.Proof_Graviton:
				txContents = convertB64ToHex(txContents, "siblings", hex.EncodeToString(
					typedTx.GetProof().GetGraviton().Siblings))
			case *models.Proof_EthereumStorage:
				siblings := []string{}
				for _, sibling := range typedTx.GetProof().GetEthereumStorage().Siblings {
					siblings = append(siblings, util.HexToString(sibling))
				}
				txContents = convertB64ToHex(txContents, "siblings", strings.Join(siblings, ", "))
			case *models.Proof_Iden3:
				txContents = convertB64ToHex(txContents, "siblings", hex.EncodeToString(
					typedTx.GetProof().GetIden3().GetSiblings()))
			default:
				logger.Info("Other proof type")
			}
		}
		txContents = convertB64ToHex(txContents, "votePackage", hex.EncodeToString(typedTx.VotePackage))

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

	return &storeutil.DecodedTransaction{
		RawTxContents: txContents,
		RawTx:         &rawTx,
		ProcessID:     processID,
		EntityID:      entityID,
		Nullifier:     nullifier,
	}
}
