module gitlab.com/vocdoni/vocexplorer/frontend

go 1.15

require (
	github.com/dustin/go-humanize v1.0.0
	github.com/hexops/vecty v0.5.1-0.20200816075853-64e387e2b2b3
	github.com/p4u/recws v1.2.2-0.20201005083112-7be7f9397e75 // indirect
	github.com/tendermint/tendermint v0.33.8 // indirect
	gitlab.com/vocdoni/go-dvote v0.6.1-0.20201113175154-a81a55b0650e
	gitlab.com/vocdoni/vocexplorer v0.0.0-20200903194749-046f8672292b
	marwan.io/vecty-router v0.0.0-20200914150808-f30c81f0deb5
)

replace gitlab.com/vocdoni/vocexplorer => ../

replace gitlab.com/vocodni/vocexplorer/frontend => ./
