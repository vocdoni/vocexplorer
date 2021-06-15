module gitlab.com/vocdoni/vocexplorer

go 1.14

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/gorilla/mux v1.7.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	go.vocdoni.io/dvote v0.6.1-0.20210614152518-02aa5537c3df
	go.vocdoni.io/proto v1.0.4-0.20210525130734-c9e1ff675866
	honnef.co/go/tools v0.0.1-2020.1.5 // indirect
	nhooyr.io/websocket v1.8.6
)

// Newer versions of the fuse module removed support for MacOS.
// Unfortunately, its downstream users don't handle this properly,
// so our builds simply break for GOOS=darwin.
// Until either upstream or downstream solve this properly,
// force a downgrade to the commit right before support was dropped.
// It's also possible to use downstream's -tags=nofuse, but that's manual.
// TODO: remove once we've untangled module dep loops.
replace bazil.org/fuse => bazil.org/fuse v0.0.0-20200407214033-5883e5a4b512
