package types

import (
	"github.com/barkisnet/barkis/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	//todo refactor name
	cdc.RegisterConcrete(IssueMsg{}, "cosmos-sdk/IssueMsg", nil)
	cdc.RegisterConcrete(MintMsg{}, "cosmos-sdk/MintMsg", nil)
	cdc.RegisterConcrete(DelayedTransferMsg{}, "cosmos-sdk/DelayedTransferMsg", nil)
}

// module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
