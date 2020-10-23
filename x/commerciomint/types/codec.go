package types

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgMintCCC{}, "commercio/MsgMintCCC", nil)
	cdc.RegisterConcrete(MsgBurnCCC{}, "commercio/MsgBurnCCC", nil)
	cdc.RegisterConcrete(MsgSetCCCConversionRate{}, "commercio/MsgSetCCCConversionRate", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
