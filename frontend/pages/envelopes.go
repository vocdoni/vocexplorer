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
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/dbapi"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
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
	envelope := dbapi.GetEnvelope(height)
	if envelope == nil || types.EnvelopeIsEmpty(envelope) {
		log.Errorf("Envelope unavailable")
		return elem.Div(
			elem.Main(vecty.Text("Envelope not available")),
		)
	}
	// Get process encryption keys:
	gwClient, cancel := client.InitGateway(home.Cfg.GatewayHost)
	defer cancel()
	if gwClient == nil {
		log.Error("Unable to connect to gateway client")
	}
	pkeys, err := gwClient.GetProcessKeys(envelope.GetProcessID())
	util.ErrPrint(err)

	// Decode vote package
	// TODO: decrypt vote package if necessary
	var decryptionStatus string
	var displayPackage bool
	var votePackage *dvotetypes.VotePackage
	keys := []string{}
	// If package is encrypted
	if len(envelope.EncryptionKeyIndexes) == 0 {
		decryptionStatus = "Vote unencrypted"
		displayPackage = true
	} else {
		decryptionStatus = "Vote decrypted"
		displayPackage = true
		for _, index := range envelope.EncryptionKeyIndexes {
			if len(pkeys.Priv) <= int(index) {
				decryptionStatus = "Process is still active, vote cannot be decrypted"
				displayPackage = false
				break
			}
			keys = append(keys, pkeys.Priv[index].Key)
		}
	}
	if len(keys) == len(envelope.EncryptionKeyIndexes) {
		votePackage, err = unmarshalVote(envelope.GetPackage(), keys)
		if util.ErrPrint(err) {
			decryptionStatus = "Unable to decode vote"
			displayPackage = false
		}
	}
	util.ErrPrint(err)
	return elem.Div(
		&components.EnvelopeContents{
			Envelope:         envelope,
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
