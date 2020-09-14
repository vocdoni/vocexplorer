module gitlab.com/vocdoni/vocexplorer/frontend

go 1.15

require (
	github.com/dustin/go-humanize v1.0.0
	github.com/golang/protobuf v1.4.2
	github.com/gopherjs/vecty v0.5.0
	github.com/hexops/vecty v0.5.1-0.20200816075853-64e387e2b2b3
	github.com/tendermint/tendermint v0.33.8
	gitlab.com/vocdoni/go-dvote v0.5.2
	gitlab.com/vocdoni/vocexplorer v0.0.0-20200903194749-046f8672292b
	golang.org/x/text v0.3.3
	marwan.io/vecty-router v0.0.0-20200914150808-f30c81f0deb5
	nhooyr.io/websocket v1.8.6
)

replace gitlab.com/vocdoni/vocexplorer => ../

replace gitlab.com/vocodni/vocexplorer/frontend => ./
