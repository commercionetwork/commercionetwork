package commercioauth

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateAccount{}, "commercioauth/CreateAccount", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
