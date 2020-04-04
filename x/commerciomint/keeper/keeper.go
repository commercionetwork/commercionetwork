package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	creditrisk "github.com/commercionetwork/commercionetwork/x/creditrisk/types"
	"github.com/commercionetwork/commercionetwork/x/government"
	"github.com/commercionetwork/commercionetwork/x/pricefeed"
)

type Keeper struct {
	cdc             *codec.Codec
	storeKey        sdk.StoreKey
	priceFeedKeeper pricefeed.Keeper
	supplyKeeper    supply.Keeper
	govKeeper       government.Keeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, supplyKeeper supply.Keeper, pk pricefeed.Keeper, govKeeper government.Keeper) Keeper {
	// ensure commerciomint module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		cdc:             cdc,
		storeKey:        key,
		priceFeedKeeper: pk,
		supplyKeeper:    supplyKeeper,
		govKeeper:       govKeeper,
	}
}

// --------------
// --- Credits
// --------------

func (k Keeper) SetCreditsDenom(ctx sdk.Context, den string) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.CreditsDenomStoreKey), []byte(den))
}

func (k Keeper) GetCreditsDenom(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get([]byte(types.CreditsDenomStoreKey)))
}

// --------------
// --- CDPs
// --------------

func (k Keeper) SetCdp(ctx sdk.Context, cdp types.Cdp) {
	store := ctx.KVStore(k.storeKey)
	key := makeCdpKey(cdp.Owner, cdp.CreatedAt)
	if bs := store.Get(key); bs != nil {
		panic(fmt.Errorf("cannot overwrite cdp at key %s", key))
	}
	store.Set(key, k.cdc.MustMarshalBinaryBare(cdp))
}

func (k Keeper) GetCdp(ctx sdk.Context, owner sdk.AccAddress, createdAt int64) (types.Cdp, bool) {
	cdp := types.Cdp{}
	key := makeCdpKey(owner, createdAt)
	store := ctx.KVStore(k.storeKey)
	bs := store.Get(key)
	if bs == nil {
		return cdp, false
	}
	k.cdc.MustUnmarshalBinaryBare(bs, &cdp)
	return cdp, true
}

func (k Keeper) CdpsByOwnerIterator(ctx sdk.Context, owner sdk.AccAddress) sdk.Iterator {
	prefix := []byte(fmt.Sprintf("%s%s:", types.CdpStorePrefix, owner.String()))
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), prefix)
}

func (k Keeper) CdpsIterator(ctx sdk.Context) sdk.Iterator {
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), []byte(types.CdpStorePrefix))
}

func (k Keeper) CdpsByOwner(ctx sdk.Context, owner sdk.AccAddress) []types.Cdp {
	cdps := []types.Cdp{}
	i := k.CdpsByOwnerIterator(ctx, owner)
	defer i.Close()
	for ; i.Valid(); i.Next() {
		var cdp types.Cdp
		k.cdc.MustUnmarshalBinaryBare(i.Value(), &cdp)
		cdps = append(cdps, cdp)
	}
	return cdps
}

// OpenCdp subtract the given token's amount from user's wallet and deposit it into the liquidity pool then,
// sending him the corresponding credits amount.
// If all these operations are done correctly, a Collateralized Debt Position is opened.
// Errors occurs if:
// 1) deposited tokens haven't been priced yet, or are negatives or invalid;
// 2) signer's funds are not enough
func (k Keeper) OpenCdp(ctx sdk.Context, depositor sdk.AccAddress, deposit sdk.Coin) error {
	if !deposit.IsValid() || !deposit.IsPositive() {
		return sdkErr.Wrap(sdkErr.ErrInvalidCoins, fmt.Sprintf("Invalid deposit amount: %s", deposit))
	}

	// Check if all the tokens inside the deposit amount have a price and calculate the total fiat value of them
	fiatValue, err := k.calculateFiatValue(ctx, deposit)
	if err != nil {
		return err
	}

	// Send the deposit from the user to the commerciomint account
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, sdk.NewCoins(deposit)); err != nil {
		return err
	}

	// Get the credits amount
	// creditsAmount = (DepositAmount value / credits price) / collateral_rate
	// Our credit price is always 1, so we simply divide the fiat value by collateral_rate
	creditsAmount := fiatValue.Quo(k.GetCollateralRate(ctx)).TruncateInt()

	// Mint the tokens and send them to the user
	credits := sdk.NewCoins(sdk.NewCoin(k.GetCreditsDenom(ctx), creditsAmount))
	if err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, credits); err != nil {
		return err
	}

	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositor, credits); err != nil {
		return err
	}

	// Create the CDP and save it
	cdp := types.NewCdp(depositor, deposit, credits, ctx.BlockHeight())
	k.SetCdp(ctx, cdp)

	return nil
}

func (k Keeper) GetCdps(ctx sdk.Context) types.Cdps {
	cdps := []types.Cdp{}
	iterator := k.CdpsIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var cdp types.Cdp
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &cdp)
		cdps = append(cdps, cdp)
	}

	return cdps
}

// CloseCdp subtract the Cdp's liquidity amount (commercio cash credits) from user's wallet, after that sends the
// deposited amount back to it. If these two operations ends without errors, the Cdp get closed.
// Errors occurs if:k.GetCdpsByOwner(ctx, testCdpOwner)
// - cdp doesnt exist
// - subtracting or adding fund to account don't end well
func (k Keeper) CloseCdp(ctx sdk.Context, user sdk.AccAddress, timestamp int64) error {
	cdp, found := k.GetCdp(ctx, user, timestamp)
	if !found {
		msg := fmt.Sprintf("CDP for user with address %s and timestamp %d does not exist", user, timestamp)
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, msg)
	}

	// Send the coins from the user to the module and then burn them
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, cdp.Owner, types.ModuleName, cdp.Credits); err != nil {
		return err
	}
	if err := k.supplyKeeper.BurnCoins(ctx, types.ModuleName, cdp.Credits); err != nil {
		return err
	}

	// Get the user the deposited amount
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, cdp.Owner, sdk.NewCoins(cdp.Deposit)); err != nil {
		return err
	}

	// Delete the CDP
	k.deleteCdp(ctx, cdp)

	return nil
}

// GetCollateralRate retrieve the cdp collateral rate.
func (k Keeper) GetCollateralRate(ctx sdk.Context) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	var rate sdk.Dec
	k.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.CollateralRateKey)), &rate)
	return rate
}

// SetCollateralRate store the cdp collateral rate.
func (k Keeper) SetCollateralRate(ctx sdk.Context, rate sdk.Dec) error {
	if err := types.ValidateCollateralRate(rate); err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.CollateralRateKey), k.cdc.MustMarshalBinaryBare(rate))
	return nil
}

// ShouldLiquidateCdp returns true if the CDP should be liquidated.
func (k Keeper) ShouldLiquidateCdp(ctx sdk.Context, cdp types.Cdp) (bool, error) {
	fiatValue, err := k.calculateFiatValue(ctx, cdp.Deposit)
	if err != nil {
		return false, err
	}

	creditsAmount := cdp.Credits.AmountOf(k.GetCreditsDenom(ctx)).ToDec()
	if creditsAmount.LTE(fiatValue) {
		return true, nil
	}
	return false, nil
}

func (k Keeper) AutoLiquidateCdps(ctx sdk.Context) {
	for _, cdp := range k.GetCdps(ctx) {
		if yes, err := k.ShouldLiquidateCdp(ctx, cdp); err != nil {
			panic(err)
		} else if !yes {
			continue
		}
		if err := k.liquidate(ctx, cdp); err != nil {
			panic(err)
		}
	}
}

func (k Keeper) calculateFiatValue(ctx sdk.Context, deposit sdk.Coin) (sdk.Dec, error) {
	// Check if all the tokens inside the deposit amount have a price and calculate the total fiat value of them
	fiatValue := sdk.ZeroDec()
	assetPrice, found := k.priceFeedKeeper.GetCurrentPrice(ctx, deposit.Denom)
	if !found {
		return fiatValue, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("no current price for given denom: %s", deposit.Denom))
	}
	fiatValue = fiatValue.Add(deposit.Amount.ToDec().Mul(assetPrice.Value))
	return fiatValue, nil
}

func (k Keeper) liquidate(ctx sdk.Context, cdp types.Cdp) error {
	// Send the coins from the user to the module and then burn them
	if err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, creditrisk.ModuleName, sdk.NewCoins(cdp.Deposit)); err != nil {
		return err
	}
	// Delete the CDP
	k.deleteCdp(ctx, cdp)
	return nil
}

func makeCdpKey(address sdk.AccAddress, height int64) []byte {
	return []byte(fmt.Sprintf("%s%s:%d", types.CdpStorePrefix, address.String(), height))
}

func (k Keeper) deleteCdp(ctx sdk.Context, cdp types.Cdp) {
	store := ctx.KVStore(k.storeKey)
	key := makeCdpKey(cdp.Owner, cdp.CreatedAt)
	if bs := store.Get(key); bs == nil {
		panic(fmt.Sprintf("no cdp stored at key %s", key))
	}
	store.Delete(key)
}
