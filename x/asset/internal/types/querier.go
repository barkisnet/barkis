package types

// querier keys
const (
	DefaultQueryLimit = 100
	QueryParams       = "params"
	GetToken          = "get"
	ListToken         = "list"

	GetDelayedTranfer      = "get_delayed_transfer"
	ListDelayedTranfer     = "list_delayed_transfer"
	ListDelayedTranferFrom = "list_delayed_transfer_from"
	ListDelayedTranferTo   = "list_delayed_transfer_to"
)

// QueryTokensParams defines the params for the following queries:
// - 'custom/asset/list'
type QueryTokensParams struct {
	Page, Limit int
}

// QueryDelayedTranferParams defines the params for the following queries:
// - 'custom/asset/list_delayed_transfer'
type QueryDelayedTranferParams struct {
	Page, Limit int
}
