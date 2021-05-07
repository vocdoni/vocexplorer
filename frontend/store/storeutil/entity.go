package storeutil

// Entities stores the current entities information
type Entities struct {
	Count             int
	CurrentEntity     Entity
	CurrentEntityID   string
	EntityIDs         []string
	Pagination        PageStore
	ProcessPagination PageStore
	ProcessHeights    map[string]int64
}

// Entity holds info about one vochain entity
type Entity struct {
	ProcessCount int
	ProcessIds   []string
}
