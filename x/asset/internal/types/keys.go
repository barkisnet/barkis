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
	TokenKeyPrefix                  = []byte{0x01}
	SequenceKey                     = []byte{0x02}
	DelayedTransferPrefix           = []byte{0x03}
	DelayedTransferMatureTimePrefix = []byte{0x04}
	DelayedTransferFromPrefix       = []byte{0x05}
	DelayedTransferToPrefix         = []byte{0x06}
)

func BuildTokenKey(symbol string) []byte {
	return append(TokenKeyPrefix, []byte(symbol)...)
}

func BuildDelayedTransferKey(sequence int64) []byte {
	sequenceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(sequenceBytes, uint64(sequence))

	return append(DelayedTransferPrefix, sequenceBytes...)
}

func BuildDelayedTransferMatureTimeKey(matureTime int64, sequence int64) []byte {
	matureTimeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(matureTimeBytes, uint64(matureTime))

	sequenceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(sequenceBytes, uint64(sequence))

	return append(append(DelayedTransferMatureTimePrefix, matureTimeBytes...), sequenceBytes...)
}

func BuildDelayedTransferFromKey(addr sdk.AccAddress, sequence int64) []byte {
	sequenceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(sequenceBytes, uint64(sequence))
	key := append(DelayedTransferFromPrefix, []byte(addr)...)
	return append(key, sequenceBytes...)
}

func BuildDelayedTransferToKey(addr sdk.AccAddress, sequence int64) []byte {
	sequenceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(sequenceBytes, uint64(sequence))
	key := append(DelayedTransferToPrefix, []byte(addr)...)
	return append(key, sequenceBytes...)
}
