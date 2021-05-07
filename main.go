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
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/router"
	"go.vocdoni.io/dvote/log"
)

func newConfig() (*config.MainCfg, error) {
	cfg := new(config.MainCfg)
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	flag.StringVar(&cfg.DataDir, "dataDir", home+"/.vocexplorer", "directory where data is stored")
	cfg.Global.RefreshTime = *flag.Int("refreshTime", 10, "Number of seconds between each content refresh")
	cfg.Global.Environment = *flag.String("environment", "", "vochain environment, ie. \"dev\", \"stg\", \"\"")
	cfg.Global.GatewayUrl = *flag.String("gatewayUrl", "http://0.0.0.0:9090/dvote", "URL for the gateway to query for data")
	cfg.DisableGzip = *flag.Bool("disableGzip", false, "use to disable gzip compression on web server")
	cfg.HostURL = *flag.String("hostURL", "http://localhost:8081", "url to host block explorer")
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

	viper.BindPFlag("global.refreshTime", flag.Lookup("refreshTime"))
	viper.BindPFlag("global.environment", flag.Lookup("environment"))
	viper.BindPFlag("global.gatewayUrl", flag.Lookup("gatewayUrl"))
	viper.BindPFlag("disableGzip", flag.Lookup("disableGzip"))
	viper.BindPFlag("hostURL", flag.Lookup("hostURL"))
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

	return cfg, cfgError
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

	urlR, err := url.Parse(cfg.HostURL)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("Server on: %v\n", urlR)

	c, err := client.New(cfg.Global.GatewayUrl)
	if err != nil {
		log.Errorf("could not connect to gateway %s: %s", cfg.Global.GatewayUrl, err)
		return
	}

	for i := 0; i < 10; i++ {
		print("lifff")
		err = c.GetGatewayInfo()
		if err == nil {
			break
		}
		if i == 9 {
			log.Errorf("could not use gateway %s: %s", cfg.Global.GatewayUrl, err)
			return
		}
		time.Sleep(5 * time.Second)
	}

	log.Infof("Connected to gateway: %v\n", cfg.Global.GatewayUrl)

	r := mux.NewRouter()
	router.RegisterRoutes(r, &cfg.Global)

	s := &http.Server{
		Addr:         urlR.Host,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 60 * time.Second,
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
