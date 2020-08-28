package storeutil

// PageStore stores information needed to display/update a pagination element
type PageStore struct {
	Tab           string
	PagChannel    chan int
	CurrentPage   int
	DisableUpdate bool
}
