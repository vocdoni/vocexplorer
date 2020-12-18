module gitlab.com/vocdoni/vocexplorer

go 1.14

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/go-kit/kit v0.10.1-0.20200710014002-02c7c016dd45 // indirect
	github.com/golang/protobuf v1.4.3
	github.com/gorilla/mux v1.7.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/tendermint/tendermint v0.34.0
	github.com/vocdoni/dvote-protobuf v0.1.3
	go.vocdoni.io/dvote v0.6.1-0.20201217155643-aa4a86ec19da
	google.golang.org/protobuf v1.25.0
	honnef.co/go/tools v0.0.1-2020.1.5 // indirect
)

replace github.com/tendermint/tendermint => github.com/vocdoni/tendermint v0.34.0-rc4.0.20201209151525-75721cb94a61
