package vochain

import (
	"time"

	"gitlab.com/vocdoni/vocexplorer/config"
	"go.vocdoni.io/dvote/log"
	"go.vocdoni.io/dvote/service"
	"go.vocdoni.io/dvote/vochain"
	"go.vocdoni.io/dvote/vochain/scrutinizer"
	"go.vocdoni.io/dvote/vochain/vochaininfo"
)

const MaxListIterations = int64(64)

// VochainService contains a scrutinizer node
type VochainService struct {
	app   *vochain.BaseApplication
	scrut *scrutinizer.Scrutinizer
	info  *vochaininfo.VochainInfo
}

// InitVochain starts up a VochainService
func InitVochain(cfg *config.MainCfg) (*VochainService, error) {
	var err error
	vs := VochainService{}
	cfg.VochainConfig.DataDir = cfg.DataDir + "/vochain"
	cfg.VochainConfig.Chain = cfg.Chain
	cfg.VochainConfig.LogLevelMemPool = "info"
	if cfg.Chain != "dev" {
		cfg.VochainConfig.Dev = false
	} else {
		cfg.VochainConfig.Dev = true
	}
	vs.app, vs.scrut, vs.info, err = service.Vochain(cfg.VochainConfig, true, false, nil, nil)
	if err != nil {
		return nil, err
	}
	// Wait for Vochain to be ready
	var h, hPrev int64
	for vs.app.Node == nil {
		hPrev = h
		time.Sleep(time.Second * 5)
		if header := vs.app.State.Header(true); header != nil {
			h = header.Height
		}
		log.Infof("[vochain info] replaying block %d at %d b/s",
			h, (h-hPrev)/5)
	}
	log.Info("Started vochain service")
	return &vs, nil
}

// Close closes the VochainService
func (vs *VochainService) Close() {
	vs.info.Close()
	vs.app.Node.Stop()
	vs.app.Node.Wait()
	vs.scrut.Storage.Close()
	vs.app.State.Store.Close()
}
