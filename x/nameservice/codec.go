package nameservice

/**
 * Any interface you create and any struct that implements an interface needs to be declared in the RegisterCodec function.
 * In this module the two Msg implementations (SetName and BuyName) need to be registered, but your Whois query return
 * type does not.
 */

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetName{}, "nameservice/SetName", nil)
	cdc.RegisterConcrete(MsgBuyName{}, "nameservice/BuyName", nil)
}