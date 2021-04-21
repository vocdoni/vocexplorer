package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/util"
	router "marwan.io/vecty-router"
)

// EnvelopeNullifier renders envelope contents
type EnvelopeNullifier struct {
	vecty.Core
	vecty.Mounter
	Nullifier   string
	Rendered    bool
	Unavailable bool
}

// Mount triggers EnvelopeNullifier  renders
func (c *EnvelopeNullifier) Mount() {
	if !c.Rendered {
		c.Rendered = true
		vecty.Rerender(c)
	}
}

// Render renders the EnvelopeNullifier component
func (c *EnvelopeNullifier) Render() vecty.ComponentOrHTML {
	if !c.Rendered {
		return LoadingBar()
	}
	if c.Unavailable {
		return elem.Div(Unavailable("This envelope does not exist", "It can take up to 10 minutes for your vote to register on the blockchain once it has been cast."))
	}

	// If envelope is not "unavailable" and page has not redirected, envelope must still be loading
	return Unavailable("Loading envelope...", "")
}

// LoadEnvelopeHeight tries to load the height of the envelope and redirect to its envelope page
func (c *EnvelopeNullifier) LoadEnvelopeHeight() {
	var envelopeHeight int64
	ok := false
	if c.Nullifier != "" {
		envelopeHeight, ok = api.GetEnvelopeHeightFromNullifier(c.Nullifier)
	}
	if !ok || envelopeHeight == 0 {
		c.Unavailable = true
		vecty.Rerender(c)
	} else {
		router.Redirect("/envelope/" + util.IntToString(envelopeHeight))
	}
}
