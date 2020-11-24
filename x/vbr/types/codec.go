package types

import "github.com/cosmos/cosmos-sdk/codec"

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgIncrementBlockRewardsPool{}, "commercio/MsgIncrementBlockRewardsPool", nil)
	cdc.RegisterConcrete(MsgSetRewardRate{}, "commercio/MsgSetRewardRate", nil)
	cdc.RegisterConcrete(MsgSetAutomaticWithdraw{}, "commercio/MsgSetAutomaticWithdraw", nil)

}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
