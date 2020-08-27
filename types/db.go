package types

//BlockIsEmpty returns true if block is empty
func BlockIsEmpty(s *StoreBlock) bool {
	if len(s.GetHash()) == 0 && s.GetHeight() == 0 && s.GetNumTxs() == 0 {
		return true
	}
	return false
}

//TxIsEmpty returns true if tx is empty
func TxIsEmpty(s *SendTx) bool {
	if len(s.GetHash()) == 0 && s.GetStore().GetTxHeight() == 0 && s.GetStore().GetHeight() == 0 && s.GetStore().GetIndex() == 0 {
		return true
	}
	return false
}

//EnvelopeIsEmpty returns true if env is empty
func EnvelopeIsEmpty(e *Envelope) bool {
	if e.GetNullifier() == "" && e.GetPackage() == "" {
		return true
	}
	return false
}

//ValidatorIsEmpty returns true if validator is empty
func ValidatorIsEmpty(v *Validator) bool {
	if len(v.GetAddress()) == 0 && v.GetHeight().GetHeight() == 0 {
		return true
	}
	return false
}
