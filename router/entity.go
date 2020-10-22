package router

import (
	"net/http"

	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/db"
)

// ListEntitiesHandler writes a list of entities from 'from'
func ListEntitiesHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsHandler(d, config.EntityHeightPrefix, nil, nil)
}

// SearchEntitiesHandler writes a list of entities by search term
func SearchEntitiesHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildSearchHandler(d,
		config.EntityIDPrefix,
		true,
		nil,
		nil,
	)
}
