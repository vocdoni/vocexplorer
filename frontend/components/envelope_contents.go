package components

import (
	"encoding/base64"
	"encoding/json"

	humanize "github.com/dustin/go-humanize"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
	"gitlab.com/vocdoni/go-dvote/crypto/nacl"
	"gitlab.com/vocdoni/go-dvote/log"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// EnvelopeContents renders envelope contents
type EnvelopeContents struct {
	vecty.Core
	vecty.Mounter
	DecryptionStatus string
	DisplayPackage   bool
	VotePackage      *dvotetypes.VotePackage
	Rendered         bool
}

// Mount triggers EnvelopeContents  renders
func (contents *EnvelopeContents) Mount() {
	if !contents.Rendered {
		contents.Rendered = true
		vecty.Rerender(contents)
	}
}

// Render renders the EnvelopeContents component
func (contents *EnvelopeContents) Render() vecty.ComponentOrHTML {
	if !contents.Rendered {
		return LoadingBar()
	}
	if store.Envelopes.CurrentEnvelope == nil || types.EnvelopeIsEmpty(store.Envelopes.CurrentEnvelope) {
		return elem.Div(
			elem.Main(vecty.Text("Envelope not available")),
		)
	}
	// Decode vote package
	var decryptionStatus string
	var displayPackage bool
	var votePackage *dvotetypes.VotePackage
	pkeys, _ := store.Processes.ProcessKeys[store.Envelopes.CurrentEnvelope.ProcessID]
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
		var err error
		votePackage, err = unmarshalVote(store.Envelopes.CurrentEnvelope.GetPackage(), keys)
		if err != nil {
			log.Error(err)
			decryptionStatus = "Unable to decode vote"
			displayPackage = false
		}
	}
	contents.DecryptionStatus = decryptionStatus
	contents.DisplayPackage = displayPackage
	contents.VotePackage = votePackage

	return elem.Main(
		renderEnvelopeHeader(),
		vecty.Text(contents.DecryptionStatus),
		contents.renderVotePackage(),
	)
}

// UpdateAndRenderEnvelopesDashboard keeps the envelope contents up to date
func UpdateAndRenderEnvelopesDashboard(d *EnvelopeContents) {
	actions.EnableUpdates()
	// Fetch actual envelope contents
	envelope, ok := api.GetEnvelope(store.Envelopes.CurrentEnvelopeHeight)
	if ok {
		dispatcher.Dispatch(&actions.SetCurrentEnvelope{Envelope: envelope})
	}
	// Ensure process keys are stored
	if _, ok := store.Processes.ProcessKeys[store.Envelopes.CurrentEnvelope.ProcessID]; !ok {
		pkeys, err := store.GatewayClient.GetProcessKeys(store.Envelopes.CurrentEnvelope.GetProcessID())
		if err != nil {
			log.Error(err)
		} else {
			dispatcher.Dispatch(&actions.SetProcessKeys{Keys: pkeys, ID: store.Envelopes.CurrentEnvelope.ProcessID})
		}
	}
}

func renderEnvelopeHeader() vecty.ComponentOrHTML {
	return elem.Div(vecty.Markup(vecty.Class("card-deck-col")),
		elem.Div(vecty.Markup(vecty.Class("card")),
			elem.Div(
				vecty.Markup(vecty.Class("card-header")),
				Link(
					"/envelope/"+util.IntToString(store.Envelopes.CurrentEnvelope.GetGlobalHeight()),
					util.IntToString(store.Envelopes.CurrentEnvelope.GetGlobalHeight()),
					"nav-link",
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("card-body")),
				elem.Div(
					vecty.Markup(vecty.Class("block-card-heading")),
					elem.Div(
						vecty.Text(humanize.Ordinal(int(store.Envelopes.CurrentEnvelope.GetProcessHeight()))+" envelope on process "),
						Link(
							"/process/"+util.StripHexString(store.Envelopes.CurrentEnvelope.GetProcessID()),
							util.StripHexString(store.Envelopes.CurrentEnvelope.GetProcessID()),
							"hash",
						),
					),
					elem.Div(
						elem.Div(
							vecty.Markup(vecty.Class("dt")),
							vecty.Text("Nullifier"),
						),
						elem.Div(
							vecty.Markup(vecty.Class("dd")),
							vecty.Text(store.Envelopes.CurrentEnvelope.GetNullifier()),
						),
					),
				),
			),
		),
	)
}

func (contents *EnvelopeContents) renderVotePackage() vecty.ComponentOrHTML {
	if contents.DisplayPackage {
		voteString, err := json.MarshalIndent(contents.VotePackage, "", "\t")
		if err != nil {
			log.Error(err)
		}
		accordionName := "accordionEnv"
		return elem.Div(
			vecty.Markup(vecty.Class("accordion"), prop.ID(accordionName)),
			renderCollapsible("Envelope Contents", accordionName, "One", elem.Preformatted(vecty.Text(string(voteString)))),
			// renderCollapsible("Data", accordionName, "Two", transactions),
			// renderCollapsible("Evidence", accordionName, "Three", elem.Preformatted(vecty.Text(string(evidence)))),
			// renderCollapsible("Last Commit", accordionName, "Four", elem.Preformatted(vecty.Text(string(commit)))),
		)
	}
	return nil
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
