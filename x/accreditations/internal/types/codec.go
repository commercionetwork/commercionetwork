package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetAccrediter{}, "commercio/MsgSetAccrediter", nil)
	cdc.RegisterConcrete(MsgDistributeReward{}, "commercio/MsgDistributeReward", nil)
	cdc.RegisterConcrete(MsgDepositIntoLiquidityPool{}, "commercio/MsgDepositIntoLiquidityPool", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
