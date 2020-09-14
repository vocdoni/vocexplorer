package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	dvotedb "gitlab.com/vocdoni/go-dvote/db"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/api/rpc"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/db"
	"gitlab.com/vocdoni/vocexplorer/router"
)

func newConfig() (*config.MainCfg, error) {
	var cfg config.MainCfg
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	flag.StringVar(&cfg.DataDir, "dataDir", home+"/.vocexplorer", "directory where data is stored")
	cfg.Global.GatewayHost = *flag.String("gatewayHost", "ws://0.0.0.0:9090/dvote", "gateway API host to connect to")
	cfg.Global.TendermintHost = *flag.String("tendermintHost", "http://0.0.0.0:26657/websocket", "gateway API host to connect to")
	cfg.Global.RefreshTime = *flag.Int("refreshTime", 10, "Number of seconds between each content refresh")
	cfg.Global.Detached = *flag.Bool("detached", false, "run database in detached mode")
	cfg.DisableGzip = *flag.Bool("disableGzip", false, "use to disable gzip compression on web server")
	cfg.HostURL = *flag.String("hostURL", "http://localhost:8081", "url to host block explorer")
	cfg.ChainID = *flag.String("chainID", "", "chain ID to fall back on if running in detached mode")
	cfg.LogLevel = *flag.String("logLevel", "error", "log level <debug, info, warn, error>")
	flag.Parse()

	// setting up viper
	viper := viper.New()
	viper.SetConfigName("vocexplorer")
	viper.SetConfigType("yml")
	viper.SetEnvPrefix("VOCEXPLORER")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set FlagVars first
	viper.BindPFlag("dataDir", flag.Lookup("dataDir"))
	cfg.DataDir = viper.GetString("dataDir")

	// Add viper config path (now we know it)
	viper.AddConfigPath(cfg.DataDir)

	viper.BindPFlag("global.detached", flag.Lookup("detached"))
	viper.BindPFlag("global.gatewayHost", flag.Lookup("gatewayHost"))
	viper.BindPFlag("global.tendermintHost", flag.Lookup("tendermintHost"))
	viper.BindPFlag("global.refreshTime", flag.Lookup("refreshTime"))
	viper.BindPFlag("disableGzip", flag.Lookup("disableGzip"))
	viper.BindPFlag("hostURL", flag.Lookup("hostURL"))
	viper.BindPFlag("chainID", flag.Lookup("chainID"))
	viper.BindPFlag("logLevel", flag.Lookup("logLevel"))

	var cfgError error
	_, err = os.Stat(cfg.DataDir + "/vocexplorer.yml")
	if os.IsNotExist(err) {
		log.Infof("creating new config file in %s", cfg.DataDir)
		// creating config folder if not exists
		err = os.MkdirAll(cfg.DataDir, os.ModePerm)
		if err != nil {
			cfgError = fmt.Errorf("cannot create data directory: %s", err)
		}
		// create config file if not exists
		if err := viper.SafeWriteConfig(); err != nil {
			cfgError = fmt.Errorf("cannot write config file into config dir: %s", err)
		}
	} else {
		// read config file
		err = viper.ReadInConfig()
		if err != nil {
			cfgError = fmt.Errorf("cannot read loaded config file in %s: %s", cfg.DataDir, err)
		}
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		cfgError = fmt.Errorf("cannot unmarshal loaded config file: %s", err)
	}

	return &cfg, cfgError
}

func main() {
	cfg, err := newConfig()
	if cfg == nil {
		log.Fatal("cannot read configuration")
	}
	if err != nil {
		log.Error(err)
	}
	log.Init(cfg.LogLevel, "stdout")

	if _, err := os.Stat("./static/wasm_exec.js"); os.IsNotExist(err) {
		panic("File not found ./static/wasm_exec.js : find it in $GOROOT/misc/wasm/ note it must be from the same version of go used during compiling")
	}

	var d *dvotedb.BadgerDB
	if !cfg.Global.Detached {
		// Get ChainID for db directory
		tmClient, ok := api.StartTendermint(cfg.Global.TendermintHost)
		if !ok {
			cfg.Global.Detached = true
		} else {
			gen, err := rpc.Genesis(tmClient)
			if err != nil {
				log.Fatal(err)
			}
			cfg.ChainID = gen.Genesis.ChainID
			d, err = db.NewDB(cfg.DataDir, cfg.ChainID)
			if err != nil {
				log.Error(err)
				cfg.Global.Detached = true
			} else {
				go db.UpdateDB(d, &cfg.Global.Detached, cfg.Global.TendermintHost, cfg.Global.GatewayHost)
			}
		}
	}
	if cfg.Global.Detached {
		log.Infof("Running in detached mode")
		d, err = db.NewDB(cfg.DataDir, cfg.ChainID)
		if err != nil {
			log.Fatal(err)
		}
	}

	//Convert host url to localhost if using internal docker network
	cfg.Global.GatewayHost = strings.Replace(cfg.Global.GatewayHost, "dvotenode", "localhost", 1)
	cfg.Global.TendermintHost = strings.Replace(cfg.Global.TendermintHost, "dvotenode", "localhost", 1)

	urlR, err := url.Parse(cfg.HostURL)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("Server on: %v\n", urlR)

	r := mux.NewRouter()
	router.RegisterRoutes(r, &cfg.Global, d)

	s := &http.Server{
		Addr:           urlR.Host,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   40 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if cfg.DisableGzip {
		s.Handler = r
		if err = s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	} else {
		h := gziphandler.GzipHandler(r)
		s.Handler = h
		if err = s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}

}
