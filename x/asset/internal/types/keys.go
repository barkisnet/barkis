package types

const (
	// module name
	ModuleName = "asset"

	// StoreKey is the store key string for asset
	StoreKey = ModuleName

	// RouterKey is the message route for asset
	RouterKey = ModuleName

	// QuerierRoute is the querier route for asset
	QuerierRoute = ModuleName
)

var (
	TokenKeyPrefix = []byte{0x01}

	ParamStoreKeyMaxDecimal = []byte("MaxDecimal")
)

func BuildTokenKey(symbol string) []byte {
	return append(TokenKeyPrefix, []byte(symbol)...)
}


