package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(upgrade.Plan{}, "cosmos-sdk/Plan", nil)
	cdc.RegisterConcrete(MsgScheduleUpgrade{}, fmt.Sprintf("upgrade/%s", ScheduleUpgradeConst), nil)
	cdc.RegisterConcrete(MsgDeleteUpgrade{}, fmt.Sprintf("upgrade/%s", DeleteUpgradeConst), nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
