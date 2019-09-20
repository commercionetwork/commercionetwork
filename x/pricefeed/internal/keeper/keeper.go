package keeper

import (
	"fmt"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/government"
	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	StoreKey         sdk.StoreKey
	GovernmentKeeper government.Keeper
	cdc              *codec.Codec
}

func NewKeeper(storekey sdk.StoreKey, govK government.Keeper, cdc *codec.Codec) Keeper {
	return Keeper{
		StoreKey:         storekey,
		GovernmentKeeper: govK,
		cdc:              cdc,
	}
}

//SetRawPrice sets the raw price for a given token after checking the validity of the signer
//If the signer hasn't the rights to set the price, then function returns error
func (keeper Keeper) SetRawPrice(ctx sdk.Context, price types.RawPrice) sdk.Error {
	err := keeper.ValidateSigner(ctx, price.Oracle)
	if err != nil {
		return err
	}
	rawPrices := keeper.GetRawPrices(ctx)
	rawPrices.UpdatePriceOrAppendIfMissing(price)
	return nil
}

//GetRawPrices retrieves all the current prices
func (keeper Keeper) GetRawPrices(ctx sdk.Context) types.RawPrices {
	store := ctx.KVStore(keeper.StoreKey)
	pricesBz := store.Get([]byte(types.RawPricesPrefix))
	var rawPrices types.RawPrices
	keeper.cdc.MustUnmarshalBinaryBare(pricesBz, &rawPrices)
	return rawPrices
}

func (keeper Keeper) SetCurrentPrice(ctx sdk.Context) sdk.Error {
	/*
		GetRawPrices
		filter not expire ones
		filter for tokenName+tokeCode
		curPrice = avg(rawPrices)
	*/
}

//GetCurrentPrices retrieves all the current prices
func (keeper Keeper) GetCurrentPrices(ctx sdk.Context) types.CurrentPrices {
	store := ctx.KVStore(keeper.StoreKey)
	pricesBz := store.Get([]byte(types.CurrentPricesPrefix))
	var curPrices types.CurrentPrices
	keeper.cdc.MustUnmarshalBinaryBare(pricesBz, &curPrices)
	return curPrices
}

//GetCurrentPrice retrieves the current price for the given token name and code
func (keeper Keeper) GetCurrentPrice(ctx sdk.Context, tokenName string, tokenCode string) (types.CurrentPrice, sdk.Error) {
	currentPrices := keeper.GetCurrentPrices(ctx)
	price, err := currentPrices.FindPrice(tokenName, tokenCode)
	if err != nil {
		return types.CurrentPrice{}, err
	}
	return price, nil
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
