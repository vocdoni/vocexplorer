// package rpc

// import rpchttp "github.com/tendermint/rpc/client/http"
// import "github.com/tendermint/tendermint/types"

// client := rpchttp.New("tcp:127.0.0.1:26657", "/websocket")
// err := client.Start()
// if err != nil {
//   handle error
// }
// defer client.Stop()
// ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
// defer cancel()
// query := "tm.event = 'Tx' AND tx.height = 3"
// txs, err := client.Subscribe(ctx, "test-client", query)
// if err != nil {
//   handle error
// }

// go func() {
//  for e := range txs {
//    fmt.Println("got ", e.Data.(types.EventDataTx))
//    }
// }()