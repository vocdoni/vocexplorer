package config

import "go.vocdoni.io/dvote/types"

// Maps to turn raw data into human readable

var (
	// TransactionTypeMap maps transaction types to readable descriptions
	TransactionTypeMap = map[string]string{
		types.TxVote:              "Vote",
		types.TxNewProcess:        "Create new process",
		types.TxCancelProcess:     "Cancel process",
		types.TxAddValidator:      "Add validator",
		types.TxRemoveValidator:   "Remove validator",
		types.TxAddOracle:         "Add oracle",
		types.TxRemoveOracle:      "Remove oracle",
		types.TxAddProcessKeys:    "Add process keys",
		types.TxRevealProcessKeys: "Reveal process keys",
		"setProcess":              "Set process metadata",
		"admin":                   "Admin",
		"unknown":                 "Unknown",
	}

	// ProcessTypeMap maps process types to readable descriptions
	ProcessTypeMap = map[string]string{
		types.PetitionSign:  "Petition",
		types.PollVote:      "Poll",
		types.EncryptedPoll: "Encrypted poll",
		types.SnarkVote:     "Anonymous poll",
	}

	// EnvelopeTypeMap maps envelope types to readable descriptions
	EnvelopeTypeMap = map[string]string{
		types.PetitionSign:  "Petition signature",
		types.PollVote:      "Poll vote",
		types.EncryptedPoll: "Encrypted poll vote",
		types.SnarkVote:     "Anonymous poll vote",
	}
)
