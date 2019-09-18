package keeper

import (
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
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

//SetPrice sets the current price for a given token after checking the validity of the signer
//If the signer hasn't the rights to set the price or if the price is zero or negative, then function returns error
func (keeper Keeper) SetPrice(ctx sdk.Context, signer sdk.AccAddress, price sdk.Dec, tokenName string) sdk.Error {
	return nil
}

//GetPrice retrieves the current price for the given token name
func (keeper Keeper) GetPrice(ctx sdk.Context, tokenName string) sdk.Int {
	return sdk.Int{}
}

//ValidateSigner makes sure the signer posting the price is an oracle
func (keeper Keeper) ValidateSigner(ctx sdk.Context, signer sdk.AccAddress) sdk.Error {
	return nil
}

// AddOracle adds an Oracle to the store
func (keeper Keeper) AddOracle(ctx sdk.Context, oracle sdk.AccAddress) {

}

func (keeper Keeper) GetOracles(ctx sdk.Context) ctypes.Addresses {}
