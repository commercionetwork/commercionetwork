package keeper

import (
	"fmt"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/types"
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
func (keeper Keeper) SetPrice(ctx sdk.Context, signer sdk.AccAddress, price sdk.Int, tokenName string, expiry sdk.Int) sdk.Error {
	return nil
}

//GetPrice retrieves the current price for the given token name
func (keeper Keeper) GetPrice(ctx sdk.Context, tokenName string) sdk.Int {
	return sdk.Int{}
}

//ValidateSigner makes sure the signer posting the price is an oracle
func (keeper Keeper) ValidateSigner(ctx sdk.Context, signer sdk.AccAddress) sdk.Error {
	oracles := keeper.GetOracles(ctx)
	isOracle := oracles.Contains(signer)
	if !isOracle {
		return sdk.ErrInvalidAddress(fmt.Sprintf("%s isn't an Oracle", signer))
	}
	return nil
}

// AddOracle adds an Oracle to the store
func (keeper Keeper) AddOracle(ctx sdk.Context, oracle sdk.AccAddress) {
	oracles := keeper.GetOracles(ctx)
	oracles = oracles.AppendIfMissing(oracle)
	store := ctx.KVStore(keeper.StoreKey)
	store.Set([]byte(types.OraclePrefix), keeper.cdc.MustMarshalBinaryBare(&oracles))
}

func (keeper Keeper) GetOracles(ctx sdk.Context) ctypes.Addresses {
	store := ctx.KVStore(keeper.StoreKey)
	oraclesBz := store.Get([]byte(types.OraclePrefix))
	var oracles ctypes.Addresses
	keeper.cdc.MustUnmarshalBinaryBare(oraclesBz, &oracles)
	return oracles
}

func (keeper Keeper) GetOracle(ctx sdk.Context, oracle sdk.AccAddress) (sdk.AccAddress, error) {
	oracles := keeper.GetOracles(ctx)
	found := oracles.GetAddress(oracle)
	if found == nil {
		return nil, sdk.ErrUnknownAddress("Oracle address not found")
	}
	return found, nil
}
