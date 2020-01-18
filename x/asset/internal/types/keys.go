package types

import (
	"encoding/binary"

	sdk "github.com/barkisnet/barkis/types"
)

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
	TokenKeyPrefix      = []byte{0x01}
	DelayTransferPrefix = []byte{0x02}
)

func BuildTokenKey(symbol string) []byte {
	return append(TokenKeyPrefix, []byte(symbol)...)
}

func BuildDelayTransferKey(addr sdk.AccAddress, sequence int64) []byte {
	sequenceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(sequenceBytes, uint64(sequence))
	key := append(DelayTransferPrefix, []byte(addr)...)
	return append(key, sequenceBytes...)
}
