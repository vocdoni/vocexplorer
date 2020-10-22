package router

import (
	"encoding/json"
	"net/http"

	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/db"
	ptypes "gitlab.com/vocdoni/vocexplorer/proto"
	"google.golang.org/protobuf/proto"
)

// GetValidatorHandler writes the validator corresponding to given address key
func GetValidatorHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildItemByIDHandler(d, "id", config.ValidatorPrefix, nil, packValidator)
}

// ListValidatorsHandler writes a list of validators from 'from'
func ListValidatorsHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsHandler(d,
		config.ValidatorHeightPrefix,
		func(key []byte) ([]byte, error) {
			return d.Db.Get(append([]byte(config.ValidatorPrefix), key...))
		},
		packValidator,
	)
}

// SearchValidatorsHandler writes a list of validators by search term
func SearchValidatorsHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildSearchHandler(d,
		config.ValidatorPrefix,
		false,
		nil,
		packValidator,
	)
}

func packValidator(raw []byte) []byte {
	var item ptypes.Validator
	err := proto.Unmarshal(raw, &item)
	if err != nil {
		log.Error(err)
	}
	new, err := json.Marshal(item.Mirror())
	if err != nil {
		log.Error(err)
	}
	return new
}
