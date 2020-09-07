package store

import (
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
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
		Entities.ProcessesIndex = a.Index

	case *actions.EntityProcessesPageChange:
		Entities.ProcessesPage = a.Index

	case *actions.SetEntityIDs:
		Entities.EntityIDs = a.EntityIDs

	case *actions.SetCurrentEntityID:
		Entities.CurrentEntityID = a.EntityID

	case *actions.EntitiesTabChange:
		Entities.Pagination.Tab = a.Tab

	case *actions.SetEntityCount:
		Entities.Count = a.Count

	case *actions.SetProcessHeights:
		Entities.ProcessHeights = a.ProcessHeights

	case *actions.SetEntityProcessList:
		Entities.CurrentEntity.ProcessIDs = a.ProcessList

	case *actions.SetEntityProcessCount:
		Entities.CurrentEntity.ProcessCount = a.Count

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

	case *actions.SetCurrentEnvelopeHeight:
		Envelopes.CurrentEnvelopeHeight = a.Height

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
		Processes.EnvelopesIndex = a.Index

	case *actions.ProcessEnvelopesPageChange:
		Processes.EnvelopesPage = a.Index

	case *actions.SetProcessIDs:
		Processes.ProcessIDs = a.ProcessIDs

	case *actions.ProcessesTabChange:
		Processes.Pagination.Tab = a.Tab

	case *actions.SetProcessCount:
		Processes.Count = a.Count

	case *actions.SetEnvelopeHeights:
		Processes.EnvelopeHeights = a.EnvelopeHeights

	case *actions.SetProcessContents:
		Processes.ProcessResults[a.ID] = a.Process

	case *actions.SetProcessKeys:
		Processes.ProcessKeys[a.ID] = a.Keys

	case *actions.SetProcessState:
		Processes.CurrentProcess.State = a.State

	case *actions.SetProcessType:
		Processes.CurrentProcess.ProcessType = a.Type

	case *actions.SetCurrentProcessEnvelopeHeight:
		Processes.CurrentProcess.EnvelopeCount = a.Height

	case *actions.SetCurrentProcess:
		Processes.CurrentProcess = a.Process

	case *actions.SetCurrentProcessID:
		Processes.CurrentProcessID = a.ID

	case *actions.SetCurrentProcessEnvelopes:
		Processes.CurrentProcess.Envelopes = a.EnvelopeList

	default:
		return // don't fire listeners
	}

	Listeners.Fire()
}

// redirectActions is the handler for all redirect-related actions
func redirectActions(action interface{}) {
	switch action.(type) {
	case *actions.SignalRedirect:
		RedirectChan <- struct{}{}

	default:
		return // don't fire listeners
	}

	Listeners.Fire()
}

// statsActions is the handler for all stats-related store actions
func statsActions(action interface{}) {
	switch a := action.(type) {
	case *actions.SetResultStatus:
		Stats.ResultStatus = a.Status

	case *actions.SetGenesis:
		Stats.Genesis = a.Genesis

	case *actions.SetGatewayInfo:
		Stats.APIList = a.APIList
		Stats.Health = a.Health

	case *actions.SetBlockStatus:
		Stats.BlockTime = a.BlockTime
		Stats.BlockTimeStamp = a.BlockTimeStamp
		Stats.Height = a.Height

	default:
		return // don't fire listeners
	}

	Listeners.Fire()
}

// transactionActions is the handler for all transaction-related store actions
func transactionActions(action interface{}) {
	switch a := action.(type) {
	case *actions.TransactionTabChange:
		Transactions.Pagination.Tab = a.Tab

	case *actions.TransactionsIndexChange:
		Transactions.Pagination.Index = a.Index

	case *actions.SetTransactionList:
		Transactions.Transactions = a.TransactionList

	case *actions.SetTransactionCount:
		Transactions.Count = a.Count

	case *actions.SetCurrentTransactionHeight:
		Transactions.CurrentTransactionHeight = a.Height

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

	case *actions.SetCurrentValidatorBlockCount:
		Validators.CurrentBlockCount = a.Count

	case *actions.SetCurrentValidatorBlockList:
		Validators.CurrentBlockList = a.BlockList

	default:
		return // don't fire listeners
	}

	Listeners.Fire()
}

// clientActions is the handler for all connection-related store actions
func clientActions(action interface{}) {
	switch a := action.(type) {
	case *actions.TendermintClientInit:
		TendermintClient = a.Client

	case *actions.GatewayClientInit:
		GatewayClient = a.Client

	case *actions.GatewayConnected:
		GatewayConnected = a.Connected

	case *actions.ServerConnected:
		ServerConnected = a.Connected

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

	case *actions.SetCurrentBlockTxHeights:
		Blocks.CurrentBlockTxHeights = a.Heights

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
		Envelopes.Pagination.DisableUpdate = false
		Transactions.Pagination.DisableUpdate = false
		Validators.Pagination.DisableUpdate = false
		Processes.Pagination.DisableUpdate = false
	}
}
