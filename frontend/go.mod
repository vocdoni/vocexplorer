module gitlab.com/vocdoni/vocexplorer/frontend

go 1.15

require (
	github.com/dustin/go-humanize v1.0.0
	github.com/hexops/vecty v0.5.1-0.20200816075853-64e387e2b2b3
	github.com/vocdoni/vocexplorer v0.0.0-20210326174540-19b196ebf6e0
	gitlab.com/vocdoni/vocexplorer v0.0.0-20200903194749-046f8672292b
	go.vocdoni.io/dvote v0.6.1-0.20210504142713-a76fcb54bf1b
	go.vocdoni.io/proto v1.0.3-0.20210427102429-33cf17b307f0
	google.golang.org/protobuf v1.25.0
	marwan.io/vecty-router v0.0.0-20200914150808-f30c81f0deb5
)

replace gitlab.com/vocdoni/vocexplorer => ../
