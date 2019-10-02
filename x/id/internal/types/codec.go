package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetIdentity{}, "commercio/MsgSetIdentity", nil)
	cdc.RegisterConcrete(MsgRequestDidDeposit{}, "commercio/MsgRequestDidDeposit", nil)
	cdc.RegisterConcrete(MsgChangeDidDepositRequestStatus{}, "commercio/MsgChangeDidDepositRequestStatus", nil)
	cdc.RegisterConcrete(MsgRequestDidPowerUp{}, "commercio/MsgRequestDidPowerUp", nil)
	cdc.RegisterConcrete(MsgChangeDidPowerUpRequestStatus{}, "commercio/MsgChangeDidPowerUpRequestStatus", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
