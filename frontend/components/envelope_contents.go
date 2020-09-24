package components

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/go-dvote/crypto/nacl"
	"gitlab.com/vocdoni/go-dvote/log"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/proto"
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

	if store.Envelopes.CurrentEnvelope == nil || proto.EnvelopeIsEmpty(store.Envelopes.CurrentEnvelope) {
		return elem.Div(
			renderGatewayConnectionBanner(),
			renderServerConnectionBanner(),
			elem.Main(vecty.Text("Envelope not available")),
		)
	}
	// Decode vote package
	var decryptionStatus string
	var displayPackage bool
	var votePackage *dvotetypes.VotePackage
	pkeys := store.Processes.ProcessKeys[store.Envelopes.CurrentEnvelope.ProcessID]
	keys := []string{}
	// If package is encrypted
	if len(store.Envelopes.CurrentEnvelope.EncryptionKeyIndexes) == 0 {
		decryptionStatus = "Vote unencrypted"
		displayPackage = true
	} else {
		if pkeys != nil {
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
		} else {
			decryptionStatus = "Unable to decrypt"
			displayPackage = true
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

	return Container(
		renderGatewayConnectionBanner(),
		renderServerConnectionBanner(),
		DetailsView(
			contents.EnvelopeView(),
			contents.EnvelopeDetails(),
		))
}

type EnvelopeTab struct {
	*Tab
}

func (t *EnvelopeTab) store() string {
	return store.Envelopes.Pagination.Tab
}
func (t *EnvelopeTab) dispatch() interface{} {
	return &actions.EnvelopesTabChange{
		Tab: t.alias(),
	}
}

// UpdateEnvelopeContents keeps the envelope contents up to date
func UpdateEnvelopeContents(d *EnvelopeContents) {
	dispatcher.Dispatch(&actions.EnableAllUpdates{})
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

func (c *EnvelopeContents) EnvelopeView() vecty.List {
	return vecty.List{
		elem.Heading1(
			vecty.Markup(vecty.Class("card-title")),
			vecty.Text("Envelope details"),
		),
		elem.Heading2(
			vecty.Text(fmt.Sprintf(
				"Envelope height: %d",
				store.Envelopes.CurrentEnvelope.GetGlobalHeight(),
			)),
		),
		elem.HorizontalRule(),
		elem.DescriptionList(
			elem.DefinitionTerm(vecty.Text("Belongs to process")),
			elem.Description(Link(
				"/process/"+util.TrimHex(store.Envelopes.CurrentEnvelope.GetProcessID()),
				util.TrimHex(store.Envelopes.CurrentEnvelope.GetProcessID()),
				"hash",
			)),
			elem.DefinitionTerm(vecty.Text("Packaged in transaction")),
			elem.Description(Link(
				"/transaction/"+util.IntToString(store.Envelopes.CurrentEnvelope.TxHeight),
				util.IntToString(store.Envelopes.CurrentEnvelope.TxHeight),
				"hash",
			)),
			elem.DefinitionTerm(vecty.Text("Position in process")),
			elem.Description(vecty.Text(
				humanize.Ordinal(int(store.Envelopes.CurrentEnvelope.GetProcessHeight())),
			)),
			elem.DefinitionTerm(vecty.Text("Nullifier")),
			elem.Description(vecty.Text(
				store.Envelopes.CurrentEnvelope.GetNullifier(),
			)),
			elem.DefinitionTerm(vecty.Text("Decryption status")),
			elem.Description(vecty.Text(
				c.DecryptionStatus,
			)),
		),
	}
}

func (c *EnvelopeContents) EnvelopeDetails() vecty.ComponentOrHTML {
	cTab := &EnvelopeTab{&Tab{
		Text:  "Contents",
		Alias: "contents",
	}}

	contents := c.renderVotePackage()

	if contents == nil {
		return nil
	}

	return vecty.List{
		elem.Navigation(
			vecty.Markup(vecty.Class("tabs")),
			elem.UnorderedList(
				TabLink(c, cTab),
			),
		),
		elem.Div(
			vecty.Markup(vecty.Class("tabs-content")),
			TabContents(cTab, contents),
		),
	}
}

func (contents *EnvelopeContents) renderVotePackage() vecty.ComponentOrHTML {
	if contents.DisplayPackage {
		voteString, err := json.MarshalIndent(contents.VotePackage, "", "\t")
		if err != nil {
			log.Error(err)
		}
		return elem.Preformatted(vecty.Text(string(voteString)))
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
