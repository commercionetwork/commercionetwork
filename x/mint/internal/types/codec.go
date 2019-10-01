package types

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgOpenCDP{}, "commercio/OpenCDP", nil)
	cdc.RegisterConcrete(MsgCloseCDP{}, "commercio/CloseCDP", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
