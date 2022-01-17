package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 2
	cdc.RegisterConcrete(&MsgInviteUser{}, "commercio/MsgInviteUser", nil)
	cdc.RegisterConcrete(&MsgDepositIntoLiquidityPool{}, "commercio/MsgDepositIntoLiquidityPool", nil)
	cdc.RegisterConcrete(&MsgBuyMembership{}, "commercio/MsgBuyMembership", nil)
	cdc.RegisterConcrete(&MsgAddTsp{}, "commercio/MsgAddTsp", nil)
	cdc.RegisterConcrete(&MsgRemoveTsp{}, "commercio/MsgRemoveTsp", nil)
	cdc.RegisterConcrete(&MsgRemoveMembership{}, "commercio/MsgRemoveMembership", nil)
	cdc.RegisterConcrete(&MsgSetMembership{}, "commercio/MsgSetMembership", nil)
	cdc.RegisterConcrete(&MsgSetParams{}, "commercio/SetParams", nil)

}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgInviteUser{},
		&MsgDepositIntoLiquidityPool{},
		&MsgBuyMembership{},
		&MsgAddTsp{},
		&MsgRemoveTsp{},
		&MsgRemoveMembership{},
		&MsgSetMembership{},
		&MsgSetParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
