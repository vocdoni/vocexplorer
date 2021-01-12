package config

import "go.vocdoni.io/dvote/types"

// Maps to turn raw data into human readable

var (
	// TransactionTypeMap maps transaction types to readable descriptions
	TransactionTypeMap = map[string]string{
		types.TxVote:                 "Vote",
		types.TxNewProcess:           "Create new process",
		types.TxCancelProcess:        "Cancel process",
		types.TxAddValidator:         "Add validator",
		types.TxRemoveValidator:      "Remove validator",
		types.TxAddOracle:            "Add oracle",
		types.TxRemoveOracle:         "Remove oracle",
		types.TxAddProcessKeys:       "Add process keys",
		types.TxRevealProcessKeys:    "Reveal process keys",
		"setProcess":                 "Set process metadata",
		"admin":                      "Admin",
		"unknown":                    "Unknown",
		"TX_UNKNOWN":                 "Unknown",
		"NEW_PROCESS":                "Create new process",
		"CANCEL_PROCESS":             "Cancel process",
		"SET_PROCESS_STATUS":         "Set process status",
		"SET_PROCESS_CENSUS":         "Set process census",
		"SET_PROCESS_QUESTION_INDEX": "Set process question index",
		"ADD_PROCESS_KEYS":           "Add process keys",
		"REVEAL_PROCESS_KEYS":        "Reveal process keys",
		"ADD_ORACLE":                 "Add oracle",
		"REMOVE_ORACLE":              "Remove oracle",
		"ADD_VALIDATOR":              "Add validator",
		"REMOVE_VALIDATOR":           "Remove validator",
		"VOTE":                       "Vote",
		"SET_PROCESS_RESULTS":        "Set process results",
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
