package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

/**
 * Any interface you create and any struct that implements an interface needs to be
 * declared in the RegisterCodec function.
 * In this module the Msg implementations (SetIdentity) need to be registered,
 * but your Identity query return type does not.
 */

// RegisterCodec registers concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgStoreDocument{}, "commerciodocs/StoreDocument", nil)
	cdc.RegisterConcrete(MsgShareDocument{}, "commerciodocs/ShareDocument", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
