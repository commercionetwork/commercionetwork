package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	StoreKey sdk.StoreKey
	cdc      *codec.Codec
}

func NewKeeper(storekey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		StoreKey: storekey,
		cdc:      cdc,
	}
}

//SetTokenPrice sets the current price for a given token after checking the validity of the signer
//If the signer hasn't the rights to set the price or if the price is zero or negative, then function returns error
func (keeper Keeper) SetTokenPrice(ctx sdk.Context, signer sdk.AccAddress, price sdk.Dec, tokenName string) sdk.Error {
	return nil
}

//GetTokenPrice retrieves the current price for the given token name
func (keeper Keeper) GetTokenPrice(ctx sdk.Context, tokenName string) sdk.Int {
	return sdk.Int{}
}

//ValidateSigner checks if the signer who's trying to post a new price has the rights to do that
func (keeper Keeper) ValidateSigner(ctx sdk.Context, signer sdk.AccAddress) sdk.Error {
	return nil
}
