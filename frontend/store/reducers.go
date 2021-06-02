package store

import (
	"fmt"

	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/logger"
)

// Reduce registers all storage actions
func Reduce() {
	dispatcher.Register(blockActions)
	dispatcher.Register(clientActions)
	dispatcher.Register(configActions)
	dispatcher.Register(validatorActions)
	dispatcher.Register(transactionActions)
	dispatcher.Register(statsActions)
	dispatcher.Register(redirectActions)
	dispatcher.Register(processActions)
	dispatcher.Register(envelopeActions)
	dispatcher.Register(entityActions)
	dispatcher.Register(disableUpdateActions)
}

// configActions is the handler for all config-related actions
func configActions(action interface{}) {
	switch a := action.(type) {
	case *actions.StoreConfig:
		Config = a.Config

	default:
		return // don't fire listeners
	}

	Listeners.Fire()
}

// entityActions is the handler for all entity-related store actions
func entityActions(action interface{}) {
	switch a := action.(type) {
	case *actions.EntitiesIndexChange:
		Entities.Pagination.Index = a.Index

	case *actions.EntityProcessesIndexChange:
		Entities.ProcessPagination.Index = a.Index

	case *actions.SetEntityIDs:
		Entities.EntityIDs = a.EntityIDs

	case *actions.SetCurrentEntityID:
		Entities.CurrentEntityID = a.EntityID

	case *actions.EntityTabChange:
		Entities.Pagination.Tab = a.Tab

	case *actions.SetEntityCount:
		Entities.Count = a.Count

	case *actions.SetCurrentEntityProcessCount:
		Entities.CurrentEntity.ProcessCount = a.Count

	case *actions.SetEntityProcessCount:
		Entities.ProcessHeights[a.EntityID] = a.Count

	case *actions.SetEntityProcessIds:
		Entities.CurrentEntity.ProcessIds = a.ProcessList
		Entities.CurrentEntity.ProcessCount = len(a.ProcessList)

	default:
		return // don't fire listeners
	}

	Listeners.Fire()
}

// envelopeActions is the handler for all envelope-related store actions
func envelopeActions(action interface{}) {
	switch a := action.(type) {
	case *actions.EnvelopesIndexChange:
		Envelopes.Pagination.Index = a.Index

	case *actions.SetEnvelopeList:
		Envelopes.Envelopes = a.EnvelopeList

	case *actions.SetEnvelopeCount:
		Envelopes.Count = a.Count

	case *actions.SetCurrentEnvelope:
		Envelopes.CurrentEnvelope = a.Envelope

	case *actions.SetCurrentEnvelopeNullifier:
		Envelopes.CurrentEnvelopeNullifier = a.Nullifier

	case *actions.EnvelopesTabChange:
		Envelopes.Pagination.Tab = a.Tab

	default:
		return // don't fire listeners
	}

	Listeners.Fire()
}

// processActions is the handler for all process-related store actions
func processActions(action interface{}) {
	switch a := action.(type) {
	case *actions.ProcessesIndexChange:
		Processes.Pagination.Index = a.Index

	case *actions.ProcessEnvelopesIndexChange:
		Processes.EnvelopePagination.Index = a.Index

	case *actions.SetProcessIds:
		Processes.ProcessIds = a.Processes

	case *actions.SetProcess:
		if proc, ok := Processes.Processes[a.PID]; ok {
			// If process exists, only update relevant fields
			if a.Process.ProcessSummary != nil {
				proc.ProcessSummary = a.Process.ProcessSummary
			}
			if len(a.Process.Envelopes) > 0 {
				proc.Envelopes = a.Process.Envelopes
			}
			if a.Process.Process != nil {
				proc.Process = a.Process.Process
			}
			if a.Process.EnvelopeCount > proc.EnvelopeCount {
				proc.EnvelopeCount = a.Process.EnvelopeCount
			}
			Processes.Processes[a.PID] = proc
		} else { // Otherwise set whatever we have
			Processes.Processes[a.PID] = a.Process
		}

	case *actions.ProcessTabChange:
		Processes.Pagination.Tab = a.Tab

	case *actions.SetProcessCount:
		Processes.Count = a.Count

	case *actions.SetCurrentProcessStruct:
		Processes.CurrentProcess = a.Process

	case *actions.SetCurrentProcessEnvelopeCount:
		if Processes.CurrentProcess != nil {
			Processes.CurrentProcess.EnvelopeCount = a.Count
		}

	case *actions.SetCurrentProcessEnvelopes:
		if Processes.CurrentProcess != nil {
			Processes.CurrentProcess.Envelopes = a.EnvelopeList
		}

	case *actions.SetProcessResults:
		Processes.ProcessResults[a.PID] = a.Results

	default:
		return // don't fire listeners
	}

	Listeners.Fire()
}

// redirectActions is the handler for all redirect-related actions
func redirectActions(action interface{}) {
	switch a := action.(type) {
	case *actions.SignalRedirect:
		RedirectChan <- struct{}{}

	case *actions.SetCurrentPage:
		CurrentPage = a.Page

	default:
		return // don't fire listeners
	}

	Listeners.Fire()
}

// statsActions is the handler for all stats-related store actions
func statsActions(action interface{}) {
	switch a := action.(type) {

	case *actions.SetStats:
		Stats = a.Stats

	default:
		return // don't fire listeners
	}

	Listeners.Fire()
}

// transactionActions is the handler for all transaction-related store actions
func transactionActions(action interface{}) {
	switch a := action.(type) {
	case *actions.SetTransactionCount:
		Transactions.Count = a.Count

	case *actions.TransactionTabChange:
		Transactions.Pagination.Tab = a.Tab

	case *actions.TransactionsIndexChange:
		Transactions.Pagination.Index = a.Index

	case *actions.SetTransactionList:
		Transactions.Transactions = a.TransactionList

	case *actions.SetCurrentTransaction:
		Transactions.CurrentTransaction = a.Transaction

	case *actions.SetTransactionBlock:
		Transactions.CurrentBlock = a.Block

	case *actions.SetCurrentDecodedTransaction:
		Transactions.CurrentDecodedTransaction = a.Transaction

	default:
		return // don't fire listeners
	}

	Listeners.Fire()
}

// validatorActions is the handler for all validator-related store actions
func validatorActions(action interface{}) {
	switch a := action.(type) {
	case *actions.ValidatorsIndexChange:
		Validators.Pagination.Index = a.Index

	case *actions.SetValidatorList:
		Validators.Validators = a.List

	case *actions.SetValidatorCount:
		Validators.Count = a.Count

	case *actions.SetCurrentValidator:
		Validators.CurrentValidator = a.Validator

	case *actions.SetCurrentValidatorID:
		Validators.CurrentValidatorID = a.ID

	default:
		return // don't fire listeners
	}

	Listeners.Fire()
}

// clientActions is the handler for all connection-related store actions
func clientActions(action interface{}) {
	switch a := action.(type) {
	case *actions.GatewayConnected:
		if a.GatewayErr != nil {
			ServerConnected = false
			logger.Error(fmt.Errorf("error getting Gateway info: %s", a.GatewayErr))
			if GatewayError == "" {
				// Don't override first gatewayError if connection hasn't been fixed
				GatewayError = a.GatewayErr.Error()
			}
		} else {
			ServerConnected = true
			GatewayError = ""
		}

	case *actions.SetLoading:
		Loading = a.Loading

	case *actions.SetSearchTerm:
		SearchTerm = a.SearchTerm

	default:
		return // don't fire listeners
	}

	Listeners.Fire()
}

// blockActions is the handler for all block-related store actions
func blockActions(action interface{}) {
	switch a := action.(type) {
	case *actions.BlocksIndexChange:
		Blocks.Pagination.Index = a.Index

	case *actions.BlockTransactionsIndexChange:
		Blocks.TransactionPagination.Index = a.Index

	case *actions.BlocksTabChange:
		Blocks.Pagination.Tab = a.Tab

	case *actions.BlocksHeightUpdate:
		Blocks.Count = a.Height

	case *actions.SetBlockList:
		Blocks.Blocks = a.BlockList

	case *actions.SetCurrentBlock:
		Blocks.CurrentBlock = a.Block

	case *actions.SetCurrentBlockHeight:
		Blocks.CurrentBlockHeight = a.Height

	case *actions.SetCurrentBlockTransactionList:
		Blocks.CurrentTxs = a.TransactionList

	default:
		return // don't fire listeners
	}

	Listeners.Fire()
}

func disableUpdateActions(action interface{}) {
	switch a := action.(type) {
	case *actions.DisableUpdate:
		*a.Updater = a.Disabled

	case *actions.EnableAllUpdates:
		Blocks.Pagination.DisableUpdate = false
		Entities.Pagination.DisableUpdate = false
		Entities.ProcessPagination.DisableUpdate = false
		Envelopes.Pagination.DisableUpdate = false
		Transactions.Pagination.DisableUpdate = false
		Validators.Pagination.DisableUpdate = false
		Processes.Pagination.DisableUpdate = false
		Processes.EnvelopePagination.DisableUpdate = false
	}
}
