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
	"gitlab.com/vocdoni/go-dvote/log"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	dvoteutil "gitlab.com/vocdoni/go-dvote/util"
	"gitlab.com/vocdoni/go-dvote/vochain"
	"gitlab.com/vocdoni/vocexplorer/config"
	voctypes "gitlab.com/vocdoni/vocexplorer/proto"
	"gitlab.com/vocdoni/vocexplorer/util"
	"google.golang.org/protobuf/proto"
)

func (d *ExplorerDB) updateBlockchainInfo() {
	bc := d.Vs.GetStatus()
	if bc == nil {
		log.Warnf("Unable to get vochain status")
		return
	}
	rawBlockchainInfo, err := proto.Marshal(bc)
	if err != nil {
		log.Warn(err)
		return
	}
	d.Db.Put([]byte(config.BlockchainInfoKey), rawBlockchainInfo)
}

func (d *ExplorerDB) updateBlockList() {
	state := BlockState{
		batch:                    d.Db.NewBatch(),
		blockHeight:              GetHeight(d.Db, config.LatestBlockHeightKey, 1).GetHeight(),
		envelopeHeight:           GetHeight(d.Db, config.LatestEnvelopeCountKey, 0).GetHeight(),
		fault:                    int32(0),
		largestBlock:             int64(0),
		largestBlockHash:         "",
		maxBlockTxs:              GetInt64(d.Db, config.MaxTxsPerBlockKey),
		maxMinuteTxs:             GetInt64(d.Db, config.MaxTxsPerMinuteKey),
		processEnvelopeHeightMap: GetHeightMap(d.Db, config.ProcessEnvelopeCountMapKey),
		validatorBlockHeightMap:  GetHeightMap(d.Db, config.ValidatorHeightMapKey),
		stateMutex:               new(sync.Mutex),
		txHeight:                 GetHeight(d.Db, config.LatestTxHeightKey, 1).GetHeight(),
		txsByMinute:              make(map[int64]int64),
	}

	status := d.Vs.GetStatus()
	gwBlockHeight := status.LatestBlockHeight

	// Wait for new blocks to be available
	for gwBlockHeight-state.blockHeight < 2 {
		time.Sleep(500 * time.Millisecond)
		return
	}

	i := int64(0)
	numNewBlocks := util.Min(config.NumBlockUpdates, int(gwBlockHeight-state.blockHeight-1))
	wg := new(sync.WaitGroup)
	// nextHeight and myHeight channels synchronize goroutines before fetching validator block height, so blocks by validator are ordered by block height
	nextHeight := make(chan struct{})
	myHeight := make(chan struct{})
	fault := int32(0)

	for ; int(i) < numNewBlocks; i++ {
		wg.Add(1)
		go d.fetchBlock(i+state.blockHeight, wg, myHeight, nextHeight, &state)
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
		if fault == 1 {
			log.Debugf("Could not fetch blocks")
			return
		}
		log.Infof("Setting block %d ", state.blockHeight+i)

		// Update the max txs per minute
		var maxTxMinute int64
		for t, n := range state.txsByMinute {
			if n > state.maxMinuteTxs {
				state.maxMinuteTxs = n
				maxTxMinute = t
			}
		}
		if maxTxMinute != 0 {
			state.batch.Put([]byte(config.MaxTxsPerMinuteKey), util.EncodeInt(state.maxMinuteTxs))
			state.batch.Put([]byte(config.MaxTxsMinuteID), util.EncodeInt(maxTxMinute))
		}

		// write largestblock
		if state.largestBlockHash != "" {
			rawLargestBlock, err := hex.DecodeString(state.largestBlockHash)
			if err != nil {
				log.Warn(err)
			}
			rawLargestBlockHeight := util.EncodeInt(state.largestBlock)
			state.batch.Put([]byte(config.MaxTxsBlockIDKey), rawLargestBlock)
			state.batch.Put([]byte(config.MaxTxsBlockHeightKey), rawLargestBlockHeight)
		}
		// write max txs per block
		state.batch.Put([]byte(config.MaxTxsPerBlockKey), util.EncodeInt(state.maxBlockTxs))

		rawValMap, err := proto.Marshal(state.validatorBlockHeightMap)
		if err != nil {
			log.Error(err)
		}
		state.batch.Put([]byte(config.ValidatorHeightMapKey), rawValMap)
		rawProcMap, err := proto.Marshal(state.processEnvelopeHeightMap)
		if err != nil {
			log.Error(err)
		}
		state.batch.Put([]byte(config.ProcessEnvelopeCountMapKey), rawProcMap)
		blockHeight := voctypes.Height{Height: state.blockHeight + i}
		encBlockHeight, err := proto.Marshal(&blockHeight)
		if err != nil {
			log.Error(err)
		}
		encTxHeight, err := proto.Marshal(&voctypes.Height{Height: state.txHeight})
		if err != nil {
			log.Error(err)
		}
		encEnvCount, err := proto.Marshal(&voctypes.Height{Height: state.envelopeHeight})
		if err != nil {
			log.Error(err)
		}

		state.batch.Put([]byte(config.LatestTxHeightKey), encTxHeight)

		state.batch.Put([]byte(config.LatestBlockHeightKey), encBlockHeight)
		state.batch.Put([]byte(config.LatestEnvelopeCountKey), encEnvCount)
		if err := state.batch.Write(); err != nil {
			log.Error(err)
		}
	}

}

func (d *ExplorerDB) fetchBlock(height int64, wg *sync.WaitGroup, myHeight, nextHeight chan struct{}, state *BlockState) {
	// Signal
	defer wg.Done()
	if atomic.LoadInt32(&state.fault) != 0 {
		close(nextHeight)
		return
	}
	res, err := d.Vs.GetBlock(height)
	if err != nil {
		log.Warn(err)
		atomic.StoreInt32(&state.fault, 1)
		close(nextHeight)
		return
	}
	var block voctypes.StoreBlock
	block.Hash = res.BlockID.Hash
	block.Height = res.Block.Header.Height
	block.Proposer = res.Block.ProposerAddress
	tm, err := ptypes.TimestampProto(res.Block.Header.Time)
	if err != nil {
		log.Error(err)
	}
	block.Time = tm

	transactions, err := d.Vs.GetTransactions(block.Height)
	if err != nil {
		log.Error(err)
	}
	block.NumTxs = int64(len(transactions))
	bodyValue, err := proto.Marshal(&block)
	if err != nil {
		log.Error(err)
	}

	// Wait for myHeight channel to close, this means fetchBlock for previous block has been assigned a validator block height
	<-myHeight
	if atomic.LoadInt32(&state.fault) != 0 {
		close(nextHeight)
		return
	}
	// Update height of validator block belongs to
	state.stateMutex.Lock()
	// If this block has the most txs, set the maxBlockTxs
	if block.NumTxs > state.maxBlockTxs {
		state.maxBlockTxs = block.NumTxs
		state.largestBlockHash = hex.EncodeToString(block.GetHash())
		state.largestBlock = block.GetHeight()
	}
	// Add numTxs to this minute's total txs
	prev, ok := state.txsByMinute[(block.Time.GetSeconds()/60)*60]
	if !ok {
		prev = 0
	}
	state.txsByMinute[(block.Time.GetSeconds()/60)*60] = prev + block.NumTxs
	height, ok = state.validatorBlockHeightMap.Heights[util.HexToString(block.Proposer)]
	if !ok {
		height = 0
	}
	height++
	state.validatorBlockHeightMap.Heights[util.HexToString(block.Proposer)] = height

	d.logTxs(transactions, state)
	state.stateMutex.Unlock()
	// Signal to next block that I have been assigned a validator block height
	close(nextHeight)

	blockHeightKey := append([]byte(config.BlockHeightPrefix), util.EncodeInt(block.GetHeight())...)
	blockHashKey := append([]byte(config.BlockHashPrefix), block.GetHash()...)
	validatorHeightKey := append([]byte(config.BlockByValidatorPrefix), block.GetProposer()...)
	validatorHeightKey = append(validatorHeightKey, util.EncodeInt(height)...)
	hashValue := block.GetHash()

	// Thread-safe batch operations
	// Store hash:body
	state.batch.Put(blockHashKey, bodyValue)
	// Store globalheight:hash
	state.batch.Put(blockHeightKey, hashValue)
	// Store validator|heightbyValidator:hash
	state.batch.Put(validatorHeightKey, hashValue)
}

func (d *ExplorerDB) updateValidatorList() {
	latestBlockHeight := GetHeight(d.Db, config.LatestBlockHeightKey, 0)
	latestValidatorCount := GetHeight(d.Db, config.LatestValidatorCountKey, 0)
	if latestBlockHeight.GetHeight() > 0 {
		d.fetchValidators(latestBlockHeight.GetHeight(), latestValidatorCount.GetHeight())

	}
}

func (d *ExplorerDB) fetchValidators(blockHeight, validatorCount int64) {
	resultValidators, err := d.Vs.GetValidators()
	if err != nil {
		log.Error(err)
		return
	}
	if len(resultValidators) < int(validatorCount) {
		return
	}
	batch := d.Db.NewBatch()
	// Cast each validator as storage struct, marshal, write to batch
	for i, validator := range resultValidators {
		if i < int(validatorCount) {
			continue
		}
		validatorCount++
		var storeValidator voctypes.Validator
		storeValidator.Address = validator.Address
		storeValidator.Height = &voctypes.Height{Height: validatorCount}
		// storeValidator.ProposerPriority = TODO
		storeValidator.VotingPower = validator.Power
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
	if err := batch.Write(); err != nil {
		log.Error(err)
	}
	log.Debugf("Retrieved %d validators at block height %d", len(resultValidators), blockHeight)
}

// does not happen concurrently, waits in queue once blocks + txs have been retrieved concurrently
func (d *ExplorerDB) logTxs(txs []*voctypes.Transaction, state *BlockState) {
	numTxs := int64(-1)
	var blockHeight int64
	for i, tx := range txs {
		numTxs = int64(i + 1)
		txHashKey := append([]byte(config.TxHashPrefix), tx.Hash...)
		tx.TxHeight = state.txHeight
		// If voteTx, get envelope nullifier
		tx.Nullifier = d.storeEnvelope(tx, state)
		txVal, err := proto.Marshal(tx)
		if err != nil {
			log.Error(err)
		}
		state.batch.Put(txHashKey, txVal)
		//Write height:tx hash
		txHeightKey := append([]byte(config.TxHeightPrefix), util.EncodeInt(tx.GetTxHeight())...)
		state.batch.Put(txHeightKey, tx.Hash)
		if i == 0 {
			blockHeight = tx.Height
		}
		state.txHeight++
	}
	if numTxs > 0 {
		log.Debugf("%d transactions logged at block %d, height %d", numTxs, blockHeight, state.txHeight)
	}
}

func (d *ExplorerDB) updateEntityList() {
	localEntityHeight := GetHeight(d.Db, config.LatestEntityCountKey, 0).GetHeight()
	globalEntityHeight := d.Vs.GetEntityCount()
	if localEntityHeight >= globalEntityHeight {
		return
	}
	latestKey := append([]byte(config.EntityHeightPrefix), util.EncodeInt(int(localEntityHeight-1))...)
	latestEntity, err := d.Db.Get(latestKey)
	if err != nil {
		latestEntity = []byte{}
	}
	log.Debugf("Getting entities from id %s", util.HexToString(latestEntity))
	newEntities, err := d.Vs.GetScrutinizerEntities(strings.ToLower(util.HexToString(latestEntity)), 100)
	if err != nil {
		log.Warn(err)
	}
	if len(newEntities) < 1 {
		log.Warn("No new entities retrieved")
		return
	}
	heightMap := GetHeightMap(d.Db, config.EntityProcessCountMapKey)

	// write new entities to db
	batch := d.Db.NewBatch()
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
			log.Warn("Retrieved entity already stored")
		}
		heightMap.Heights[entity] = 0
	}

	rawValMap, err := proto.Marshal(heightMap)
	if err != nil {
		log.Error(err)
	}
	log.Debugf("Retrieved %d new entities at height %d", len(newEntities), int(localEntityHeight)+i+1)

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

func (d *ExplorerDB) updateProcessList() {
	localProcessHeight := GetHeight(d.Db, config.LatestProcessCountKey, 0).GetHeight()
	globalProcessHeight := d.Vs.GetProcessCount()
	if localProcessHeight >= globalProcessHeight {
		return
	}

	// Get height map for list of entities, current heights stored
	heightMap := GetHeightMap(d.Db, config.EntityProcessCountMapKey)
	// Initialize concurrency helper variables
	heightMapMutex := new(sync.Mutex)
	requestMutex := new(sync.Mutex)
	numNewProcesses := 0
	numEntities := len(heightMap.Heights)
	wg := new(sync.WaitGroup)

	for entity, localHeight := range heightMap.Heights {
		wg.Add(1)
		go d.fetchProcesses(entity, localHeight, localProcessHeight, heightMap, heightMapMutex, requestMutex, &numNewProcesses, wg)
	}
	log.Debugf("Found %d stored entities", numEntities)

	// Sync: wait here for all goroutines to complete
	wg.Wait()
	log.Debugf("Retrieved %d new processes", numNewProcesses)

	batch := d.Db.NewBatch()
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

func (d *ExplorerDB) fetchProcesses(entity string, localHeight, height int64, heightMap *voctypes.HeightMap, heightMapMutex, requestMutex *sync.Mutex, numNew *int, wg *sync.WaitGroup) {
	defer wg.Done()
	batch := d.Db.NewBatch()

	var lastRawProcess []byte
	rawEntity, err := hex.DecodeString(util.TrimHex(entity))
	if err != nil {
		log.Warn(err)
	}
	// Get Entity|LocalHeight:ProcessHeight
	entityProcessKey := append([]byte(config.ProcessByEntityPrefix), rawEntity...)
	entityProcessKey = append(entityProcessKey, util.EncodeInt(int(localHeight-1))...)
	rawGlobalHeight, err := d.Db.Get(entityProcessKey)
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
	lastRawProcess, err = d.Db.Get(lastProcessKey)
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
	newProcessList, err := d.Vs.GetProcessList(entity, lastPID, 100)
	if err != nil {
		log.Error(err)
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
	batch.Write()
}

func (d *ExplorerDB) storeEnvelope(tx *voctypes.Transaction, state *BlockState) string {
	var rawTx dvotetypes.Tx
	err := json.Unmarshal(tx.Tx, &rawTx)
	if err != nil {
		log.Error(err)
	}
	if rawTx.Type == "vote" {
		state.envelopeHeight++
		var voteTx dvotetypes.VoteTx
		err = json.Unmarshal(tx.Tx, &voteTx)
		if err != nil {
			log.Error(err)
		}

		// Write vote package
		votePackage := voctypes.Envelope{
			GlobalHeight: state.envelopeHeight,
			Package:      voteTx.VotePackage,
			ProcessID:    voteTx.ProcessID,
			TxHeight:     tx.TxHeight,
		}

		// Update height of process env belongs to
		procHeight, ok := state.processEnvelopeHeightMap.Heights[util.TrimHex(votePackage.GetProcessID())]
		if !ok {
			procHeight = 0
		}
		procHeight++
		state.processEnvelopeHeightMap.Heights[util.TrimHex(votePackage.GetProcessID())] = procHeight

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
		rawPID, err := hex.DecodeString(votePackage.ProcessID)
		if err != nil {
			log.Errorf("cannot generate nullifier: (%s)", err)
		}
		votePackage.Nullifier = hex.EncodeToString(vochain.GenerateNullifier(addr, rawPID))
		for _, index := range voteTx.EncryptionKeyIndexes {
			votePackage.EncryptionKeyIndexes = append(votePackage.EncryptionKeyIndexes, int32(index))
		}

		// Write globalHeight:package
		rawEnvelope, err := proto.Marshal(&votePackage)
		if err != nil {
			log.Error(err)
		}
		packageKey := append([]byte(config.EnvPackagePrefix), util.EncodeInt(state.envelopeHeight)...)
		state.batch.Put(packageKey, rawEnvelope)

		// Write nullifier:globalHeight
		storeHeight := voctypes.Height{Height: state.envelopeHeight}
		rawHeight, err := proto.Marshal(&storeHeight)
		if err != nil {
			log.Error(err)
		}
		nullifier, err := hex.DecodeString(util.TrimHex(votePackage.Nullifier))
		if err != nil {
			log.Error(err)
		}
		nullifierKey := append([]byte(config.EnvNullifierPrefix), nullifier...)
		state.batch.Put(nullifierKey, rawHeight)

		// Write pid|heightbyPID:globalHeight
		heightBytes := util.EncodeInt(procHeight)
		PIDBytes, err := hex.DecodeString(util.TrimHex(votePackage.ProcessID))
		if err != nil {
			log.Error(err)
		}
		heightKey := append([]byte(config.EnvPIDPrefix), PIDBytes...)
		heightKey = append(heightKey, heightBytes...)
		state.batch.Put(heightKey, rawHeight)

		return votePackage.Nullifier
	}
	return ""
}
