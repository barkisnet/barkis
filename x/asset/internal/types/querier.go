package types

// querier keys
const (
	DefaultQueryLimit = 100
	QueryParams = "params"
	QueryToken  = "get"
	QueryList   = "list"
)

// QueryTokensParams defines the params for the following queries:
// - 'custom/asset/list'
type QueryTokensParams struct {
	Page, Limit int
}