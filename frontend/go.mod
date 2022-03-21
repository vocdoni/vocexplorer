module gitlab.com/vocdoni/vocexplorer/frontend

go 1.15

require (
	github.com/dustin/go-humanize v1.0.0
	github.com/hexops/vecty v0.5.1-0.20200816075853-64e387e2b2b3
	github.com/vocdoni/vocexplorer v0.0.0-20210326174540-19b196ebf6e0
	gitlab.com/vocdoni/vocexplorer v0.0.0-20200903194749-046f8672292b
	go.vocdoni.io/dvote v1.0.4-0.20220321130928-65cfa3e0ac55
	go.vocdoni.io/proto v1.13.3-0.20220203130255-cbdb9679ec7c
	google.golang.org/protobuf v1.27.1
	marwan.io/vecty-router v0.0.0-20200914150808-f30c81f0deb5
)

replace gitlab.com/vocdoni/vocexplorer => ../
