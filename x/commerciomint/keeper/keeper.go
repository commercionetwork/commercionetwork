package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/pkg/errors"

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
// --- Positions
// --------------

func (k Keeper) SetPosition(ctx sdk.Context, cdp types.Position) {
	store := ctx.KVStore(k.storeKey)
	key := makeCdpKey(cdp.Owner, cdp.CreatedAt)
	if bs := store.Get(key); bs != nil {
		panic(fmt.Errorf("cannot overwrite cdp at key %s", key))
	}
	store.Set(key, k.cdc.MustMarshalBinaryBare(cdp))
}

func (k Keeper) GetPosition(ctx sdk.Context, owner sdk.AccAddress, createdAt int64) (types.Position, bool) {
	cdp := types.Position{}
	key := makeCdpKey(owner, createdAt)
	store := ctx.KVStore(k.storeKey)
	bs := store.Get(key)
	if bs == nil {
		return cdp, false
	}
	k.cdc.MustUnmarshalBinaryBare(bs, &cdp)
	return cdp, true
}

func (k Keeper) GetAllPositionsOwnedBy(ctx sdk.Context, owner sdk.AccAddress) []types.Position {
	cdps := []types.Position{}
	i := k.newPositionsByOwnerIterator(ctx, owner)
	defer i.Close()
	for ; i.Valid(); i.Next() {
		var cdp types.Position
		k.cdc.MustUnmarshalBinaryBare(i.Value(), &cdp)
		cdps = append(cdps, cdp)
	}
	return cdps
}

// NewPosition subtract the given token's amount from user's wallet and deposit it into the liquidity pool then,
// sending him the corresponding credits amount.
// If all these operations are done correctly, a Collateralized Debt Position is opened.
// Errors occurs if:
// 1) deposited tokens haven't been priced yet, or are negatives or invalid;
// 2) signer's funds are not enough
func (k Keeper) NewPosition(ctx sdk.Context, depositor sdk.AccAddress, deposit sdk.Coins) error {
	fiatValue, err := k.calculateFiatValue(ctx, deposit)
	if err != nil {
		return err
	}

	// Get the credits amount
	// creditsAmount = (DepositAmount value / credits price) / collateral_rate
	// Our credit price is always 1, so we simply divide the fiat value by collateral_rate
	creditsAmount := fiatValue.Quo(k.GetCollateralRate(ctx)).TruncateInt()

	// Create credits token
	credits := sdk.NewCoin(k.GetCreditsDenom(ctx), creditsAmount)
	// Create the CDP and validate it
	position := types.NewPosition(depositor, deposit, credits, ctx.BlockHeight())
	if err := position.Validate(); err != nil {
		return errors.Wrap(err, "invalid position")
	}

	// Mint the tokens and send them to the user
	creditsCoins := sdk.NewCoins(credits)
	if err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, creditsCoins); err != nil {
		return errors.Wrap(err, "couldn't mint credits")
	}
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositor, creditsCoins); err != nil {
		return err
	}

	// Send the deposit from the user to the commerciomint account
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, deposit); err != nil {
		return err
	}

	// Create positon
	k.SetPosition(ctx, position)

	return nil
}

func (k Keeper) GetAllPositions(ctx sdk.Context) []types.Position {
	cdps := []types.Position{}
	iterator := k.newPositionsIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var cdp types.Position
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &cdp)
		cdps = append(cdps, cdp)
	}

	return cdps
}

// CloseCdp subtract the Position's liquidity amount (commercio cash credits) from user's wallet, after that sends the
// deposited amount back to it. If these two operations ends without errors, the Position get closed.
// Errors occurs if:k.GetCdpsByOwner(ctx, testCdpOwner)
// - cdp doesnt exist
// - subtracting or adding fund to account don't end well
func (k Keeper) CloseCdp(ctx sdk.Context, user sdk.AccAddress, timestamp int64) error {
	cdp, found := k.GetPosition(ctx, user, timestamp)
	if !found {
		msg := fmt.Sprintf("position for user with address %s and timestamp %d does not exist", user, timestamp)
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, msg)
	}

	// Send the coins from the user to the module and then burn them
	creditsCoins := sdk.NewCoins(cdp.Credits)
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, cdp.Owner, types.ModuleName, creditsCoins); err != nil {
		return err
	}
	if err := k.supplyKeeper.BurnCoins(ctx, types.ModuleName, creditsCoins); err != nil {
		return err
	}

	// Get the user the deposited amount
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, cdp.Owner, cdp.Deposit); err != nil {
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
func (k Keeper) ShouldLiquidateCdp(ctx sdk.Context, cdp types.Position) (bool, error) {
	fiatValue, err := k.calculateFiatValue(ctx, cdp.Deposit)
	if err != nil {
		return false, err
	}

	creditsAmount := cdp.Credits.Amount.ToDec()
	if creditsAmount.LTE(fiatValue) {
		return true, nil
	}
	return false, nil
}

func (k Keeper) AutoLiquidateCdps(ctx sdk.Context) {
	for _, cdp := range k.GetAllPositions(ctx) {
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

func (k Keeper) newPositionsByOwnerIterator(ctx sdk.Context, owner sdk.AccAddress) sdk.Iterator {
	prefix := []byte(fmt.Sprintf("%s%s:", types.CdpStorePrefix, owner.String()))
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), prefix)
}

func (k Keeper) newPositionsIterator(ctx sdk.Context) sdk.Iterator {
	return sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), []byte(types.CdpStorePrefix))
}

func (k Keeper) calculateFiatValue(ctx sdk.Context, deposits sdk.Coins) (sdk.Dec, error) {
	fiatValue := sdk.ZeroDec()
	for _, deposit := range deposits {
		assetPrice, found := k.priceFeedKeeper.GetCurrentPrice(ctx, deposit.Denom)
		if !found {
			return fiatValue, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("no current price for given denom: %s", deposit.Denom))
		}
		fiatValue = fiatValue.Add(deposit.Amount.ToDec().Mul(assetPrice.Value))
	}
	return fiatValue, nil
}

func (k Keeper) liquidate(ctx sdk.Context, cdp types.Position) error {
	// Send the coins from the user to the module and then burn them
	if err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, creditrisk.ModuleName, cdp.Deposit); err != nil {
		return err
	}
	// Delete the CDP
	k.deleteCdp(ctx, cdp)
	return nil
}

func makeCdpKey(address sdk.AccAddress, height int64) []byte {
	return []byte(fmt.Sprintf("%s%s:%d", types.CdpStorePrefix, address.String(), height))
}

func (k Keeper) deleteCdp(ctx sdk.Context, cdp types.Position) {
	store := ctx.KVStore(k.storeKey)
	key := makeCdpKey(cdp.Owner, cdp.CreatedAt)
	if bs := store.Get(key); bs == nil {
		panic(fmt.Sprintf("no cdp stored at key %s", key))
	}
	store.Delete(key)
}
