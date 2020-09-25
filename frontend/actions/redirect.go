package actions

// SetCurrentPage is the action to set the current page title
type SetCurrentPage struct {
	Page string
}

// SignalRedirect is the action to signal a page redirect
type SignalRedirect struct {
}
