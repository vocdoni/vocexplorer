package components

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"

	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
	"gitlab.com/vocdoni/vocexplorer/logger"
	"gitlab.com/vocdoni/vocexplorer/util"
	"go.vocdoni.io/dvote/crypto/nacl"
	indexertypes "go.vocdoni.io/dvote/vochain/scrutinizer/indexertypes"
	"go.vocdoni.io/proto/build/go/models"
)

// EnvelopeContents renders envelope contents
type EnvelopeContents struct {
	vecty.Core
	vecty.Mounter
	DecryptionStatus string
	DisplayPackage   bool
	VotePackage      *indexertypes.VotePackage
	Rendered         bool
	Unavailable      bool
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
	if c.Unavailable {
		return Unavailable("Envelope unavailable", "")
	}
	if store.Envelopes.CurrentEnvelope == nil {
		return Unavailable("Loading envelope...", "")
	}

	// Decode vote package
	var decryptionStatus string
	displayPackage := false
	var votePackage *indexertypes.VotePackage
	process, pok := store.Processes.Processes[util.HexToString(store.Envelopes.CurrentEnvelope.Meta.ProcessId)]
	results, rok := store.Processes.ProcessResults[util.HexToString(store.Envelopes.CurrentEnvelope.Meta.ProcessId)]
	if process == nil || process.Process == nil || !pok || !rok {
		return Unavailable("Loading process...", "")
	}
	pkeys := process.Process.PrivateKeys
	keys := []string{}
	// If package is encrypted
	if !strings.Contains(strings.ToLower(results.Type), "encrypted") {
		decryptionStatus = "Vote unencrypted"
		displayPackage = true
	} else { // process is/was encrypted
		// If not ended or results, keys must be available
		if s := strings.ToLower(results.State); s != "ended" && s != "results" {
			decryptionStatus = fmt.Sprintf("Vote cannot be decrypted yet: process in state %s", s)
			displayPackage = false
		} else if pkeys != nil {
			// If ended or results then check for the keys
			if len(store.Envelopes.CurrentEnvelope.EncryptionKeyIndexes) > 0 {
				for _, key := range pkeys {
					keys = append(keys, key)
				}
				if len(store.Envelopes.CurrentEnvelope.EncryptionKeyIndexes) != len(keys) {
					decryptionStatus = fmt.Sprintf("Vote cannot be decrypted yet: %d keys expected, %d provided",
						len(store.Envelopes.CurrentEnvelope.EncryptionKeyIndexes), len(keys))
					displayPackage = false
				} else {
					decryptionStatus = "Vote decrypted"
					displayPackage = true
				}
			}
		} else {
			decryptionStatus = "Unable to decrypt: no keys available"
			displayPackage = false
		}
	}
	if len(keys) == len(store.Envelopes.CurrentEnvelope.EncryptionKeyIndexes) {
		var err error
		votePackage, err = unmarshalVote(store.Envelopes.CurrentEnvelope.VotePackage, keys)
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
	// Set current envelope to nil so previous one is not displayed
	dispatcher.Dispatch(&actions.SetCurrentEnvelope{Envelope: nil})
	dispatcher.Dispatch(&actions.EnableAllUpdates{})
	d.fetchEnvelope()
	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	if !update.CheckCurrentPage("envelope", ticker) {
		return
	}
	for {
		select {
		case <-store.RedirectChan:
			if !update.CheckCurrentPage("envelope", ticker) {
				return
			}
		case <-ticker.C:
			if !update.CheckCurrentPage("envelope", ticker) {
				return
			}
			// If envelope never loaded, load it
			if d.Unavailable {
				d.fetchEnvelope()
			}
		}
	}
}

func (c *EnvelopeContents) fetchEnvelope() {
	// Fetch actual envelope contents
	envelope, err := store.Client.GetEnvelope(store.Envelopes.CurrentEnvelopeNullifier)
	if err != nil {
		c.Unavailable = true
		dispatcher.Dispatch(&actions.SetCurrentEnvelope{Envelope: nil})
		logger.Error(err)
	} else {
		c.Unavailable = false
		dispatcher.Dispatch(&actions.SetCurrentEnvelope{Envelope: envelope})
	}
	// Ensure process & results are stored
	if store.Envelopes.CurrentEnvelope == nil {
		return
	}
	process, err := store.Client.GetProcess(store.Envelopes.CurrentEnvelope.Meta.ProcessId)
	if err != nil {
		logger.Error(err)
	}
	pubKeys, privKeys, err := store.Client.GetProcessKeys(store.Envelopes.CurrentEnvelope.Meta.ProcessId)
	if err != nil {
		logger.Error(err)
	}
	for _, key := range pubKeys {
		process.PublicKeys = append(process.PublicKeys, key.Key)
	}
	for _, key := range privKeys {
		process.PrivateKeys = append(process.PrivateKeys, key.Key)
	}
	if process != nil {
		dispatcher.Dispatch(&actions.SetProcess{
			PID: util.HexToString(store.Envelopes.CurrentEnvelope.Meta.ProcessId),
			Process: &storeutil.Process{
				Process: process,
			},
		})
	}
	if _, ok := store.Processes.ProcessResults[util.HexToString(store.Envelopes.CurrentEnvelope.Meta.ProcessId)]; !ok {
		results, state, tp, final, err := store.Client.GetResults(store.Envelopes.CurrentEnvelope.Meta.ProcessId)
		if err != nil {
			logger.Error(err)
		}
		dispatcher.Dispatch(&actions.SetProcessResults{
			Results: storeutil.ProcessResults{
				Results: results,
				State:   state,
				Type:    tp,
				Final:   final,
			},
			PID: util.HexToString(store.Envelopes.CurrentEnvelope.Meta.ProcessId),
		})
	}
}

// EnvelopeView renders one envelope
func (c *EnvelopeContents) EnvelopeView() vecty.List {
	return vecty.List{
		elem.Heading1(
			vecty.Markup(vecty.Class("card-title")),
			vecty.Text("Envelope details"),
		),
		// elem.Heading2(
		// 	vecty.Text(fmt.Sprintf(
		// 		"Envelope height: %d",
		// 		store.Envelopes.CurrentEnvelope.H,
		// 	)),
		// ),
		elem.HorizontalRule(),
		elem.DescriptionList(
			elem.DefinitionTerm(vecty.Text("Belongs to process")),
			elem.Description(Link(
				"/process/"+util.TrimHex(util.HexToString(store.Envelopes.CurrentEnvelope.Meta.ProcessId)),
				util.HexToString(store.Envelopes.CurrentEnvelope.Meta.ProcessId),
				"hash",
			)),
			elem.DefinitionTerm(vecty.Text("Packaged in transaction")),
			elem.Description(Link(
				"/transaction/"+util.IntToString(store.Envelopes.CurrentEnvelope.Meta.Height)+"/"+util.IntToString(store.Envelopes.CurrentEnvelope.Meta.TxIndex),
				util.IntToString(store.Envelopes.CurrentEnvelope.Meta.TxIndex+1)+" on block "+util.IntToString(store.Envelopes.CurrentEnvelope.Meta.Height),
				"hash",
			)),
			// elem.DefinitionTerm(vecty.Text("Position in process")),
			// elem.Description(vecty.Text(
			// 	humanize.Ordinal(int(store.Envelopes.CurrentEnvelope.ProcessHeight)),
			// )),
			elem.DefinitionTerm(vecty.Text("Nullifier")),
			elem.Description(vecty.Text(
				util.HexToString(store.Envelopes.CurrentEnvelope.Meta.Nullifier),
			)),
			elem.DefinitionTerm(vecty.Text("Vote type")),
			elem.Description(vecty.Text(
				util.GetEnvelopeName(store.Processes.ProcessResults[util.HexToString(store.Envelopes.CurrentEnvelope.Meta.ProcessId)].Type),
			)),
			elem.DefinitionTerm(vecty.Text("Encryption key indexes")),
			elem.Description(vecty.Text(
				fmt.Sprintf("%v", store.Envelopes.CurrentEnvelope.EncryptionKeyIndexes),
			)),

			vecty.If(store.Envelopes.CurrentEnvelope.Weight != "", elem.DefinitionTerm(vecty.Text("Envelope weight"))),
			vecty.If(store.Envelopes.CurrentEnvelope.Weight != "", elem.Description(vecty.Text(store.Envelopes.CurrentEnvelope.Weight))),
			elem.DefinitionTerm(vecty.Text("Process status")),
			elem.Description(vecty.Text(strings.Title(store.Processes.ProcessResults[util.HexToString(store.Envelopes.CurrentEnvelope.Meta.ProcessId)].State))),
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
	details := &EnvelopeTab{&Tab{
		Text:  "Details",
		Alias: "details",
	}}

	if store.Envelopes.CurrentEnvelope == nil {
		return nil
	}
	contents := c.renderVotePackage()
	if contents == nil {
		return nil
	}

	process := store.Processes.Processes[util.HexToString(store.Envelopes.CurrentEnvelope.Meta.ProcessId)]
	if process == nil {
		return nil
	}

	envelopeDetails := elem.Div(vecty.Markup(vecty.Class("poll-details")), renderEnvelopeType(process.Process.Envelope))

	return vecty.List{
		elem.Navigation(
			vecty.Markup(vecty.Attribute("aria-label", "Tab navigation")),
			vecty.Markup(vecty.Class("tabs")),
			elem.UnorderedList(
				TabLink(c, cTab),
				TabLink(c, details),
			),
		),
		elem.Div(
			vecty.Markup(vecty.Class("tabs-content")),
			TabContents(cTab, contents),
			TabContents(details, envelopeDetails),
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
func unmarshalVote(votePackage []byte, keys []string) (*indexertypes.VotePackage, error) {
	var vote indexertypes.VotePackage
	rawVote := make([]byte, len(votePackage))
	copy(rawVote, votePackage)
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
		return nil, fmt.Errorf("cannot unmarshal vote: %w", err)
	}
	return &vote, nil
}

func renderEnvelopeType(envelopeType *models.EnvelopeType) vecty.ComponentOrHTML {
	if envelopeType == nil {
		return vecty.Text("Envelope Type unavailable")
	}
	return elem.Div(
		elem.Span(
			vecty.Markup(vecty.Class("detail")),
			vecty.Text("Envelope Type"),
		),
		elem.OrderedList(
			elem.ListItem(vecty.Text(fmt.Sprintf("Serial: %t", envelopeType.Serial))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Anonymous: %t", envelopeType.Anonymous))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Encrypted votes: %t", envelopeType.EncryptedVotes))),
			elem.ListItem(vecty.Text(fmt.Sprintf("Unique values: %t", envelopeType.UniqueValues))),
		))
}
