package tmtypes

import (
	"reflect"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

var cdc = amino.NewCodec()

func init() {
	RegisterBlockAmino(cdc)
}

func RegisterBlockAmino(cdc *amino.Codec) {
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PubKeyEd25519{},
		ed25519.PubKeyAminoName, nil)
	// cryptoamino.RegisterAmino(cdc)
	// RegisterEvidences(cdc)
}

// func RegisterEvidences(cdc *amino.Codec) {
// 	cdc.RegisterInterface((*Evidence)(nil), nil)
// 	cdc.RegisterConcrete(&DuplicateVoteEvidence{}, "tendermint/DuplicateVoteEvidence", nil)
// }

// Go lacks a simple and safe way to see if something is a typed nil.
// See:
//  - https://dave.cheney.net/2017/08/09/typed-nils-in-go-2
//  - https://groups.google.com/forum/#!topic/golang-nuts/wnH302gBa4I/discussion
//  - https://github.com/golang/go/issues/21538
func isTypedNil(o interface{}) bool {
	rv := reflect.ValueOf(o)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

// Returns true if it has zero length.
func isEmpty(o interface{}) bool {
	rv := reflect.ValueOf(o)
	switch rv.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return rv.Len() == 0
	default:
		return false
	}
}

// cdcEncode returns nil if the input is nil, otherwise returns
// cdc.MustMarshalBinaryBare(item)
func cdcEncode(item interface{}) []byte {
	if item != nil && !isTypedNil(item) && !isEmpty(item) {
		return cdc.MustMarshalBinaryBare(item)
	}
	return nil
}
