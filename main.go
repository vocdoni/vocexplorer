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
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/db"
	"gitlab.com/vocdoni/vocexplorer/router"
	"gitlab.com/vocdoni/vocexplorer/util"
)

func newConfig() (*config.MainCfg, error) {
	var cfg config.MainCfg
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	flag.StringVar(&cfg.DataDir, "dataDir", home+"/.vocexplorer", "directory where data is stored")
	cfg.Global.GatewayHost = *flag.String("gatewayHost", "0.0.0.0:9090", "gateway API host to connect to")
	cfg.Global.TendermintHost = *flag.String("tendermintHost", "0.0.0.0:26657", "gateway API host to connect to")
	cfg.Global.RefreshTime = *flag.Int("refreshTime", 5, "Number of seconds between each content refresh")
	cfg.DisableGzip = *flag.Bool("disableGzip", false, "use to disable gzip compression on web server")
	cfg.HostURL = *flag.String("hostURL", "http://localhost:8081", "url to host block explorer")
	cfg.LogLevel = *flag.String("logLevel", "error", "log level <debug, info, warn, error>")
	cfg.Detached = *flag.Bool("detached", false, "run database in detached mode")
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

	viper.BindPFlag("global.gatewayHost", flag.Lookup("gatewayHost"))
	viper.BindPFlag("global.tendermintHost", flag.Lookup("tendermintHost"))
	viper.BindPFlag("global.refreshTime", flag.Lookup("refreshTime"))
	viper.BindPFlag("disableGzip", flag.Lookup("disableGzip"))
	viper.BindPFlag("hostURL", flag.Lookup("hostURL"))
	viper.BindPFlag("logLevel", flag.Lookup("logLevel"))
	viper.BindPFlag("detached", flag.Lookup("detached"))

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

	d, err := db.NewDB(cfg.DataDir)
	if err != nil {
		log.Fatal(err)
	}
	if !cfg.Detached {
		go db.UpdateDB(d, cfg.Global.GatewayHost, cfg.Global.TendermintHost)
	} else {
		log.Infof("Running in detached mode")
	}

	//Convert host url to localhost if using internal docker network
	if strings.Contains(cfg.Global.GatewayHost, "dvotenode") {
		sub := strings.Split(cfg.Global.GatewayHost, ":")
		port := "9090"
		if len(sub) < 1 {
			port = sub[1]
		}
		cfg.Global.GatewayHost = "ws://localhost:" + port + "/dvote"
	} else {
		cfg.Global.GatewayHost = "ws://" + cfg.Global.GatewayHost + "/dvote"
	}

	//Convert host url to localhost if using internal docker network
	if strings.Contains(cfg.Global.TendermintHost, "dvotenode") {
		sub := strings.Split(cfg.Global.TendermintHost, ":")
		port := "26657"
		if len(sub) < 1 {
			port = sub[1]
		}
		cfg.Global.TendermintHost = "http://localhost:" + port
	} else {
		cfg.Global.TendermintHost = "http://" + cfg.Global.TendermintHost
	}
	urlR, err := url.Parse(cfg.HostURL)
	if util.ErrPrint(err) {
		return
	}
	log.Infof("Server on: %v\n", urlR)

	r := mux.NewRouter()
	router.RegisterRoutes(r, &cfg.Global, d)

	s := &http.Server{
		Addr:           urlR.Host,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if cfg.DisableGzip {
		s.Handler = r
		err = s.ListenAndServe()
		util.ErrFatal(err)
	} else {
		h := gziphandler.GzipHandler(r)
		s.Handler = h
		err = s.ListenAndServe()
		util.ErrFatal(err)
	}

}
