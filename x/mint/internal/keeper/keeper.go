package keeper

import (
	"fmt"
	"time"

	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
	"github.com/commercionetwork/commercionetwork/x/pricefeed"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

type Keeper struct {
	cdc             *codec.Codec
	storeKey        sdk.StoreKey
	priceFeedKeeper pricefeed.Keeper
	supplyKeeper    supply.Keeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, supplyKeeper supply.Keeper, pk pricefeed.Keeper) Keeper {

	// ensure mint module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		cdc:             cdc,
		storeKey:        key,
		priceFeedKeeper: pk,
		supplyKeeper:    supplyKeeper,
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
	return []byte(types.UserCdpsStorePrefix + address.String())
}

// AddCdp adds a Cdp to the user's Cdps list
func (k Keeper) AddCdp(ctx sdk.Context, cdp types.Cdp) {
	var cdps types.Cdps
	store := ctx.KVStore(k.storeKey)
	cdpsBz := store.Get(k.getCdpKey(cdp.Owner))
	k.cdc.MustUnmarshalBinaryBare(cdpsBz, &cdps)
	cdps, found := cdps.AppendIfMissing(cdp)
	if !found {
		store.Set(k.getCdpKey(cdp.Owner), k.cdc.MustMarshalBinaryBare(cdps))
	}
}

// OpenCdp subtract the given token's amount from user's wallet and deposit it into the liquidity pool then,
// sending him the corresponding credits amount.
// If all these operations are done correctly, a Collateralized Debt Position is opened.
// Errors occurs if:
// 1) deposited tokens haven't been priced yet, or are negatives or invalid;
// 2) signer's funds are not enough
func (k Keeper) OpenCdp(ctx sdk.Context, cdpRequest types.CdpRequest) sdk.Error {

	depositAmount := cdpRequest.DepositedAmount
	if !depositAmount.IsValid() || depositAmount.IsAnyNegative() || depositAmount.IsZero() {
		return sdk.ErrInvalidCoins(depositAmount.String())
	}

	// Check if all the tokens inside the deposit amount have a price and calculate the total fiat value of them
	fiatValue := sdk.NewInt(0)
	for _, token := range depositAmount {
		assetPrice, found := k.priceFeedKeeper.GetCurrentPrice(ctx, token.Denom)
		if !found {
			return sdk.ErrUnknownRequest(fmt.Sprintf("No current price for given token: %s", token.Denom))
		}
		fiatValue = fiatValue.Add(token.Amount.Mul(assetPrice.Price.RoundInt()))
	}

	// Send the deposit from the user to the mint account
	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, cdpRequest.Signer, types.ModuleName, depositAmount)
	if err != nil {
		return err
	}

	// Get the credits amount
	// creditsAmount = (DepositAmount value / credits price) / 2
	// Our credit price is always 1 euro, so we simply divide the fiat value by 2
	creditsAmount := fiatValue.Quo(sdk.NewInt(2))

	// Mint the tokens and send them to the user
	credits := sdk.NewCoins(sdk.NewCoin(k.GetCreditsDenom(ctx), creditsAmount))
	err = k.supplyKeeper.MintCoins(ctx, types.ModuleName, credits)
	if err != nil {
		return err
	}

	err = k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, cdpRequest.Signer, credits)
	if err != nil {
		return err
	}

	// Create the CDP and save it
	cdp := types.NewCdp(cdpRequest, credits)
	k.AddCdp(ctx, cdp)

	return nil
}

func (k Keeper) GetCdpsByOwner(ctx sdk.Context, owner sdk.AccAddress) (cdps types.Cdps) {
	store := ctx.KVStore(k.storeKey)
	k.cdc.MustUnmarshalBinaryBare(store.Get(k.getCdpKey(owner)), &cdps)
	return cdps
}

func (k Keeper) GetCdpByOwnerAndTimeStamp(ctx sdk.Context, owner sdk.AccAddress, timestamp time.Time) (cdp types.Cdp, found bool) {
	cdps := k.GetCdpsByOwner(ctx, owner)
	for _, ele := range cdps {
		if ele.Timestamp.Equal(timestamp) {
			return ele, true
		}
	}
	return types.Cdp{}, false
}

func (k Keeper) GetTotalCdps(ctx sdk.Context) types.Cdps {
	store := ctx.KVStore(k.storeKey)

	cdps := types.Cdps{}
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.UserCdpsStorePrefix))

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
// Errors occurs if:
// - cdp doesnt exist
// - subtracting or adding fund to account don't end well
func (k Keeper) CloseCdp(ctx sdk.Context, user sdk.AccAddress, timestamp time.Time) sdk.Error {
	cdp, found := k.GetCdpByOwnerAndTimeStamp(ctx, user, timestamp)
	if !found {
		msg := fmt.Sprintf("CDP for user with address %s and timestamp %s does not exist", user.String(), timestamp)
		return sdk.ErrUnknownRequest(msg)
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
