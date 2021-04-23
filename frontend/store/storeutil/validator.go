package storeutil

import (
	"go.vocdoni.io/proto/build/go/models"
)

// Validators stores all data about blockchain validators
type Validators struct {
	Count              int
	CurrentValidator   *models.Validator
	CurrentValidatorID string
	Pagination         PageStore
	Validators         []*models.Validator
}
