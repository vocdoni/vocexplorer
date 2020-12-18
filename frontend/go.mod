module gitlab.com/vocdoni/vocexplorer/frontend

go 1.15

require (
	github.com/dustin/go-humanize v1.0.0
	github.com/hexops/vecty v0.5.1-0.20200816075853-64e387e2b2b3
	github.com/tendermint/tendermint v0.34.0
	github.com/vocdoni/dvote-protobuf v0.1.3
	gitlab.com/vocdoni/vocexplorer v0.0.0-20200903194749-046f8672292b
	go.vocdoni.io/dvote v0.6.1-0.20201217155643-aa4a86ec19da
	google.golang.org/protobuf v1.25.0
	marwan.io/vecty-router v0.0.0-20200914150808-f30c81f0deb5
)

replace gitlab.com/vocdoni/vocexplorer => ../

replace gitlab.com/vocodni/vocexplorer/frontend => ./

replace github.com/tendermint/tendermint => github.com/vocdoni/tendermint v0.34.0-rc4.0.20201209151525-75721cb94a61
