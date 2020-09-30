package db

import (
	"encoding/hex"
	"encoding/json"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/ptypes"
	"gitlab.com/vocdoni/go-dvote/crypto/ethereum"
	dvotedb "gitlab.com/vocdoni/go-dvote/db"
	"gitlab.com/vocdoni/go-dvote/log"
	dvoteutil "gitlab.com/vocdoni/go-dvote/util"
	"gitlab.com/vocdoni/go-dvote/vochain"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/api/dvotetypes"
	"gitlab.com/vocdoni/vocexplorer/api/rpc"
	"gitlab.com/vocdoni/vocexplorer/api/tmtypes"
	"gitlab.com/vocdoni/vocexplorer/config"
	voctypes "gitlab.com/vocdoni/vocexplorer/proto"
	"gitlab.com/vocdoni/vocexplorer/util"
	"google.golang.org/protobuf/proto"
)

func updateBlockchainInfo(d *dvotedb.BadgerDB, t *rpc.TendermintRPC, g *api.GatewayClient) {
	bc := new(voctypes.BlockchainInfo)

	status := api.GetHealth(t)
	if status == nil {
		log.Warnf("Unable to get vochain status")
	} else {
		bc.Network = status.NodeInfo.Network
		bc.Version = status.NodeInfo.Version
		rawSync, err := json.Marshal(status.SyncInfo)
		if err != nil {
			log.Warn(err)
		} else {
			bc.SyncInfo = rawSync
		}
	}

	genesis := api.GetGenesis(t)
	if genesis == nil {
		log.Warnf("Unable to get genesis block")
	} else {
		timeStamp, err := ptypes.TimestampProto(genesis.GenesisTime)
		if err != nil {
			log.Warn(err)
		} else {
			bc.GenesisTimeStamp = timeStamp
		}
		bc.ChainID = genesis.ChainID
	}

	blockTime, blockTimeStamp, height, err := g.GetBlockStatus()
	if err != nil {
		log.Warn(err)
	} else {
		bc.BlockTime = blockTime[:]
		bc.BlockTimeStamp = blockTimeStamp
		bc.Height = height
	}

	rawBlockchainInfo, err := proto.Marshal(bc)
	if err != nil {
		log.Warn(err)
	} else {
		d.Put([]byte(config.BlockchainInfoKey), rawBlockchainInfo)
	}
}

func updateBlockList(d *dvotedb.BadgerDB, t *rpc.TendermintRPC) {
	// Fetch latest block & tx heights
	latestBlockHeight := GetHeight(d, config.LatestBlockHeightKey, 1)
	latestTxHeight := GetHeight(d, config.LatestTxHeightKey, 1)
	latestEnvelopeCount := GetHeight(d, config.LatestEnvelopeCountKey, 0)
	maxBlockTxs := GetInt64(d, config.MaxTxsPerBlockKey)
	maxMinuteTxs := GetInt64(d, config.MaxTxsPerMinuteKey)
	largestBlock := ""
	largestBlockHeight := int64(0)

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
		for errs := 0; ; errs++ {
			if errs > 2 {
				log.Errorf("Unable to get status: %s", err.Error())
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
	wg := new(sync.WaitGroup)
	// nextHeight and myHeight channels synchronize goroutines before fetching validator block height, so blocks by validator are ordered by block height
	nextHeight := make(chan struct{})
	myHeight := make(chan struct{})
	txsByMinute := make(map[int64]int64)
	for ; int(i) < numNewBlocks; i++ {
		wg.Add(1)
		go fetchBlock(i+latestBlockHeight.GetHeight(), &maxBlockTxs, &largestBlockHeight, &largestBlock, &batch, t, wg, myHeight, nextHeight, &txsList[i], valMap, valMapMutex, txsByMinute)
		if i == 0 {
			//Signal to the first block to start
			close(myHeight)
		}
		myHeight = nextHeight
		nextHeight = make(chan struct{})
	}

	if i > 0 {
		// Sync: wait here for all goroutines to complete
		wg.Wait()
		log.Debugf("Setting block %d ", latestBlockHeight.GetHeight()+i)

		wg := new(sync.WaitGroup)
		for _, txs := range txsList {
			if len(txs) > 0 {
				wg.Add(1)
				go updateTxs(latestTxHeight.GetHeight(), txs, t, batch, wg, latestEnvelopeCount, procEnvHeightMap, procEnvHeightMapMutex)
				latestTxHeight.Height += int64(len(txs))
			}
		}

		// Sync: wait here for all goroutines to complete
		wg.Wait()

		// Update the max txs per minute
		var maxTxMinute int64
		for t, n := range txsByMinute {
			if n > maxMinuteTxs {
				maxMinuteTxs = n
				maxTxMinute = t
			}
		}
		if maxTxMinute != 0 {
			batch.Put([]byte(config.MaxTxsPerMinuteKey), util.EncodeInt(maxMinuteTxs))
			batch.Put([]byte(config.MaxTxsMinuteID), util.EncodeInt(maxTxMinute))
		}

		// write largestblock
		if largestBlock != "" {
			rawLargestBlock, err := hex.DecodeString(largestBlock)
			if err != nil {
				log.Warn(err)
			}
			rawLargestBlockHeight := util.EncodeInt(largestBlockHeight)
			batch.Put([]byte(config.MaxTxsBlockIDKey), rawLargestBlock)
			batch.Put([]byte(config.MaxTxsBlockHeightKey), rawLargestBlockHeight)
		}
		// write max txs per block
		batch.Put([]byte(config.MaxTxsPerBlockKey), util.EncodeInt(maxBlockTxs))

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

func fetchBlock(height int64, maxBlockTxs, largestBlockHeight *int64, largestBlock *string, batch *dvotedb.Batch, t *rpc.TendermintRPC, wg *sync.WaitGroup, myHeight, nextHeight chan struct{}, txs *tmtypes.Txs, valMap *voctypes.HeightMap, valMapMutex *sync.Mutex, maxMinuteTxs map[int64]int64) {
	// Signal
	defer wg.Done()
	// Thread-safe api request
	res, err := t.Block(&height)
	// If error is returned, try the request more times, then exit
	if err != nil {
		if strings.Contains(err.Error(), "closed") {
			exit <- struct{}{}
			return
		}
		for errs := 0; ; errs++ {
			if errs > 10 {
				log.Errorf("Unable to get block: %s", err.Error())
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
	// If this block has the most txs, set the maxBlockTxs
	if block.NumTxs > *maxBlockTxs {
		*maxBlockTxs = block.NumTxs
		*largestBlock = hex.EncodeToString(block.GetHash())
		*largestBlockHeight = block.GetHeight()
	}

	// Add numTxs to this minute's total txs
	prev, ok := maxMinuteTxs[(block.Time.GetSeconds()/60)*60]
	if !ok {
		prev = 0
	}
	maxMinuteTxs[(block.Time.GetSeconds()/60)*60] = prev + block.NumTxs

	height, ok = valMap.Heights[util.HexToString(block.Proposer)]
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
		if strings.Contains(err.Error(), "closed") {
			exit <- struct{}{}
			return
		}
		for errs := 0; ; errs++ {
			if errs > 2 {
				log.Errorf("Unable to get validators: %s", err.Error())
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
			if strings.Contains(err.Error(), "closed") {
				exit <- struct{}{}
				return
			}
			for errs := 0; ; errs++ {
				if errs > 2 {
					log.Errorf("Unable to get validators: %s", err.Error())
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

func updateTxs(startTxHeight int64, txs tmtypes.Txs, t *rpc.TendermintRPC, batch dvotedb.Batch, wg *sync.WaitGroup, envHeight *voctypes.Height, procHeightMap *voctypes.HeightMap, procHeightMapMutex *sync.Mutex) {
	defer wg.Done()
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
		txStore.Nullifier = storeEnvelope(txStore.Tx, txStore.TxHeight, envHeight, procHeightMap, procHeightMapMutex, batch)
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
}

func updateEntityList(d *dvotedb.BadgerDB, c *api.GatewayClient) {
	localEntityHeight := GetHeight(d, config.LatestEntityCountKey, 0).GetHeight()
	gatewayEntityHeight, err := c.GetEntityCount()
	// If error is returned, try the request more times, then fatal.
	if err != nil {
		if strings.Contains(err.Error(), "closed") {
			exit <- struct{}{}
			return
		}
		for errs := 0; ; errs++ {
			if errs > 2 {
				log.Errorf("Unable to get entity height: %s", err.Error())
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
	if err != nil {
		log.Warn(err)
	}
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
		entity = util.TrimHex(entity)
		rawEntity, err := hex.DecodeString(entity)
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
		if strings.Contains(err.Error(), "closed") {
			exit <- struct{}{}
			return
		}
		for errs := 0; ; errs++ {
			if errs > 2 {
				log.Errorf("Unable to get process height: %s", err.Error())
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
	wg := new(sync.WaitGroup)

	batch := d.NewBatch()

	for entity, localHeight := range heightMap.Heights {
		wg.Add(1)
		go fetchProcesses(entity, localHeight, localProcessHeight, d, batch, heightMap, heightMapMutex, requestMutex, &numNewProcesses, c, wg)
	}
	log.Debugf("Found %d stored entities", numEntities)

	// Sync: wait here for all goroutines to complete
	wg.Wait()
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

func fetchProcesses(entity string, localHeight, height int64, db *dvotedb.BadgerDB, batch dvotedb.Batch, heightMap *voctypes.HeightMap, heightMapMutex, requestMutex *sync.Mutex, numNew *int, c *api.GatewayClient, wg *sync.WaitGroup) {
	defer wg.Done()

	var lastRawProcess []byte
	rawEntity, err := hex.DecodeString(util.TrimHex(entity))
	if err != nil {
		log.Warn(err)
	}
	// Get Entity|LocalHeight:ProcessHeight
	entityProcessKey := append([]byte(config.ProcessByEntityPrefix), rawEntity...)
	entityProcessKey = append(entityProcessKey, util.EncodeInt(int(localHeight-1))...)
	rawGlobalHeight, err := db.Get(entityProcessKey)
	if err != nil {
		log.Debugf("Height Key not found: %s", err.Error())
		rawGlobalHeight = []byte{}
	}
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
	var lastProcess voctypes.Process
	if len(lastRawProcess) > 0 {
		err := proto.Unmarshal(lastRawProcess, &lastProcess)
		if err != nil {
			log.Error(err)
		}
	}
	requestMutex.Lock()
	lastPID := lastProcess.ID
	entity = strings.ToLower(util.TrimHex(entity))
	log.Debugf("Getting processes from id %s, entity %s", lastPID, entity)
	if !dvoteutil.IsHexEncodedStringWithLength(lastPID, dvotetypes.ProcessIDsize) {
		lastPID = ""
	}
	newProcessList, err := c.GetProcessList(entity, lastPID)
	// If error is returned, try the request more times, then fatal.
	if err != nil {
		if strings.Contains(err.Error(), "closed") {
			exit <- struct{}{}
			requestMutex.Unlock()
			return
		}
		for errs := 0; ; errs++ {
			if errs > 2 {
				log.Errorf("Unable to get process list: %s", err.Error())
				requestMutex.Unlock()
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

func storeEnvelope(tx tmtypes.Tx, txHeight int64, height *voctypes.Height, procHeightMap *voctypes.HeightMap, procHeightMapMutex *sync.Mutex, batch dvotedb.Batch) string {
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
			GlobalHeight: globalHeight,
			Package:      voteTx.VotePackage,
			ProcessID:    voteTx.ProcessID,
			TxHeight:     txHeight,
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
