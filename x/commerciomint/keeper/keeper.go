package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
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

func (k Keeper) getCdpKey(address sdk.AccAddress) []byte {
	return []byte(types.CdpStorePrefix + address.String())
}

// AddCdp adds a Cdp to the user's Cdps list
func (k Keeper) AddCdp(ctx sdk.Context, cdp types.Cdp) {
	store := ctx.KVStore(k.storeKey)
	storeKey := k.getCdpKey(cdp.Owner)

	var cdps types.Cdps
	k.cdc.MustUnmarshalBinaryBare(store.Get(storeKey), &cdps)
	if cdps, edited := cdps.AppendIfMissing(cdp); edited {
		store.Set(storeKey, k.cdc.MustMarshalBinaryBare(cdps))
	}
}

// OpenCdp subtract the given token's amount from user's wallet and deposit it into the liquidity pool then,
// sending him the corresponding credits amount.
// If all these operations are done correctly, a Collateralized Debt Position is opened.
// Errors occurs if:
// 1) deposited tokens haven't been priced yet, or are negatives or invalid;
// 2) signer's funds are not enough
func (k Keeper) OpenCdp(ctx sdk.Context, depositor sdk.AccAddress, depositAmount sdk.Coins) error {

	if depositAmount.Empty() || !depositAmount.IsValid() {
		return sdkErr.Wrap(sdkErr.ErrInvalidCoins, fmt.Sprintf("Invalid deposit amount: %s", depositAmount))
	}

	// Check if all the tokens inside the deposit amount have a price and calculate the total fiat value of them
	fiatValue := sdk.NewInt(0)
	for _, token := range depositAmount {
		assetPrice, found := k.priceFeedKeeper.GetCurrentPrice(ctx, token.Denom)
		if !found {
			return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("No current price for given token: %s", token.Denom))
		}
		fiatValue = fiatValue.Add(token.Amount.Mul(assetPrice.Value.RoundInt()))
	}

	// Send the deposit from the user to the commerciomint account
	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, depositAmount)
	if err != nil {
		return err
	}

	// Get the credits amount
	// creditsAmount = (DepositAmount value / credits price) / 2
	// Our credit price is always 1 euro, so we simply divide the fiat value by 2

	// collateralRate is 2 here (cashcredit = fiat / collateralRate(government calls the shots here))
	creditsAmount := fiatValue.Quo(sdk.NewInt(2))

	// Mint the tokens and send them to the user
	credits := sdk.NewCoins(sdk.NewCoin(k.GetCreditsDenom(ctx), creditsAmount))
	err = k.supplyKeeper.MintCoins(ctx, types.ModuleName, credits)
	if err != nil {
		return err
	}

	err = k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, depositor, credits)
	if err != nil {
		return err
	}

	// Create the CDP and save it
	cdp := types.NewCdp(depositor, depositAmount, credits, ctx.BlockHeight())
	k.AddCdp(ctx, cdp)

	return nil
}

func (k Keeper) GetCdpsByOwner(ctx sdk.Context, owner sdk.AccAddress) (cdps types.Cdps) {
	store := ctx.KVStore(k.storeKey)
	k.cdc.MustUnmarshalBinaryBare(store.Get(k.getCdpKey(owner)), &cdps)
	return cdps
}

func (k Keeper) GetCdpByOwnerAndTimeStamp(ctx sdk.Context, owner sdk.AccAddress, timestamp int64) (cdp types.Cdp, found bool) {
	cdps := k.GetCdpsByOwner(ctx, owner)
	for _, ele := range cdps {
		if ele.Timestamp == timestamp {
			return ele, true
		}
	}
	return types.Cdp{}, false
}

func (k Keeper) GetCdps(ctx sdk.Context) types.Cdps {
	store := ctx.KVStore(k.storeKey)

	cdps := types.Cdps{}
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.CdpStorePrefix))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var cdp types.Cdps
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &cdp)
		cdps = append(cdps, cdp...)
	}

	return cdps
}

// CloseCdp subtract the Cdp's liquidity amount (commercio cash credits) from user's wallet, after that sends the
// deposited amount back to it. If these two operations ends without errors, the Cdp get closed.
// Errors occurs if:k.GetCdpsByOwner(ctx, testCdpOwner)
// - cdp doesnt exist
// - subtracting or adding fund to account don't end well
func (k Keeper) CloseCdp(ctx sdk.Context, user sdk.AccAddress, timestamp int64) error {
	cdp, found := k.GetCdpByOwnerAndTimeStamp(ctx, user, timestamp)
	if !found {
		msg := fmt.Sprintf("CDP for user with address %s and timestamp %d does not exist", user, timestamp)
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, msg)
	}

	// Send the coins from the user to the module and then burn them
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, cdp.Owner, types.ModuleName, cdp.CreditsAmount); err != nil {
		return err
	}
	if err := k.supplyKeeper.BurnCoins(ctx, types.ModuleName, cdp.CreditsAmount); err != nil {
		return err
	}

	// Get the user the deposited amount
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, cdp.Owner, cdp.DepositedAmount); err != nil {
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

func (k Keeper) deleteCdp(ctx sdk.Context, cdp types.Cdp) {
	store := ctx.KVStore(k.storeKey)

	var cdps types.Cdps
	k.cdc.MustUnmarshalBinaryBare(store.Get(k.getCdpKey(cdp.Owner)), &cdps)
	if cdps, found := cdps.RemoveWhenFound(cdp.Timestamp); found {
		if len(cdps) == 0 {
			store.Delete(k.getCdpKey(cdp.Owner))
		} else {
			store.Set(k.getCdpKey(cdp.Owner), k.cdc.MustMarshalBinaryBare(cdps))
		}
	}
}
