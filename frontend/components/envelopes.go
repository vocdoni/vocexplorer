package components

import (
	"encoding/base64"
	"encoding/json"
	"strconv"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/go-dvote/log"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/dbapi"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
	router "marwan.io/vecty-router"
)

// EnvelopesView renders the Envelopes page
type EnvelopesView struct {
	vecty.Core
	cfg *config.Cfg
}

// Render renders the EnvelopesView component
func (home *EnvelopesView) Render() vecty.ComponentOrHTML {
	height, err := strconv.ParseInt(router.GetNamedVar(home)["id"], 0, 64)
	util.ErrPrint(err)
	envelope := dbapi.GetEnvelope(height)
	if envelope == nil || types.EnvelopeIsEmpty(envelope) {
		log.Errorf("Envelope unavailable")
		return elem.Div(
			&Header{},
			elem.Main(vecty.Text("Envelope not available")),
		)
	}
	// Decode vote package
	// TODO: decrypt vote package if necessary
	packageBytes, err := base64.StdEncoding.DecodeString(envelope.GetPackage())
	util.ErrPrint(err)
	votePackage := new(dvotetypes.VotePackageStruct)
	util.ErrPrint(json.Unmarshal(packageBytes, votePackage))
	return elem.Div(
		&Header{},
		&EnvelopeContents{
			Envelope:    envelope,
			VotePackage: votePackage,
		},
	)
}
