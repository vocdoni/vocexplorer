package components

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/go-dvote/crypto/nacl"
	"gitlab.com/vocdoni/go-dvote/types"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/api/dbtypes"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/logger"
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
func (c *EnvelopeContents) Mount() {
	if !c.Rendered {
		c.Rendered = true
		vecty.Rerender(c)
	}
}

// Render renders the EnvelopeContents component
func (c *EnvelopeContents) Render() vecty.ComponentOrHTML {
	if !c.Rendered {
		return LoadingBar()
	}

	if store.Envelopes.CurrentEnvelope == nil || dbtypes.EnvelopeIsEmpty(store.Envelopes.CurrentEnvelope) {
		return elem.Div(
			vecty.Markup(vecty.Attribute("id", "main")),
			renderServerConnectionBanner(),
			elem.Main(vecty.Text("Envelope not available")),
		)
	}
	// Decode vote package
	var decryptionStatus string
	displayPackage := false
	var votePackage *dvotetypes.VotePackage
	pkeys := store.Processes.ProcessKeys[store.Envelopes.CurrentEnvelope.ProcessID]
	results := store.Processes.ProcessResults[store.Envelopes.CurrentEnvelope.ProcessID]
	keys := []string{}
	// If package is encrypted
	if !types.ProcessIsEncrypted[results.ProcessType] {
		decryptionStatus = "Vote unencrypted"
		displayPackage = true
	} else { // process is/was encrypted
		if pkeys != nil {
		indexLoop:
			for _, index := range store.Envelopes.CurrentEnvelope.EncryptionKeyIndexes {
				for _, key := range pkeys.Priv {
					if key.Idx == int(index) {
						keys = append(keys, key.Key)
						break
					} else {
						decryptionStatus = "Process is " + results.State + ", vote cannot be decrypted"
						displayPackage = false
						break indexLoop
					}
				}
				decryptionStatus = "Vote decrypted"
				displayPackage = true
			}
		} else {
			decryptionStatus = "Unable to decrypt: no keys available"
			displayPackage = false
		}
	}
	if len(keys) == len(store.Envelopes.CurrentEnvelope.EncryptionKeyIndexes) {
		var err error
		votePackage, err = unmarshalVote(store.Envelopes.CurrentEnvelope.Package, keys)
		if err != nil {
			logger.Error(err)
			decryptionStatus = "Unable to decode vote"
			displayPackage = false
		}
	}
	c.DecryptionStatus = decryptionStatus
	c.DisplayPackage = displayPackage
	c.VotePackage = votePackage

	return Container(
		vecty.Markup(vecty.Attribute("id", "main")),
		renderServerConnectionBanner(),
		DetailsView(
			c.EnvelopeView(),
			c.EnvelopeDetails(),
		))
}

// EnvelopeTab is the current active tab on the envelopes page
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
		pkeys, ok := api.GetProcessKeys(store.Envelopes.CurrentEnvelope.ProcessID)
		if ok {
			dispatcher.Dispatch(&actions.SetProcessKeys{Keys: pkeys, ID: store.Envelopes.CurrentEnvelope.ProcessID})
		}
	}
	// Ensure process keys are stored
	if _, ok := store.Processes.ProcessResults[store.Envelopes.CurrentEnvelope.ProcessID]; !ok {
		results, ok := api.GetProcessResults(strings.ToLower(store.Envelopes.CurrentEnvelope.ProcessID))
		if ok && results != nil {
			dispatcher.Dispatch(&actions.SetProcessContents{
				ID: store.Envelopes.CurrentEnvelope.ProcessID,
				Process: storeutil.Process{
					ProcessType: results.Type,
					State:       results.State,
					Results:     results.Results},
			})
		}
	}
}

// EnvelopeView renders one envelope
func (c *EnvelopeContents) EnvelopeView() vecty.List {
	return vecty.List{
		elem.Heading1(
			vecty.Markup(vecty.Class("card-title")),
			vecty.Text("Envelope details"),
		),
		elem.Heading2(
			vecty.Text(fmt.Sprintf(
				"Envelope height: %d",
				store.Envelopes.CurrentEnvelope.GlobalHeight,
			)),
		),
		elem.HorizontalRule(),
		elem.DescriptionList(
			elem.DefinitionTerm(vecty.Text("Belongs to process")),
			elem.Description(Link(
				"/process/"+util.TrimHex(store.Envelopes.CurrentEnvelope.ProcessID),
				util.TrimHex(store.Envelopes.CurrentEnvelope.ProcessID),
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
				humanize.Ordinal(int(store.Envelopes.CurrentEnvelope.ProcessHeight)),
			)),
			elem.DefinitionTerm(vecty.Text("Nullifier")),
			elem.Description(vecty.Text(
				store.Envelopes.CurrentEnvelope.Nullifier,
			)),
			elem.DefinitionTerm(vecty.Text("Vote type")),
			elem.Description(vecty.Text(
				util.GetEnvelopeName(store.Processes.ProcessResults[store.Envelopes.CurrentEnvelope.ProcessID].ProcessType),
			)),
			elem.DefinitionTerm(vecty.Text("Process status")),
			elem.Description(vecty.Text(strings.Title(store.Processes.ProcessResults[store.Envelopes.CurrentEnvelope.ProcessID].State))),
			elem.DefinitionTerm(vecty.Text("Decryption status")),
			elem.Description(vecty.Text(
				c.DecryptionStatus,
			)),
		),
	}
}

// EnvelopeDetails renders the details of an envelope contents
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
			vecty.Markup(vecty.Attribute("aria-label", "Tab navigation")),
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

func (c *EnvelopeContents) renderVotePackage() vecty.ComponentOrHTML {
	if c.DisplayPackage {
		voteString, err := json.MarshalIndent(c.VotePackage, "", "\t")
		if err != nil {
			logger.Error(err)
		}
		if len(voteString) > 0 {
			return elem.Preformatted(vecty.Text(string(voteString)))
		}
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
				logger.Warn("cannot create private key cipher: " + err.Error())
				continue
			}
			if rawVote, err = priv.Decrypt(rawVote); err != nil {
				logger.Warn("cannot decrypt vote with key " + util.IntToString(i))
			}
		}
	}
	if err := json.Unmarshal(rawVote, &vote); err != nil {
		return nil, err
	}
	return &vote, nil
}
