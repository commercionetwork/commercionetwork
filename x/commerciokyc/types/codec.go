package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgInviteUser{}, "commercio/MsgInviteUser", nil)
	cdc.RegisterConcrete(&MsgDepositIntoLiquidityPool{}, "commercio/MsgDepositIntoLiquidityPool", nil)
	cdc.RegisterConcrete(&MsgBuyMembership{}, "commercio/MsgBuyMembership", nil)
	cdc.RegisterConcrete(&MsgAddTsp{}, "commercio/MsgAddTsp", nil)
	cdc.RegisterConcrete(&MsgRemoveTsp{}, "commercio/MsgRemoveTsp", nil)
	cdc.RegisterConcrete(&MsgRemoveMembership{}, "commercio/MsgRemoveMembership", nil)
	cdc.RegisterConcrete(&MsgSetMembership{}, "commercio/MsgSetMembership", nil)

}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgInviteUser{},
		&MsgDepositIntoLiquidityPool{},
		&MsgBuyMembership{},
		&MsgAddTsp{},
		&MsgRemoveTsp{},
		&MsgRemoveMembership{},
		&MsgSetMembership{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
