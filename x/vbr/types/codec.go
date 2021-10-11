package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 2
	cdc.RegisterConcrete(&MsgSetAutomaticWithdraw{}, "vbr/SetAutomaticWithdraw", nil)

	cdc.RegisterConcrete(&MsgSetRewardRate{}, "vbr/SetRewardRate", nil)

	cdc.RegisterConcrete(&MsgIncrementBlockRewardsPool{}, "vbr/MsgIncrementBlockRewardsPool", nil)

}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetAutomaticWithdraw{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetRewardRate{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgIncrementBlockRewardsPool{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
