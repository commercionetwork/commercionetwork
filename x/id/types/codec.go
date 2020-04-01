package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetIdentity{}, "commercio/MsgSetIdentity", nil)
	cdc.RegisterConcrete(MsgRequestDidPowerUp{}, "commercio/MsgRequestDidPowerUp", nil)
	cdc.RegisterConcrete(MsgChangePowerUpStatus{}, "commercio/MsgChangePowerUpStatus", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
