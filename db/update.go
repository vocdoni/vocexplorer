package db

import (
	"encoding/hex"
	"encoding/json"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/ptypes"
	tmtypes "github.com/tendermint/tendermint/types"
	"gitlab.com/vocdoni/go-dvote/crypto/ethereum"
	dvotedb "gitlab.com/vocdoni/go-dvote/db"
	"gitlab.com/vocdoni/go-dvote/log"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	dvoteutil "gitlab.com/vocdoni/go-dvote/util"
	"gitlab.com/vocdoni/go-dvote/vochain"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/api/rpc"
	"gitlab.com/vocdoni/vocexplorer/config"
	voctypes "gitlab.com/vocdoni/vocexplorer/proto"
	"gitlab.com/vocdoni/vocexplorer/util"
	"google.golang.org/protobuf/proto"
)

func updateBlockList(d *dvotedb.BadgerDB, t *rpc.TendermintRPC) {
	// Fetch latest block & tx heights
	latestBlockHeight := GetHeight(d, config.LatestBlockHeightKey, 1)
	latestTxHeight := GetHeight(d, config.LatestTxHeightKey, 1)
	latestEnvelopeCount := GetHeight(d, config.LatestEnvelopeCountKey, 0)

	// Get Height maps: stored in map object so each update isn't slow db-write/get
	// Map of validator:num blocks
	valMap := GetHeightMap(d, config.ValidatorHeightMapKey)
	valMapMutex := new(sync.Mutex)
	// Map of pid:num envelopes
	procEnvHeightMap := GetHeightMap(d, config.ProcessEnvelopeCountMapKey)
	procEnvHeightMapMutex := new(sync.Mutex)

	status, err := t.Status()
	// If error is returned, try the request more times, then fatal.
	if err != nil {
		if strings.Contains(err.Error(), "WebSocket closed") {
			exit <- struct{}{}
			return
		}
		log.Error(err)
		for errs := 0; ; errs++ {
			if errs > 10 {
				log.Error("Gateway Disconnected")
				exit <- struct{}{}
				return
			}
			status, err = t.Status()
			if err == nil {
				break
			}
		}
	}
	gwBlockHeight := status.SyncInfo.LatestBlockHeight

	// Wait for new blocks to be available
	for gwBlockHeight-latestBlockHeight.GetHeight() < 1 {
		time.Sleep(500 * time.Millisecond)
		return
	}

	batch := d.NewBatch()

	i := int64(0)
	numNewBlocks := util.Min(config.NumBlockUpdates, int(gwBlockHeight-latestBlockHeight.GetHeight()))
	// Array of new tx id's. Each goroutine can only access its assigned index, making this array thread-safe as long as all goroutines exit before read access
	txsList := make([]tmtypes.Txs, numNewBlocks)
	complete := make(chan struct{}, config.NumBlockUpdates)
	// nextHeight and myHeight channels synchronize goroutines before fetching validator block height, so blocks by validator are ordered by block height
	nextHeight := make(chan struct{})
	myHeight := make(chan struct{})
	for ; int(i) < numNewBlocks; i++ {
		go fetchBlock(i+latestBlockHeight.GetHeight(), &batch, t, complete, myHeight, nextHeight, &txsList[i], valMap, valMapMutex)
		if i == 0 {
			//Signal to the first block to start
			close(myHeight)
		}
		myHeight = nextHeight
		nextHeight = make(chan struct{})
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
				go updateTxs(latestTxHeight.GetHeight(), txs, t, batch, complete, latestEnvelopeCount, procEnvHeightMap, procEnvHeightMapMutex)
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
		rawValMap, err := proto.Marshal(valMap)
		if err != nil {
			log.Error(err)
		}
		batch.Put([]byte(config.ValidatorHeightMapKey), rawValMap)
		rawProcMap, err := proto.Marshal(procEnvHeightMap)
		if err != nil {
			log.Error(err)
		}
		batch.Put([]byte(config.ProcessEnvelopeCountMapKey), rawProcMap)
		blockHeight := voctypes.Height{Height: latestBlockHeight.GetHeight() + i}
		encBlockHeight, err := proto.Marshal(&blockHeight)
		if err != nil {
			log.Error(err)
		}
		encTxHeight, err := proto.Marshal(latestTxHeight)
		if err != nil {
			log.Error(err)
		}
		encEnvCount, err := proto.Marshal(latestEnvelopeCount)
		if err != nil {
			log.Error(err)
		}

		batch.Put([]byte(config.LatestTxHeightKey), encTxHeight)
		batch.Put([]byte(config.LatestBlockHeightKey), encBlockHeight)
		batch.Put([]byte(config.LatestEnvelopeCountKey), encEnvCount)
		if err := batch.Write(); err != nil {
			log.Error(err)
		}
	}

}

func fetchBlock(height int64, batch *dvotedb.Batch, t *rpc.TendermintRPC, complete, myHeight, nextHeight chan struct{}, txs *tmtypes.Txs, valMap *voctypes.HeightMap, valMapMutex *sync.Mutex) {
	// Signal
	defer func() {
		complete <- struct{}{}
	}()
	// Thread-safe api request
	res, err := t.Block(&height)
	// If error is returned, try the request more times, then fatal.
	if err != nil {
		log.Warn(err)
		for errs := 0; ; errs++ {
			if errs > 10 {
				log.Error("Gateway Disconnected")
				exit <- struct{}{}
				return
			}
			res, err = t.Block(&height)
			if err == nil {
				break
			}
		}
	}
	var block voctypes.StoreBlock
	block.NumTxs = int64(len(res.Block.Data.Txs))
	block.Hash = res.BlockID.Hash
	block.Height = res.Block.Header.Height
	block.Proposer = res.Block.ProposerAddress
	tm, err := ptypes.TimestampProto(res.Block.Header.Time)
	if err != nil {
		log.Error(err)
	}
	block.Time = tm

	*txs = res.Block.Data.Txs

	bodyValue, err := proto.Marshal(&block)
	if err != nil {
		log.Error(err)
	}

	// Wait for myHeight channel to close, this means fetchBlock for previous block has been assigned a validator block height
	<-myHeight
	// Update height of validator block belongs to
	valMapMutex.Lock()
	height, ok := valMap.Heights[util.HexToString(block.Proposer)]
	if !ok {
		height = 0
	}
	height++
	valMap.Heights[util.HexToString(block.Proposer)] = height
	valMapMutex.Unlock()
	// Signal to next block that I have been assigned a validator block height
	close(nextHeight)

	blockHeightKey := append([]byte(config.BlockHeightPrefix), util.EncodeInt(block.GetHeight())...)
	blockHashKey := append([]byte(config.BlockHashPrefix), block.GetHash()...)
	validatorHeightKey := append([]byte(config.BlockByValidatorPrefix), block.GetProposer()...)
	validatorHeightKey = append(validatorHeightKey, util.EncodeInt(height)...)
	hashValue := block.GetHash()

	// Thread-safe batch operations
	// Store hash:body
	(*batch).Put(blockHashKey, bodyValue)
	// Store globalheight:hash
	(*batch).Put(blockHeightKey, hashValue)
	// Store validator|heightbyValidator:hash
	(*batch).Put(validatorHeightKey, hashValue)
}

func updateValidatorList(d *dvotedb.BadgerDB, t *rpc.TendermintRPC) {
	latestBlockHeight := GetHeight(d, config.LatestBlockHeightKey, 0)
	latestValidatorCount := GetHeight(d, config.LatestValidatorCountKey, 0)
	if latestBlockHeight.GetHeight() > 0 {
		batch := d.NewBatch()
		fetchValidators(latestBlockHeight.GetHeight(), latestValidatorCount.GetHeight(), t, batch)
		if err := batch.Write(); err != nil {
			log.Error(err)
		}
	}
}

func fetchValidators(blockHeight, validatorCount int64, t *rpc.TendermintRPC, batch dvotedb.Batch) {
	maxPerPage := 100
	page := 0
	resultValidators, err := t.Validators(&blockHeight, page, 100)
	// If error is returned, try the request more times, then fatal.
	if err != nil {
		log.Error(err)
		for errs := 0; ; errs++ {
			if errs > 10 {
				log.Error("Gateway Disconnected")
				exit <- struct{}{}
				return
			}
			resultValidators, err = t.Validators(&blockHeight, page, 100)
			if err == nil {
				break
			}
		}
	}
	// Check if there are more validators.
	for len(resultValidators.Validators) == maxPerPage {
		moreValidators, err := t.Validators(&blockHeight, page, maxPerPage)
		// If error is returned, try the request more times, then fatal.
		if err != nil {
			log.Error(err)
			for errs := 0; ; errs++ {
				if errs > 10 {
					log.Error("Gateway Disconnected")
					exit <- struct{}{}
					return
				}
				moreValidators, err = t.Validators(&blockHeight, page, maxPerPage)
				if err == nil {
					break
				}
			}
		}

		if len(moreValidators.Validators) > 0 {
			resultValidators.Validators = append(resultValidators.Validators, moreValidators.Validators...)
		}
		page++
	}
	// Cast each validator as storage struct, marshal, write to batch
	for i, validator := range resultValidators.Validators {
		if i < int(validatorCount) {
			continue
		}
		validatorCount++
		var storeValidator voctypes.Validator
		storeValidator.Address = validator.Address
		storeValidator.Height = &voctypes.Height{Height: validatorCount}
		storeValidator.ProposerPriority = validator.ProposerPriority
		storeValidator.VotingPower = validator.VotingPower
		storeValidator.PubKey = validator.PubKey.Bytes()
		encValidator, err := proto.Marshal(&storeValidator)
		if err != nil {
			log.Error(err)
		}
		// Write id:validator
		batch.Put(append([]byte(config.ValidatorPrefix), validator.Address...), encValidator)
		// Write height:id
		batch.Put(append([]byte(config.ValidatorHeightPrefix), util.EncodeInt(storeValidator.Height.GetHeight())...), validator.Address)
	}
	// Write latest validator height
	rawHeight, err := proto.Marshal(&voctypes.Height{Height: validatorCount})
	if err != nil {
		log.Error(err)
	}
	batch.Put([]byte(config.LatestValidatorCountKey), rawHeight)
	log.Debugf("Fetched %d validators at block height %d", len(resultValidators.Validators), blockHeight)
}

func updateTxs(startTxHeight int64, txs tmtypes.Txs, t *rpc.TendermintRPC, batch dvotedb.Batch, complete chan<- struct{}, envHeight *voctypes.Height, procHeightMap *voctypes.HeightMap, procHeightMapMutex *sync.Mutex) {
	numTxs := int64(-1)
	var blockHeight int64
	for i, tx := range txs {
		numTxs = int64(i)
		txRes := api.GetTransaction(t, tx.Hash())

		txHashKey := append([]byte(config.TxHashPrefix), tx.Hash()...)
		// Marshal TxResult to bytes for protobuf encoding
		result, err := json.Marshal(txRes.TxResult)
		if err != nil {
			log.Error(err)
		}
		txStore := voctypes.StoreTx{
			Height:   txRes.Height,
			TxHeight: startTxHeight,
			Tx:       txRes.Tx,
			TxResult: result,
			Index:    txRes.Index,
		}
		// If voteTx, get envelope nullifier
		txStore.Nullifier = storeEnvelope(txStore.Tx, envHeight, procHeightMap, procHeightMapMutex, batch)
		txVal, err := proto.Marshal(&txStore)
		if err != nil {
			log.Error(err)
		}
		if err != nil {
			log.Error(err)
		}
		batch.Put(txHashKey, txVal)
		//Write height:tx hash
		txHeightKey := append([]byte(config.TxHeightPrefix), util.EncodeInt(txStore.GetTxHeight())...)
		batch.Put(txHeightKey, tx.Hash())
		if i == 0 {
			blockHeight = txRes.Height
		}
		startTxHeight++
	}
	if numTxs > 0 {
		log.Debugf("%d transactions logged at block %d, height %d", numTxs+1, blockHeight, startTxHeight)
	}
	complete <- struct{}{}
}

func updateEntityList(d *dvotedb.BadgerDB, c *api.GatewayClient) {
	localEntityHeight := GetHeight(d, config.LatestEntityCountKey, 0).GetHeight()
	gatewayEntityHeight, err := c.GetEntityCount()
	// If error is returned, try the request more times, then fatal.
	if err != nil {
		log.Error(err)
		for errs := 0; ; errs++ {
			if errs > 10 {
				log.Error("Gateway Disconnected")
				exit <- struct{}{}
				return
			}
			gatewayEntityHeight, err = c.GetEntityCount()
			if err == nil {
				break
			}
		}
	}
	if localEntityHeight >= gatewayEntityHeight {
		return
	}
	latestKey := append([]byte(config.EntityHeightPrefix), util.EncodeInt(int(localEntityHeight-1))...)
	latestEntity, err := d.Get(latestKey)
	if err != nil {
		latestEntity = []byte{}
	}
	log.Debugf("Getting entities from id %s", util.HexToString(latestEntity))
	newEntities, err := c.GetScrutinizerEntities(strings.ToLower(util.HexToString(latestEntity)))
	if len(newEntities) < 1 {
		log.Warn("No new entities fetched")
		return
	}
	heightMap := GetHeightMap(d, config.EntityProcessCountMapKey)

	// write new entities to db
	batch := d.NewBatch()
	i := 0
	entity := ""
	for i, entity = range newEntities {
		// Write entityHeight:entityID
		heightKey := append([]byte(config.EntityHeightPrefix), util.EncodeInt(int(localEntityHeight)+i)...)
		rawEntity, err := hex.DecodeString(util.TrimHex(entity))
		if err != nil {
			log.Error(err)
			break
		}
		batch.Put(heightKey, rawEntity)

		// Write entityID:[]
		// Write entityID:[]
		entityIDKey := append([]byte(config.EntityIDPrefix), rawEntity...)
		batch.Put(entityIDKey, []byte{})

		// Add new entity to height map with height of 0 so db will get new entity's processes
		if _, ok := heightMap.GetHeights()[entity]; ok {
			log.Warn("Fetched entity already stored")
		}
		heightMap.Heights[entity] = 0
	}

	rawValMap, err := proto.Marshal(heightMap)
	if err != nil {
		log.Error(err)
	}
	log.Debugf("Fetched %d new entities at height %d", len(newEntities), int(localEntityHeight)+i+1)

	// Write latest entity height
	encHeight := voctypes.Height{Height: localEntityHeight + int64(i) + 1}
	rawHeight, err := proto.Marshal(&encHeight)
	if err != nil {
		log.Error(err)
	}
	heightKey := []byte(config.LatestEntityCountKey)
	batch.Put(heightKey, rawHeight)

	// Write entity/process height map
	heightMapKey := []byte(config.EntityProcessCountMapKey)
	batch.Put(heightMapKey, rawValMap)
	batch.Write()
}

func updateProcessList(d *dvotedb.BadgerDB, c *api.GatewayClient) {
	localProcessHeight := GetHeight(d, config.LatestProcessCountKey, 0).GetHeight()
	gatewayProcessHeight, err := c.GetProcessCount()
	// If error is returned, try the request more times, then fatal.
	if err != nil {
		log.Error(err)
		for errs := 0; ; errs++ {
			if errs > 10 {
				log.Error("Gateway Disconnected")
				exit <- struct{}{}
				return
			}
			gatewayProcessHeight, err = c.GetProcessCount()
			if err == nil {
				break
			}
		}
	}
	if localProcessHeight >= gatewayProcessHeight {
		return
	}

	// Get height map for list of entities, current heights stored
	heightMap := GetHeightMap(d, config.EntityProcessCountMapKey)
	// Initialize concurrency helper variables
	heightMapMutex := new(sync.Mutex)
	requestMutex := new(sync.Mutex)
	numNewProcesses := 0
	numEntities := len(heightMap.Heights)
	complete := make(chan struct{}, numEntities)

	batch := d.NewBatch()

	for entity, localHeight := range heightMap.Heights {
		go fetchProcesses(entity, localHeight, localProcessHeight, d, batch, heightMap, heightMapMutex, requestMutex, &numNewProcesses, c, complete)
	}
	log.Debugf("Found %d stored entities", numEntities)

	// Sync: wait here for all goroutines to complete
	num := 0
	for range complete {
		if num >= numEntities-1 {
			break
		}
		num++
	}
	log.Debugf("Fetched %d new processes", numNewProcesses)

	// Write updated entity process height map
	rawHeightMap, err := proto.Marshal(heightMap)
	if err != nil {
		log.Error(err)
	}
	heightMapKey := []byte(config.EntityProcessCountMapKey)
	batch.Put(heightMapKey, rawHeightMap)
	// Write global process height
	encHeight := voctypes.Height{Height: localProcessHeight + int64(numNewProcesses)}
	rawHeight, err := proto.Marshal(&encHeight)
	if err != nil {
		log.Error(err)
	}
	heightKey := []byte(config.LatestProcessCountKey)
	batch.Put(heightKey, rawHeight)
	batch.Write()
}

func fetchProcesses(entity string, localHeight, height int64, db *dvotedb.BadgerDB, batch dvotedb.Batch, heightMap *voctypes.HeightMap, heightMapMutex, requestMutex *sync.Mutex, numNew *int, c *api.GatewayClient, complete chan struct{}) {
	defer func() {
		complete <- struct{}{}
	}()

	var lastRawProcess []byte
	rawEntity, err := hex.DecodeString(util.TrimHex(entity))
	// Get Entity|LocalHeight:ProcessHeight
	entityProcessKey := append([]byte(config.ProcessByEntityPrefix), rawEntity...)
	entityProcessKey = append(entityProcessKey, util.EncodeInt(int(localHeight-1))...)
	rawGlobalHeight, err := db.Get(entityProcessKey)
	if err != nil {
		log.Debugf("Height Key not found: %s", err.Error())
		rawGlobalHeight = []byte{}
	} else {
		var globalHeight voctypes.Height
		err = proto.Unmarshal(rawGlobalHeight, &globalHeight)
		if err != nil {
			globalHeight.Height = -1
		}
		// Get ProcessHeight:Process
		lastProcessKey := append([]byte(config.ProcessHeightPrefix), util.EncodeInt(globalHeight.GetHeight())...)
		lastRawProcess, err = db.Get(lastProcessKey)
		if err != nil {
			log.Debugf("Process Key not found: %s", err.Error())
			lastRawProcess = []byte{}
		}
	}
	var lastProcess voctypes.Process
	if len(lastRawProcess) > 0 {
		err := proto.Unmarshal(lastRawProcess, &lastProcess)
		if err != nil {
			log.Error(err)
		}
	}
	requestMutex.Lock()
	lastPID := lastProcess.ID
	log.Debugf("Getting processes from id %s", lastPID)
	if !dvoteutil.IsHexEncodedStringWithLength(lastPID, dvotetypes.ProcessIDsize) {
		lastPID = ""
	}
	newProcessList, err := c.GetProcessList(entity, lastPID)
	// If error is returned, try the request more times, then fatal.
	if err != nil {
		log.Error(err)
		for errs := 0; ; errs++ {
			if errs > 10 {
				log.Error("Gateway Disconnected")
				exit <- struct{}{}
				return
			}
			newProcessList, err = c.GetProcessList(entity, lastPID)
			if err == nil {
				break
			}
		}
	}
	requestMutex.Unlock()
	if len(newProcessList) < 1 {
		return
	}
	var process voctypes.Process
	for _, processID := range newProcessList {
		heightMapMutex.Lock()
		*numNew++
		globalHeight := int(height) + *numNew
		localHeight := heightMap.Heights[entity]
		heightMap.Heights[entity]++
		heightMapMutex.Unlock()

		process.ID = processID
		process.EntityID = entity
		process.LocalHeight = &voctypes.Height{Height: localHeight}
		rawProcess, err := proto.Marshal(&process)
		if err != nil {
			log.Error(err)
		}
		rawPID, err := hex.DecodeString(util.TrimHex(processID))
		if err != nil {
			log.Error(err)
		}

		// Write Height:Process
		processKey := append([]byte(config.ProcessHeightPrefix), util.EncodeInt(globalHeight)...)
		batch.Put(processKey, rawProcess)

		storeHeight := &voctypes.Height{Height: int64(globalHeight)}
		rawStoreHeight, err := proto.Marshal(storeHeight)
		if err != nil {
			log.Error(err)
		}
		// Write PID:Processheight
		processIDKey := append([]byte(config.ProcessIDPrefix), rawPID...)
		batch.Put(processIDKey, rawStoreHeight)

		// Write Entity|LocalHeight:ProcessHeight
		entityProcessKey := append([]byte(config.ProcessByEntityPrefix), rawEntity...)
		entityProcessKey = append(entityProcessKey, util.EncodeInt(int(localHeight))...)

		batch.Put(entityProcessKey, rawStoreHeight)
	}
}

func storeEnvelope(tx tmtypes.Tx, height *voctypes.Height, procHeightMap *voctypes.HeightMap, procHeightMapMutex *sync.Mutex, batch dvotedb.Batch) string {
	var rawTx dvotetypes.Tx
	err := json.Unmarshal(tx, &rawTx)
	if err != nil {
		log.Error(err)
	}
	if rawTx.Type == "vote" {
		globalHeight := atomic.AddInt64(&height.Height, 1)
		var voteTx dvotetypes.VoteTx
		err = json.Unmarshal(tx, &voteTx)
		if err != nil {
			log.Error(err)
		}

		// Write vote package
		votePackage := voctypes.Envelope{
			ProcessID:    voteTx.ProcessID,
			Package:      voteTx.VotePackage,
			GlobalHeight: globalHeight,
		}

		// Update height of process env belongs to
		procHeightMapMutex.Lock()
		procHeight, ok := procHeightMap.Heights[util.TrimHex(votePackage.GetProcessID())]
		if !ok {
			procHeight = 0
		}
		procHeight++
		procHeightMap.Heights[util.TrimHex(votePackage.GetProcessID())] = procHeight
		procHeightMapMutex.Unlock()

		votePackage.ProcessHeight = procHeight

		// Generate nullifier as in go-dvote vochain/transaction.go
		signature := voteTx.Signature
		voteTx.Signature = ""
		voteTx.Type = ""
		voteBytes, err := json.Marshal(&voteTx)
		if err != nil {
			log.Error(err)
		}
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

		// Write globalHeight:package
		rawEnvelope, err := proto.Marshal(&votePackage)
		if err != nil {
			log.Error(err)
		}
		packageKey := append([]byte(config.EnvPackagePrefix), util.EncodeInt(globalHeight)...)
		batch.Put(packageKey, rawEnvelope)

		// Write nullifier:globalHeight
		storeHeight := voctypes.Height{Height: globalHeight}
		rawHeight, err := proto.Marshal(&storeHeight)
		if err != nil {
			log.Error(err)
		}
		nullifier, err := hex.DecodeString(util.TrimHex(votePackage.Nullifier))
		if err != nil {
			log.Error(err)
		}
		nullifierKey := append([]byte(config.EnvNullifierPrefix), nullifier...)
		batch.Put(nullifierKey, rawHeight)

		// Write pid|heightbyPID:globalHeight
		heightBytes := util.EncodeInt(procHeight)
		PIDBytes, err := hex.DecodeString(util.TrimHex(votePackage.ProcessID))
		if err != nil {
			log.Error(err)
		}
		heightKey := append([]byte(config.EnvPIDPrefix), PIDBytes...)
		heightKey = append(heightKey, heightBytes...)
		batch.Put(heightKey, rawHeight)

		return votePackage.Nullifier
	}
	return ""
}