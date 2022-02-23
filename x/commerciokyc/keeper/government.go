package keeper

import (
	"fmt"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	eventAddTsp    = "add_tsp"
	eventRemoveTsp = "remove_tsp"
)

// AddTrustedServiceProvider allows to add the given signer as a trusted entity
// that can sign transactions setting an accrediter for a user.
func (k Keeper) AddTrustedServiceProvider(ctx sdk.Context, tsp sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	var trustedServiceProviders types.TrustedServiceProviders
	var signers ctypes.Strings
	signers = k.GetTrustedServiceProviders(ctx).Addresses
	if signersNew, inserted := signers.AppendIfMissing(tsp.String()); inserted {
		trustedServiceProviders.Addresses = signersNew
		newSignersBz, _ := k.cdc.MarshalBinaryBare(&trustedServiceProviders)
		store.Set([]byte(types.TrustedSignersStoreKey), newSignersBz)

	}
	// TODO Emits only if add
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventAddTsp,
		sdk.NewAttribute("tsp", tsp.String()),
	))
	logger := k.Logger(ctx)
	logger.Info("Tsp successfully added")

}

// RemoveTrustedServiceProvider allows to remove the given tsp from trusted entity
// list that can sign transactions setting an accrediter for a user.
func (k Keeper) RemoveTrustedServiceProvider(ctx sdk.Context, tsp sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	var trustedServiceProviders types.TrustedServiceProviders
	var signers ctypes.Strings
	signers = k.GetTrustedServiceProviders(ctx).Addresses
	if signersNew, find := signers.RemoveIfExisting(tsp.String()); find {
		trustedServiceProviders.Addresses = signersNew
		newSignersBz := k.cdc.MustMarshalBinaryBare(&trustedServiceProviders)
		store.Set([]byte(types.TrustedSignersStoreKey), newSignersBz)
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventRemoveTsp,
		sdk.NewAttribute("tsp", tsp.String()),
	))
	logger := k.Logger(ctx)
	logger.Info("Tsp successfully removed")

}

// DepositIntoPool allows the depositor to deposit the specified amount inside the rewards pool
func (k Keeper) DepositIntoPool(ctx sdk.Context, depositor sdk.AccAddress, amount sdk.Coins) error {
	// Send the coins from the user wallet to the pool
	var amountCoins sdk.Coins
	for _, coin := range amount {
		if coin.Denom != stakeDenom { // TODO change with constant
			return sdkErr.Wrap(sdkErr.ErrInsufficientFunds, fmt.Sprintf("deposit into membership pool can only be expressed in %s", stakeDenom))
		}
		amountCoins = append(amountCoins, coin)
	}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, amountCoins); err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventDeposit,
		sdk.NewAttribute("depositor", depositor.String()),
		sdk.NewAttribute("amount", amount.String()),
	))
	logger := k.Logger(ctx)
	logger.Info("Amount successfully deposited into pool")

	return nil
}

// GetTrustedServiceProviders returns the list of signers that are allowed to sign
// transactions setting a specific accrediter for a user.
// NOTE. Any user which is not present inside the returned list SHOULD NOT
// be allowed to send a transaction setting an accrediter for another user.
func (k Keeper) GetTrustedServiceProviders(ctx sdk.Context) (signers types.TrustedServiceProviders) {
	store := ctx.KVStore(k.storeKey)

	signersBz := store.Get([]byte(types.TrustedSignersStoreKey))
	k.cdc.UnmarshalBinaryBare(signersBz, &signers)

	//k.Cdc.MustUnmarshalBinaryBare(signersBz, &signers)
	// Cannot use add govAddress: trust service provider doesn't work proprerly
	//signers = append(signers, k.governmentKeeper.GetGovernmentAddress(ctx))
	return signers
}

// IsTrustedServiceProvider tells if the given signer is a trusted one or not
func (k Keeper) IsTrustedServiceProvider(ctx sdk.Context, signer sdk.Address) bool {
	var signers ctypes.Strings
	signers = k.GetTrustedServiceProviders(ctx).Addresses
	return signers.Contains(signer.String()) || signer.Equals(k.GovKeeper.GetGovernmentAddress(ctx))
}

// GetPoolFunds returns the funds currently present inside the rewards pool
func (k Keeper) GetPoolFunds(ctx sdk.Context) sdk.Coins {
	moduleAccount := k.GetMembershipModuleAccount(ctx)
	var coins sdk.Coins
	for _, coin := range k.bankKeeper.GetAllBalances(ctx, moduleAccount.GetAddress()) {
		coins = append(coins, coin)
	}
	return coins
}
