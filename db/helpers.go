package db

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"strings"

	dvotedb "gitlab.com/vocdoni/go-dvote/db"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	voctypes "gitlab.com/vocdoni/vocexplorer/proto"
	"gitlab.com/vocdoni/vocexplorer/util"
	"google.golang.org/protobuf/proto"
)

// GetInt64 fetches a int64 value from the database corresponding to given key
func GetInt64(d *dvotedb.BadgerDB, key string) int64 {
	var val int64
	var num int
	has, err := d.Has([]byte(key))
	if err != nil {
		log.Error(err)
	}
	if has {
		raw, err := d.Get([]byte(key))
		if err != nil {
			log.Error(err)
		}
		val, num = binary.Varint(raw)
		if num < 1 {
			log.Error("Could not decode int")
		}
	}
	return val
}

// GetHeight fetches a height value from the database corresponding to given key
func GetHeight(d *dvotedb.BadgerDB, key string, def int64) *voctypes.Height {
	height := &voctypes.Height{}
	has, err := d.Has([]byte(key))
	if err != nil {
		log.Error(err)
	}
	if has {
		val, err := d.Get([]byte(key))
		if err != nil {
			log.Error(err)
		}
		err = proto.Unmarshal(val, height)
		if err != nil {
			log.Error(err)
		}
	}
	if def > height.GetHeight() {
		height.Height = def
	}
	return height
}

// GetHeightMap fetches a height map from the database
func GetHeightMap(d *dvotedb.BadgerDB, key string) *voctypes.HeightMap {
	var valMap voctypes.HeightMap
	valMapKey := []byte(key)
	has, err := d.Has(valMapKey)
	if err != nil {
		log.Error(err)
	}
	if has {
		rawValMap, err := d.Get(valMapKey)
		if err != nil {
			log.Error(err)
		}
		proto.Unmarshal(rawValMap, &valMap)
	}
	if valMap.Heights == nil {
		valMap.Heights = make(map[string]int64)
	}
	return &valMap
}

// ListItemsByHeight returns a list of items given integer keys
func ListItemsByHeight(d *dvotedb.BadgerDB, max, height int, prefix []byte) [][]byte {
	if max > 64 {
		max = 64
	}
	var hashList [][]byte
	for ; max > 0 && height >= 0; max-- {
		heightKey := util.EncodeInt(height)
		key := append(prefix, heightKey...)
		has, err := d.Has(key)
		if !has || err != nil {
			if err != nil {
				log.Error(err)
			}
			height--
			continue
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

// SearchItems returns a list of items given search term, starting with given prefix
func SearchItems(d *dvotedb.BadgerDB, max int, term string, prefix []byte) [][]byte {
	return searchIter(d, max, term, prefix, false)
}

// SearchKeys returns a list of key values including the search term, starting with the given prefix
func SearchKeys(d *dvotedb.BadgerDB, max int, term string, prefix []byte) [][]byte {
	return searchIter(d, max, term, prefix, true)
}

// SearchBlocksByValidator returns a list of blocks given the search term and validator
func SearchBlocksByValidator(d *dvotedb.BadgerDB, max int, term, validator string) [][]byte {
	rawValidator, err := hex.DecodeString(validator)
	if err != nil {
		log.Warn(err)
		return nil
	}
	prefix := []byte(config.BlockHashPrefix)
	if max > 64 {
		max = 64
	}

	var itemList [][]byte
	iter := d.NewIterator().(*dvotedb.BadgerIterator)
	// dvote badgerdb iterator has bug: rewind() on first Seek call rewinds to start of db
	var valid bool
	iter.Seek(prefix)
	valid = iter.Next()
	for iter.Seek(prefix); valid && bytes.HasPrefix(iter.Key(), prefix); valid = iter.Next() {
		if max < 1 {
			break
		}
		// Converting each key to string is a Non-ideal solution. Without converting to string, we cannot analyze each character because each hex byte represents two characters. Alternative would be to decode hex byte array to byte array (one byte per character), but this may be no faster.
		keyString := hex.EncodeToString(iter.Key())
		if strings.Contains(keyString, term) {
			block := &voctypes.StoreBlock{}
			err := proto.Unmarshal(iter.Value(), block)
			if err != nil {
				log.Warn(err)
				continue
			}
			// Check if block found belongs to validator
			if bytes.Equal(block.GetProposer(), rawValidator) {
				itemList = append(itemList, iter.Value())
				max--
			}
		}
	}
	iter.Release()
	return itemList
}

func searchIter(d *dvotedb.BadgerDB, max int, term string, prefix []byte, getKey bool) [][]byte {
	if max > 64 {
		max = 64
	}

	var itemList [][]byte
	iter := d.NewIterator().(*dvotedb.BadgerIterator)
	// dvote badgerdb iterator has bug: rewind() on first Seek call rewinds to start of db
	var valid bool
	iter.Seek(prefix)
	valid = iter.Next()
	for iter.Seek(prefix); valid && bytes.HasPrefix(iter.Key(), prefix); valid = iter.Next() {
		if max < 1 {
			break
		}
		// Converting each key to string is a Non-ideal solution. Without converting to string, we cannot analyze each character because each hex byte represents two characters. Alternative would be to decode hex byte array to byte array (one byte per character), but this may be no faster.
		keyString := hex.EncodeToString(iter.Key())
		if strings.Contains(keyString, term) {
			if getKey {
				// Append key, cutting off the prefix bytes
				// Safe-copy of key
				itemList = append(itemList, append([]byte{}, iter.Key()[len(prefix):]...))
			} else {
				itemList = append(itemList, iter.Value())
			}
			max--
		}
	}
	iter.Release()
	return itemList
}
