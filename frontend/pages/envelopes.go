package pages

import (
	"encoding/base64"
	"encoding/json"
	"strconv"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/go-dvote/crypto/nacl"
	"gitlab.com/vocdoni/go-dvote/log"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
	router "marwan.io/vecty-router"
)

// EnvelopesView renders the Envelopes page
type EnvelopesView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the EnvelopesView component
func (home *EnvelopesView) Render() vecty.ComponentOrHTML {
	height, err := strconv.ParseInt(router.GetNamedVar(home)["id"], 0, 64)
	util.ErrPrint(err)
	envelope, ok := api.GetEnvelope(height)
	if envelope == nil || types.EnvelopeIsEmpty(envelope) || !ok {
		log.Errorf("Envelope unavailable")
		return elem.Div(
			elem.Main(vecty.Text("Envelope not available")),
		)
	}
	dispatcher.Dispatch(&actions.SetCurrentEnvelope{Envelope: envelope})
	var pkeys *api.Pkeys
	if pkeys, ok := store.Processes.ProcessKeys[store.Envelopes.CurrentEnvelope.ProcessID]; !ok {
		pkeys, err = store.GatewayClient.GetProcessKeys(store.Envelopes.CurrentEnvelope.GetProcessID())
		if err != nil {
			log.Error(err)
		} else {
			dispatcher.Dispatch(&actions.SetProcessKeys{Keys: pkeys, ID: store.Envelopes.CurrentEnvelope.ProcessID})
		}
	}

	// Decode vote package
	var decryptionStatus string
	var displayPackage bool
	var votePackage *dvotetypes.VotePackage
	keys := []string{}
	// If package is encrypted
	if len(store.Envelopes.CurrentEnvelope.EncryptionKeyIndexes) == 0 {
		decryptionStatus = "Vote unencrypted"
		displayPackage = true
	} else {
		decryptionStatus = "Vote decrypted"
		displayPackage = true
		for _, index := range store.Envelopes.CurrentEnvelope.EncryptionKeyIndexes {
			if len(pkeys.Priv) <= int(index) {
				decryptionStatus = "Process is still active, vote cannot be decrypted"
				displayPackage = false
				break
			}
			keys = append(keys, pkeys.Priv[index].Key)
		}
	}
	if len(keys) == len(store.Envelopes.CurrentEnvelope.EncryptionKeyIndexes) {
		votePackage, err = unmarshalVote(store.Envelopes.CurrentEnvelope.GetPackage(), keys)
		if util.ErrPrint(err) {
			decryptionStatus = "Unable to decode vote"
			displayPackage = false
		}
	}
	util.ErrPrint(err)
	return elem.Div(
		&components.EnvelopeContents{
			DecryptionStatus: decryptionStatus,
			DisplayPackage:   displayPackage,
			VotePackage:      votePackage,
		},
	)
}

// From go-dvote keykeepercli.go
func unmarshalVote(votePackage string, keys []string) (*dvotetypes.VotePackage, error) {
	rawVote, err := base64.StdEncoding.DecodeString(votePackage)
	if err != nil {
		return nil, err
	}
	var vote dvotetypes.VotePackage
	// if encryption keys, decrypt the vote
	if len(keys) > 0 {
		for i := len(keys) - 1; i >= 0; i-- {
			priv, err := nacl.DecodePrivate(keys[i])
			if err != nil {
				log.Warnf("cannot create private key cipher: (%s)", err)
				continue
			}
			if rawVote, err = priv.Decrypt(rawVote); err != nil {
				log.Warnf("cannot decrypt vote with index key %d", i)
			}
		}
	}
	if err := json.Unmarshal(rawVote, &vote); err != nil {
		return nil, err
	}
	return &vote, nil
}
