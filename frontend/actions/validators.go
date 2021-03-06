package actions

import (
	"go.vocdoni.io/proto/build/go/models"
)

// ValidatorsIndexChange is the action to set the pagination index
type ValidatorsIndexChange struct {
	Index int
}

// SetValidatorList is the action to set the validator list
type SetValidatorList struct {
	List []*models.Validator
}

// SetValidatorCount is the action to set the Validator count
type SetValidatorCount struct {
	Count int
}

// SetCurrentValidator is the action to set the currently displayed validator
type SetCurrentValidator struct {
	Validator *models.Validator
}

// SetCurrentValidatorID is the action to set the currently displayed validator ID
type SetCurrentValidatorID struct {
	ID string
}
