package actions

//GatewayConnected is the action to change the connection status of the gateway
type GatewayConnected struct {
	GatewayErr error
}

// SetLoading is the action to change the loading status
type SetLoading struct {
	Loading bool
}

// SetSearchTerm is the action to set the current search term
type SetSearchTerm struct {
	SearchTerm string
}
