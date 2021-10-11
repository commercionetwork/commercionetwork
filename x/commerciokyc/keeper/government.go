package keeper

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

// AddTrustedServiceProvider allows to add the given signer as a trusted entity
// that can sign transactions setting an accrediter for a user.
func (k Keeper) AddTrustedServiceProvider(ctx sdk.Context, tsp sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	signers := k.GetTrustedServiceProviders(ctx)
	if signers, success := signers.AppendIfMissing(tsp); success {
		newSignersBz, _ := json.Marshal(signers) // TODO control this conversion
		//newSignersBz := k.cdc.MustMarshalBinaryBare(&signers)
		store.Set([]byte(types.TrustedSignersStoreKey), newSignersBz)
	}

	// TODO emits events
	//ctx.EventManager().EmitEvent(sdk.NewEvent(
	//	eventAddTsp,
	//	sdk.NewAttribute("tsp", tsp.String()),
	//))

}

// RemoveTrustedServiceProvider allows to remove the given tsp from trusted entity
// list that can sign transactions setting an accrediter for a user.
func (k Keeper) RemoveTrustedServiceProvider(ctx sdk.Context, tsp sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	signers := k.GetTrustedServiceProviders(ctx)
	if signers, success := signers.RemoveIfExisting(tsp); success {
		newSignersBz, _ := json.Marshal(signers) // TODO control this conversion
		//newSignersBz := k.cdc.MustMarshalBinaryBare(&signers)
		store.Set([]byte(types.TrustedSignersStoreKey), newSignersBz)
	}

	/*
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			eventRemoveTsp,
			sdk.NewAttribute("tsp", tsp.String()),
		))
	*/
}

// DepositIntoPool allows the depositor to deposit the specified amount inside the rewards pool
func (k Keeper) DepositIntoPool(ctx sdk.Context, depositor sdk.AccAddress, amount []*sdk.Coin) error {
	// Send the coins from the user wallet to the pool
	bondDenom := "ucommercio"
	var amountCoins sdk.Coins
	for _, coin := range amount {
		if coin.Denom != bondDenom { // TODO change with constant
			return sdkErr.Wrap(sdkErr.ErrInsufficientFunds, fmt.Sprintf("deposit into membership pool can only be expressed in %s", bondDenom))
		}
		amountCoins = append(amountCoins, *coin)
	}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, amountCoins); err != nil {
		return err
	}
	/*
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			eventDeposit,
			sdk.NewAttribute("depositor", depositor.String()),
			sdk.NewAttribute("amount", amount.String()),
		))
	*/

	return nil
}
