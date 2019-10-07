package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
	"github.com/commercionetwork/commercionetwork/x/pricefeed"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type Keeper struct {
	StoreKey        sdk.StoreKey
	bankKeeper      bank.Keeper
	priceFeedKeeper pricefeed.Keeper
	cdc             *codec.Codec
}

func NewKeeper(sk sdk.StoreKey, bk bank.Keeper, pk pricefeed.Keeper, cdc *codec.Codec) Keeper {
	return Keeper{
		StoreKey:        sk,
		bankKeeper:      bk,
		priceFeedKeeper: pk,
		cdc:             cdc,
	}
}

// ---------------------
// --- Liquidity pool
// ---------------------

// SetLiquidityPool allows to set the given pool amount as the current liquidity pool amount
func (keeper Keeper) SetLiquidityPool(ctx sdk.Context, poolAmount sdk.Coins) {
	store := ctx.KVStore(keeper.StoreKey)
	storeKey := []byte(types.LiquidityPoolStorePrefix)

	if poolAmount == nil || poolAmount.Empty() {
		store.Delete(storeKey)
	} else {
		store.Set(storeKey, keeper.cdc.MustMarshalBinaryBare(&poolAmount))
	}
}

// GetLiquidityPool returns the amount of the liquidity pool for the given context
func (keeper Keeper) GetLiquidityPool(ctx sdk.Context) sdk.Coins {
	var lPool sdk.Coins
	store := ctx.KVStore(keeper.StoreKey)
	lpBz := store.Get([]byte(types.LiquidityPoolStorePrefix))
	keeper.cdc.MustUnmarshalBinaryBare(lpBz, &lPool)
	return lPool
}

// --------------
// --- Credits
// --------------

func (keeper Keeper) SetCreditsDenom(ctx sdk.Context, den string) {
	store := ctx.KVStore(keeper.StoreKey)
	store.Set([]byte(types.CreditsDenomStoreKey), []byte(den))
}

func (keeper Keeper) GetCreditsDenom(ctx sdk.Context) string {
	store := ctx.KVStore(keeper.StoreKey)
	return string(store.Get([]byte(types.CreditsDenomStoreKey)))
}

// --------------
// --- CDPs
// --------------

func (keeper Keeper) getCdpKey(address sdk.AccAddress) []byte {
	return []byte(types.UserCdpsStorePrefix + address.String())
}

// AddCdp adds a Cdp to the user's Cdps list
func (keeper Keeper) AddCdp(ctx sdk.Context, cdp types.Cdp) {
	var cdps types.Cdps
	store := ctx.KVStore(keeper.StoreKey)
	cdpsBz := store.Get(keeper.getCdpKey(cdp.Owner))
	keeper.cdc.MustUnmarshalBinaryBare(cdpsBz, &cdps)
	cdps, found := cdps.AppendIfMissing(cdp)
	if !found {
		store.Set(keeper.getCdpKey(cdp.Owner), keeper.cdc.MustMarshalBinaryBare(cdps))
	}
}

// OpenCdp subtract the given token's amount from user's wallet and deposit it into the liquidity pool then,
// sending him the corresponding credits amount.
// If all these operations are done correctly, a Collateralized Debt Position is opened.
// Errors occurs if:
// 1) deposited tokens haven't been priced yet, or are negatives or invalid;
// 2) signer's funds are not enough
func (keeper Keeper) OpenCdp(ctx sdk.Context, cdpRequest types.CdpRequest) sdk.Error {

	if !cdpRequest.DepositedAmount.IsValid() || cdpRequest.DepositedAmount.IsAnyNegative() ||
		cdpRequest.DepositedAmount.IsZero() {
		return sdk.ErrInvalidCoins(cdpRequest.DepositedAmount.String())
	}

	fiatValue := sdk.NewInt(0)

	//Check if all tokens in deposit amount have a price and calculate the total FIAT value of them
	for _, token := range cdpRequest.DepositedAmount {
		assetPrice, found := keeper.priceFeedKeeper.GetCurrentPrice(ctx, token.Denom)
		if found == false {
			return sdk.ErrInvalidCoins(fmt.Sprintf("no current price for given token: %s", token.Denom))
		}
		fiatValue = fiatValue.Add(token.Amount.Mul(assetPrice.Price.RoundInt()))
	}

	//Subtract the given deposit amount from user's wallet
	_, err := keeper.bankKeeper.SubtractCoins(ctx, cdpRequest.Signer, cdpRequest.DepositedAmount)
	if err != nil {
		return err
	}

	liquidityPool := keeper.GetLiquidityPool(ctx)
	//depositing the amount to the liquidity pool
	liquidityPool = liquidityPool.Add(cdpRequest.DepositedAmount)
	keeper.SetLiquidityPool(ctx, liquidityPool)

	//get credits' amount = DepositAmount value / credits price (always 1 euro) / 2 is the power of collateral which is 2:1 (comm -> ccc)
	creditsAmount := fiatValue.Quo(sdk.NewInt(2))

	//add credits to users wallet
	credits := sdk.NewCoins(sdk.NewCoin(keeper.GetCreditsDenom(ctx), creditsAmount))
	credits, _ = keeper.bankKeeper.AddCoins(ctx, cdpRequest.Signer, credits)

	//Creating Cdp and adding to the store
	cdp := types.NewCdp(cdpRequest, credits)
	keeper.AddCdp(ctx, cdp)

	return nil
}

func (keeper Keeper) GetCdpsByOwner(ctx sdk.Context, owner sdk.AccAddress) (cdps types.Cdps) {
	store := ctx.KVStore(keeper.StoreKey)
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(keeper.getCdpKey(owner)), &cdps)
	return cdps
}

func (keeper Keeper) GetCdpByOwnerAndTimeStamp(ctx sdk.Context, owner sdk.AccAddress, timestamp string) (cdp types.Cdp, found bool) {
	cdps := keeper.GetCdpsByOwner(ctx, owner)
	for _, ele := range cdps {
		if ele.Timestamp == timestamp {
			return ele, true
		}
	}
	return types.Cdp{}, false
}

func (keeper Keeper) GetTotalCdps(ctx sdk.Context) types.Cdps {
	store := ctx.KVStore(keeper.StoreKey)

	cdps := types.Cdps{}
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.UserCdpsStorePrefix))
	for ; iterator.Valid(); iterator.Next() {
		var cdp types.Cdp
		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &cdp)
		cdps = append(cdps, cdp)
	}

	return cdps
}

// CloseCdp subtract the Cdp's liquidity amount (commercio cash credits) from user's wallet, after that sends the
// deposited amount back to it. If these two operations ends without errors, the Cdp get closed.
// Errors occurs if:
// - cdp doesnt exist
// - subtracting or adding fund to account don't end well
func (keeper Keeper) CloseCdp(ctx sdk.Context, user sdk.AccAddress, timestamp string) sdk.Error {
	cdp, found := keeper.GetCdpByOwnerAndTimeStamp(ctx, user, timestamp)
	if !found {
		msg := fmt.Sprintf("CDP for user with address %s and timestamp %s does not exist", user.String(), timestamp)
		return sdk.ErrUnknownRequest(msg)
	}

	// Subtract the liquidity amount from the user's wallet
	_, err := keeper.bankKeeper.SubtractCoins(ctx, user, cdp.CreditsAmount)
	if err != nil {
		return err
	}

	// Withdraw the previously deposited amount from the liquidity pool
	liquidityPool := keeper.GetLiquidityPool(ctx)
	liquidityPool = liquidityPool.Sub(cdp.DepositedAmount)
	keeper.SetLiquidityPool(ctx, liquidityPool)

	// Add the liquidity amount to the user's wallet
	_, err = keeper.bankKeeper.AddCoins(ctx, user, cdp.DepositedAmount)
	if err != nil {
		return err
	}

	keeper.deleteCdp(ctx, cdp)

	return nil
}

func (keeper Keeper) deleteCdp(ctx sdk.Context, cdp types.Cdp) {
	store := ctx.KVStore(keeper.StoreKey)

	var cdps types.Cdps
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(keeper.getCdpKey(cdp.Owner)), &cdps)
	cdps, found := cdps.RemoveWhenFound(cdp.Timestamp)

	if found {
		if len(cdps) == 0 {
			store.Delete(keeper.getCdpKey(cdp.Owner))
		} else {
			store.Set(keeper.getCdpKey(cdp.Owner), keeper.cdc.MustMarshalBinaryBare(cdps))
		}
	}
}
