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
	"github.com/vocdoni/vocexplorer/config"
	"github.com/vocdoni/vocexplorer/db"
	"github.com/vocdoni/vocexplorer/router"
	dvotecfg "go.vocdoni.io/dvote/config"
	"go.vocdoni.io/dvote/log"
)

func newConfig() (*config.MainCfg, error) {
	var cfg config.MainCfg
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	flag.StringVar(&cfg.DataDir, "dataDir", home+"/.vocexplorer", "directory where data is stored")
	cfg.Global.RefreshTime = *flag.Int("refreshTime", 10, "Number of seconds between each content refresh")
	cfg.DisableGzip = *flag.Bool("disableGzip", false, "use to disable gzip compression on web server")
	cfg.HostURL = *flag.String("hostURL", "http://localhost:8081", "url to host block explorer")
	cfg.LogLevel = *flag.String("logLevel", "error", "log level <debug, info, warn, error>")
	cfg.Chain = *flag.String("chain", "main", "vochain network to connect to (eg. main, dev")

	// Vochain config
	cfg.VochainConfig = new(dvotecfg.VochainCfg)
	cfg.VochainConfig.P2PListen = *flag.String("vochainP2PListen", "0.0.0.0:26656", "p2p host and port to listent for the voting chain")
	cfg.VochainConfig.CreateGenesis = *flag.Bool("vochainCreateGenesis", false, "create own/testing genesis file on vochain")
	cfg.VochainConfig.Genesis = *flag.String("vochainGenesis", "", "use alternative genesis file for the voting chain")
	cfg.VochainConfig.LogLevel = *flag.String("vochainLogLevel", "error", "voting chain node log level")
	cfg.VochainConfig.Peers = *flag.StringArray("vochainPeers", []string{}, "coma separated list of p2p peers")
	cfg.VochainConfig.Seeds = *flag.StringArray("vochainSeeds", []string{}, "coma separated list of p2p seed nodes")
	cfg.VochainConfig.NodeKey = *flag.String("vochainNodeKey", "", "user alternative vochain private key (hexstring[64])")
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

	viper.BindPFlag("global.refreshTime", flag.Lookup("refreshTime"))
	viper.BindPFlag("disableGzip", flag.Lookup("disableGzip"))
	viper.BindPFlag("hostURL", flag.Lookup("hostURL"))
	viper.BindPFlag("chain", flag.Lookup("chain"))
	viper.BindPFlag("logLevel", flag.Lookup("logLevel"))

	viper.BindPFlag("vochainConfig.P2PListen", flag.Lookup("vochainP2PListen"))
	viper.BindPFlag("vochainConfig.CreateGenesis", flag.Lookup("vochainCreateGenesis"))
	viper.BindPFlag("vochainConfig.Genesis", flag.Lookup("vochainGenesis"))
	viper.BindPFlag("vochainConfig.LogLevel", flag.Lookup("vochainLogLevel"))
	viper.BindPFlag("vochainConfig.Peers", flag.Lookup("vochainPeers"))
	viper.BindPFlag("vochainConfig.Seeds", flag.Lookup("vochainSeeds"))
	viper.BindPFlag("vochainConfig.NodeKey", flag.Lookup("vochainNodeKey"))

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

	cfg.Global.Dev = cfg.Chain == "dev"

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
	d := db.NewDB(cfg)
	go d.UpdateDB()

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
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if cfg.DisableGzip {
		s.Handler = r
		if err = s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	} else {
		h, err := gziphandler.NewGzipLevelHandler(9)
		if err != nil {
			log.Error(err)
		}
		s.Handler = h(r)
		if err = s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}
}
