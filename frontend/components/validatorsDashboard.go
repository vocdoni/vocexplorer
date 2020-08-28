package components

import (
	"time"

	"github.com/gopherjs/vecty"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/dbapi"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// ValidatorsDashboardView renders the validators list dashboard
type ValidatorsDashboardView struct {
	vecty.Core
	totalValidators  int
	validatorList    [config.ListSize]*types.Validator
	serverConnected  bool
	gatewayConnected bool
	validatorIndex   int
	validatorRefresh chan int
	disableUpdate    bool
	quitCh           chan struct{}
	t                *rpc.TendermintInfo
}

// Render renders the ValidatorsDashboardView component
func (dash *ValidatorsDashboardView) Render() vecty.ComponentOrHTML {
	return Container(
		renderGatewayConnectionBanner(dash.gatewayConnected),
		renderServerConnectionBanner(dash.serverConnected),
		&ValidatorList{
			totalValidators: &dash.totalValidators,
			validatorList:   &dash.validatorList,
			refreshCh:       dash.validatorRefresh,
			disableUpdate:   &dash.disableUpdate,
		},
		&BlockchainInfo{
			T: dash.t,
		},
	)
}

// InitValidatorsDashboardView initializes the Validators dashboard view
func InitValidatorsDashboardView(t *rpc.TendermintInfo, dash *ValidatorsDashboardView, cfg *config.Cfg) *ValidatorsDashboardView {
	dash.t = t
	dash.quitCh = make(chan struct{})
	dash.validatorRefresh = make(chan int, 50)
	dash.validatorIndex = 0
	dash.disableUpdate = false
	dash.serverConnected = true
	BeforeUnload(func() {
		close(dash.quitCh)
	})
	go updateAndRenderValidatorsDashboard(dash, cfg)
	return dash
}

func updateAndRenderValidatorsDashboard(d *ValidatorsDashboardView, cfg *config.Cfg) {
	ticker := time.NewTicker(time.Duration(cfg.RefreshTime) * time.Second)
	updateValidatorsDashboard(d)
	vecty.Rerender(d)
	for {
		select {
		case <-d.quitCh:
			ticker.Stop()
			return
		case <-ticker.C:
			updateValidatorsDashboard(d)
			vecty.Rerender(d)
		case i := <-d.validatorRefresh:
		loop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-d.validatorRefresh:
				default:
					break loop
				}
			}
			d.validatorIndex = i
			oldValidators := d.totalValidators
			newHeight, _ := dbapi.GetValidatorCount()
			d.totalValidators = int(newHeight) - 1
			if i < 1 {
				oldValidators = d.totalValidators
			}
			updateValidators(d, util.Max(oldValidators-d.validatorIndex, config.ListSize))

			vecty.Rerender(d)
		}
	}
}

func updateValidatorsDashboard(d *ValidatorsDashboardView) {
	if !rpc.Ping(store.Tendermint) {
		d.gatewayConnected = false
	} else {
		d.gatewayConnected = true
	}
	if !dbapi.Ping() {
		d.serverConnected = false
	} else {
		d.serverConnected = true
	}
	updateHeight(d.t)
	rpc.UpdateTendermintInfo(store.Tendermint, d.t)
	newVal, ok := dbapi.GetValidatorCount()
	if ok {
		d.totalValidators = int(newVal)
	}
	if !d.disableUpdate {
		updateValidators(d, util.Max(d.totalValidators-d.validatorIndex, config.ListSize))
	}
}

func updateValidators(d *ValidatorsDashboardView, index int) {
	log.Infof("Getting Blocks from index %d", util.IntToString(index))
	list, ok := dbapi.GetValidatorList(index)
	if ok {
		reverseValidatorList(&list)
		d.validatorList = list
	}
}

func reverseValidatorList(list *[config.ListSize]*types.Validator) {
	for i := len(list)/2 - 1; i >= 0; i-- {
		opp := len(list) - 1 - i
		list[i], list[opp] = list[opp], list[i]
	}
}
