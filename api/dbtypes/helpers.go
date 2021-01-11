package dbtypes

//BlockIsEmpty returns true if block is empty
func BlockIsEmpty(s *StoreBlock) bool {
	if s == nil || len(s.Hash) == 0 && s.Height == 0 && s.NumTxs == 0 {
		return true
	}
	return false
}

//TxIsEmpty returns true if tx is empty
func TxIsEmpty(s *Transaction) bool {
	if s == nil || len(s.Hash) == 0 && s.TxHeight == 0 && s.Height == 0 && s.Index == 0 && s == nil {
		return true
	}
	return false
}

//EnvelopeIsEmpty returns true if env is empty
func EnvelopeIsEmpty(e *Envelope) bool {
	if e == nil || e.Nullifier == "" && len(e.Package) == 0 {
		return true
	}
	return false
}

//ValidatorIsEmpty returns true if validator is empty
func ValidatorIsEmpty(v *Validator) bool {
	if v == nil || len(v.Address) == 0 && v.Height.Height == 0 {
		return true
	}
	return false
}
