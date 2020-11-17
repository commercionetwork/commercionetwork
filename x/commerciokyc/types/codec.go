package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgInviteUser{}, "commercio/MsgInviteUser", nil)
	cdc.RegisterConcrete(MsgDepositIntoLiquidityPool{}, "commercio/MsgDepositIntoLiquidityPool", nil)
	cdc.RegisterConcrete(MsgAddTsp{}, "commercio/MsgAddTsp", nil)
	cdc.RegisterConcrete(MsgRemoveTsp{}, "commercio/MsgRemoveTsp", nil)
	cdc.RegisterConcrete(MsgBuyMembership{}, "commercio/MsgBuyMembership", nil)
	cdc.RegisterConcrete(MsgSetMembership{}, "commercio/MsgSetMembership", nil)
	cdc.RegisterConcrete(MsgRemoveMembership{}, "commercio/MsgRemoveMembership", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
