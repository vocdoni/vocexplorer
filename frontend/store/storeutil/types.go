package storeutil

// PageStore stores information needed to display/update a pagination element
type PageStore struct {
	CurrentPage   int
	DisableUpdate bool
	Index         int
	PagChannel    chan int
	Search        bool
	SearchChannel chan string
	Tab           string
}
