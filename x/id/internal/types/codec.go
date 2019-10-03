package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetIdentity{}, "commercio/MsgSetIdentity", nil)
	cdc.RegisterConcrete(MsgRequestDidDeposit{}, "commercio/MsgRequestDidDeposit", nil)
	cdc.RegisterConcrete(MsgInvalidateDidDepositRequest{}, "commercio/MsgInvalidateDidDepositRequest", nil)
	cdc.RegisterConcrete(MsgRequestDidPowerUp{}, "commercio/MsgRequestDidPowerUp", nil)
	cdc.RegisterConcrete(MsgInvalidateDidPowerUpRequest{}, "commercio/MsgInvalidateDidPowerUpRequest", nil)
	cdc.RegisterConcrete(MsgWithdrawDeposit{}, "commercio/MsgWithdrawDeposit", nil)
	cdc.RegisterConcrete(MsgPowerUpDid{}, "commercio/MsgPowerUpDid", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
