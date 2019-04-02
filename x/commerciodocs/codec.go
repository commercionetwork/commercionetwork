package commerciodocs

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgStoreDocument{}, "commerciodocs/StoreDocument", nil)
	cdc.RegisterConcrete(MsgShareDocument{}, "commerciodocs/ShareDocument", nil)
}

var msgCdc = codec.New()
