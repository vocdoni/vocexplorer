package rpcinit

import (
	"fmt"
	"time"

	gohttp "net/http"

	"github.com/tendermint/tendermint/rpc/client/http"
	jsonrpcclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// StartClient initializes an http tendermint api client on websockets
func StartClient(host string) *http.HTTP {
	fmt.Println("connecting to " + host)
	tClient, err := initClient(host)
	if util.ErrPrint(err) {
		return nil
	}
	return tClient
}

func initClient(host string) (*http.HTTP, error) {
	httpClient, err := jsonrpcclient.DefaultHTTPClient(host)
	if err != nil {
		return nil, err
	}
	httpClient.Timeout = 2 * time.Second
	// Increase max idle connections. This fixes issue with too many concurrent requests, as described here: https://github.com/golang/go/issues/16012
	httpClient.Transport.(*gohttp.Transport).MaxIdleConnsPerHost = 10000
	c, err := http.NewWithClient(host, "/websocket", httpClient)
	// c, err := http.NewWithTimeout(host, "/websocket", 2)
	if err != nil {
		return nil, err
	}
	_, err = c.Status()
	if err != nil {
		return nil, err
	}
	return c, nil
}
