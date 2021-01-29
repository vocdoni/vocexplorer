module github.com/vocdoni/vocexplorer/frontend

go 1.15

require (
	github.com/dustin/go-humanize v1.0.0
	github.com/hexops/vecty v0.5.1-0.20200816075853-64e387e2b2b3
	github.com/vocdoni/vocexplorer v0.0.0-20200903194749-046f8672292b
	go.vocdoni.io/dvote v0.6.1-0.20210129171551-b0e1cf9426b4
	go.vocdoni.io/proto v0.1.8
	google.golang.org/protobuf v1.25.0
	marwan.io/vecty-router v0.0.0-20200914150808-f30c81f0deb5
)

replace github.com/vocdoni/vocexplorer => ../

replace gitlab.com/vocodni/vocexplorer/frontend => ./
