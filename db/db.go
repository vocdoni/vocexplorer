package db

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/ptypes"
	tmhttp "github.com/tendermint/tendermint/rpc/client/http"
	tmtypes "github.com/tendermint/tendermint/types"
	"gitlab.com/vocdoni/go-dvote/crypto/ethereum"
	dvotedb "gitlab.com/vocdoni/go-dvote/db"
	"gitlab.com/vocdoni/go-dvote/log"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/go-dvote/vochain"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
	"google.golang.org/protobuf/proto"
)

// NewDB initializes a badger database at the given path
func NewDB(path, chainID string) (*dvotedb.BadgerDB, error) {
	log.Infof("Initializing database at " + path + "/" + chainID)
	return dvotedb.NewBadgerDB(path + "/" + chainID)
}

// UpdateDB continuously updates the database by calling dvote & tendermint apis
func UpdateDB(d *dvotedb.BadgerDB, gwHost, tmHost string) {

	ping := pingGateway(gwHost)
	if !ping {
		log.Warn("Gateway Client is not running. Running as detached database")
		return
	}
	// Init tendermint client
	tClient, up := StartTendermint(tmHost)
	if !up {
		log.Warn("Cannot connect to tendermint client. Running as detached database")
		return
	}

	log.Debugf("Connected to " + tmHost)
	// defer (*cancel)()
	i := 0
	for {
		updateBlockList(d, tClient)
		// Update validators less frequently than blocks
		if i%20 == 0 {
			updateValidatorList(d, tClient)
		}
		updateEntityList(d)
		updateProcessList(d)
		time.Sleep(config.DBWaitTime * time.Millisecond)
		i++
	}
}

func updateValidatorList(d *dvotedb.BadgerDB, c *tmhttp.HTTP) {
	latestBlockHeight := getHeight(d, config.LatestBlockHeightKey, 1)

	batch := d.NewBatch()
	fetchValidators(latestBlockHeight.GetHeight(), c, batch)
	util.ErrPrint(batch.Write())
}

func updateBlockList(d *dvotedb.BadgerDB, c *tmhttp.HTTP) {
	// Fetch latest block & tx heights
	latestBlockHeight := getHeight(d, config.LatestBlockHeightKey, 1)
	latestTxHeight := getHeight(d, config.LatestTxHeightKey, 1)
	latestEnvelopeHeight := getHeight(d, config.LatestEnvelopeHeightKey, 0)

	status, err := c.Status()
	if err != nil {
		log.Error(err)
	}
	gwBlockHeight := status.SyncInfo.LatestBlockHeight

	// Wait for new blocks to be available
	for gwBlockHeight-latestBlockHeight.GetHeight() < 1 {
		time.Sleep(500 * time.Millisecond)
		status, err := c.Status()
		if err != nil {
			log.Error(err)
		}
		gwBlockHeight = status.SyncInfo.LatestBlockHeight
	}

	batch := d.NewBatch()

	i := int64(0)
	numNewBlocks := util.Min(config.NumBlockUpdates, int(gwBlockHeight-latestBlockHeight.GetHeight()))
	// Array of new tx id's. Each goroutine can only access its assigned index, making this array thread-safe as long as all goroutines exit before read access
	txsList := make([]tmtypes.Txs, numNewBlocks)
	complete := make(chan struct{}, config.NumBlockUpdates)
	for ; int(i) < numNewBlocks; i++ {
		go fetchBlock(i+latestBlockHeight.GetHeight(), &batch, c, complete, &txsList[i])
	}
	num := 0
	if i > 0 {
		// Sync: wait here for all goroutines to complete
		for range complete {
			if num >= numNewBlocks-1 {
				break
			}
			num++
		}
		log.Debugf("Setting block %d ", latestBlockHeight.GetHeight()+i)

		complete = make(chan struct{}, len(txsList))
		for _, txs := range txsList {
			if len(txs) > 0 {
				go updateTxs(latestTxHeight.GetHeight(), txs, c, batch, complete, latestEnvelopeHeight)
				latestTxHeight.Height += int64(len(txs))
			} else {
				complete <- struct{}{}
			}
		}

		// Sync: wait here for all goroutines to complete
		num = 0
		for range complete {
			if num >= len(txsList)-1 {
				break
			}
			num++
		}
		blockHeight := types.Height{Height: latestBlockHeight.GetHeight() + i}
		encBlockHeight, err := proto.Marshal(&blockHeight)
		util.ErrPrint(err)
		encTxHeight, err := proto.Marshal(latestTxHeight)
		util.ErrPrint(err)
		encEnvHeight, err := proto.Marshal(latestEnvelopeHeight)
		util.ErrPrint(err)

		batch.Put([]byte(config.LatestTxHeightKey), encTxHeight)
		batch.Put([]byte(config.LatestBlockHeightKey), encBlockHeight)
		batch.Put([]byte(config.LatestEnvelopeHeightKey), encEnvHeight)
		util.ErrPrint(batch.Write())
	}

}

func fetchValidators(blockHeight int64, c *tmhttp.HTTP, batch dvotedb.Batch) {
	maxPerPage := 100
	page := 0
	resultValidators, err := c.Validators(&blockHeight, page, 100)
	util.ErrPrint(err)
	// Check if there are more validators.
	for len(resultValidators.Validators) == maxPerPage {
		moreValidators, err := c.Validators(&blockHeight, page, maxPerPage)
		util.ErrPrint(err)

		if len(resultValidators.Validators) > 0 {
			resultValidators.Validators = append(resultValidators.Validators, moreValidators.Validators...)
		}
		page++
	}
	// Cast each validator as storage struct, marshal, write to batch
	for _, validator := range resultValidators.Validators {
		var storeValidator types.Validator
		storeValidator.Address = validator.Address
		storeValidator.ProposerPriority = validator.ProposerPriority
		storeValidator.VotingPower = validator.VotingPower
		storeValidator.PubKey = validator.PubKey.Bytes()
		encValidator, err := proto.Marshal(&storeValidator)
		util.ErrPrint(err)
		batch.Put(append([]byte(config.ValidatorPrefix), validator.Address...), encValidator)
		log.Debugf("Validator address: %s", util.HexToString(storeValidator.GetAddress()))
	}
	log.Debugf("Fetched %d validators at block height %d", len(resultValidators.Validators), blockHeight)
}

func updateTxs(startTxHeight int64, txs tmtypes.Txs, c *tmhttp.HTTP, batch dvotedb.Batch, complete chan<- struct{}, envHeight *types.Height) {
	numTxs := int64(-1)
	var height int64
	for i, tx := range txs {
		numTxs = int64(i)
		txRes := rpc.GetTransaction(c, tx.Hash())

		txHashKey := append([]byte(config.TxHashPrefix), tx.Hash()...)
		// Marshal TxResult to bytes for protobuf encoding
		result, err := json.Marshal(txRes.TxResult)
		util.ErrPrint(err)
		txStore := types.StoreTx{
			Height:   txRes.Height,
			TxHeight: startTxHeight,
			Tx:       txRes.Tx,
			TxResult: result,
			Index:    txRes.Index,
		}
		storeEnvelope(txStore.Tx, envHeight, batch)
		txVal, err := proto.Marshal(&txStore)
		util.ErrPrint(err)
		util.ErrPrint(err)
		batch.Put(txHashKey, txVal)
		txHeightKey := []byte(config.TxHeightPrefix + util.IntToString(txStore.GetTxHeight()))
		batch.Put(txHeightKey, tx.Hash())
		if i == 0 {
			height = txRes.Height
		}
		startTxHeight++
	}
	if numTxs > 0 {
		log.Debugf("%d transactions logged at block %d, height %d", numTxs+1, height, startTxHeight)
	}
	complete <- struct{}{}
}

func fetchBlock(height int64, batch *dvotedb.Batch, c *tmhttp.HTTP, complete chan<- struct{}, txs *tmtypes.Txs) {
	// Signal
	defer func() {
		complete <- struct{}{}
	}()
	// Thread-safe api request
	res, err := c.Block(&height)
	// If error is returned, try the request more times, then fatal.
	if err != nil {
		log.Error(err)
		for errs := 0; ; errs++ {
			if errs > 10 {
				log.Fatal("Blockchain RPC Disconnected")
				return
			}
			res, err = c.Block(&height)
			if err == nil {
				break
			}
		}
	}
	var block types.StoreBlock
	block.NumTxs = int64(len(res.Block.Data.Txs))
	block.Hash = res.BlockID.Hash
	block.Height = res.Block.Header.Height
	block.Proposer = res.Block.ProposerAddress
	tm, err := ptypes.TimestampProto(res.Block.Header.Time)
	util.ErrPrint(err)
	block.Time = tm

	*txs = res.Block.Data.Txs

	bodyValue, err := proto.Marshal(&block)
	if err != nil {
		log.Error(err)
	}

	blockHeightKey := append([]byte(config.BlockHeightPrefix), []byte(util.IntToString(block.GetHeight()))...)
	blockHashKey := append([]byte(config.BlockHashPrefix), block.GetHash()...)
	hashValue := block.Hash

	// Thread-safe batch operations
	(*batch).Put(blockHashKey, bodyValue)
	(*batch).Put(blockHeightKey, hashValue)
}

func updateEntityList(d *dvotedb.BadgerDB) {

}

func updateProcessList(d *dvotedb.BadgerDB) {

}

// listHashesByHeight returns a list of hashes given integer keys
func listHashesByHeight(d *dvotedb.BadgerDB, max, height int, prefix string) [][]byte {
	if max > 64 {
		max = 64
	}
	var hashList [][]byte
	for ; max > 0; max-- {
		key := []byte(prefix + util.IntToString(height))
		has, err := d.Has(key)
		if !has || util.ErrPrint(err) {
			break
		}
		val, err := d.Get(key)
		if err != nil {
			log.Error(err)
		}
		hashList = append(hashList, val)
		height--
	}
	return hashList
}

func pingGateway(host string) bool {
	pingClient := http.Client{
		Timeout: 5 * time.Second,
	}
	for i := 0; ; i++ {
		if i > 10 {
			return false
		}
		resp, err := pingClient.Get("http://" + host + "/ping")
		if err != nil {
			log.Debug(err.Error())
			time.Sleep(2 * time.Second)
			continue
		}
		body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
		if err != nil {
			log.Debug(err.Error())
			time.Sleep(time.Second)
			continue
		}
		if string(body) != "pong" {
			log.Warn("Gateway ping not yet available")
		} else {
			return true
		}
	}
}

// func startGateway(host string) (*client.Client, *context.CancelFunc, bool) {
// 	for i := 0; ; i++ {
// 		if i > 20 {
// 			return nil, nil, false
// 		}
// 		gwClient, cancel := client.InitGateway(host)
// 		if gwClient == nil {
// 			time.Sleep(5 * time.Second)
// 			continue
// 		} else {
// 			return gwClient, &cancel, true
// 		}
// 	}
// }

//StartTendermint starts the tendermint client
func StartTendermint(host string) (*tmhttp.HTTP, bool) {
	for i := 0; ; i++ {
		if i > 20 {
			return nil, false
		}
		tmClient := rpc.StartClient("http://" + host)
		if tmClient == nil {
			time.Sleep(1 * time.Second)
			continue
		} else {
			return tmClient, true
		}
	}
}

func getHeight(d *dvotedb.BadgerDB, key string, def int64) *types.Height {
	height := &types.Height{Height: def}
	has, err := d.Has([]byte(key))
	if err != nil {
		log.Error(err)
	}
	if has {
		val, err := d.Get([]byte(key))
		util.ErrPrint(err)
		err = proto.Unmarshal(val, height)
		util.ErrPrint(err)
	}
	return height
}

func storeEnvelope(tx tmtypes.Tx, height *types.Height, batch dvotedb.Batch) {
	var rawTx dvotetypes.Tx
	err := json.Unmarshal(tx, &rawTx)
	util.ErrPrint(err)
	if rawTx.Type == "vote" {
		myHeight := atomic.AddInt64(&height.Height, 1)
		var voteTx dvotetypes.VoteTx
		err = json.Unmarshal(tx, &voteTx)
		util.ErrPrint(err)

		// Write vote package
		votePackage := types.Envelope{
			ProcessID: voteTx.ProcessID,
			Package:   voteTx.VotePackage,
		}

		// Generate nullifier as in go-dvote vochain/transaction.go
		signature := voteTx.Signature
		voteTx.Signature = ""
		voteTx.Type = ""
		voteBytes, err := json.Marshal(&voteTx)
		util.ErrPrint(err)
		pubKey, err := ethereum.PubKeyFromSignature(voteBytes, signature)
		if err != nil {
			log.Errorf("cannot extract public key from signature (%s)", err)
		}
		addr, err := ethereum.AddrFromPublicKey(pubKey)
		if err != nil {
			log.Errorf("cannot extract address from public key: (%s)", err)
		}
		votePackage.Nullifier, err = vochain.GenerateNullifier(addr, votePackage.ProcessID)
		if err != nil {
			log.Errorf("cannot generate nullifier: (%s)", err)
		}
		for _, index := range voteTx.EncryptionKeyIndexes {
			votePackage.EncryptionKeyIndexes = append(votePackage.EncryptionKeyIndexes, int32(index))
		}
		rawEnvelope, err := proto.Marshal(&votePackage)
		util.ErrPrint(err)
		packageKey := append([]byte(config.EnvPackagePrefix), []byte(util.IntToString(myHeight))...)
		batch.Put(packageKey, rawEnvelope)

		// Write nullifier:height
		storeHeight := types.Height{Height: myHeight}
		rawHeight, err := proto.Marshal(&storeHeight)
		util.ErrPrint(err)
		nullifier, err := hex.DecodeString(util.StripHexString(voteTx.Nullifier))
		util.ErrPrint(err)
		nullifierKey := append([]byte(config.EnvNullifierPrefix), nullifier...)
		batch.Put(nullifierKey, rawHeight)

		// Write pid|height:height
		heightBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(heightBytes, uint64(myHeight))
		PIDBytes, err := hex.DecodeString(util.StripHexString(voteTx.ProcessID))
		util.ErrPrint(err)
		heightKey := append([]byte(config.EnvPIDPrefix), heightBytes...)
		heightKey = append(heightKey, PIDBytes...)
		batch.Put(heightKey, rawHeight)

		log.Debugf("Stored envelope %s of process %s at height %d", votePackage.Nullifier, voteTx.ProcessID, myHeight)
	}
}
