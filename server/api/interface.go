package api

// Connector Is for interface connector
type Connector interface {
	Connect(*Client) error
	Verify(*Client) error
}
