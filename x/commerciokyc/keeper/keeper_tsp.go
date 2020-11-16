package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
)

const (
	eventAddTsp    = "add_tsp"
	eventRemoveTsp = "remove_tsp"
)

// AddTrustedServiceProvider allows to add the given signer as a trusted entity
// that can sign transactions setting an accrediter for a user.
func (k Keeper) AddTrustedServiceProvider(ctx sdk.Context, tsp sdk.AccAddress) {
	store := ctx.KVStore(k.StoreKey)

	signers := k.GetTrustedServiceProviders(ctx)
	if signers, success := signers.AppendIfMissing(tsp); success {
		newSignersBz := k.Cdc.MustMarshalBinaryBare(&signers)
		store.Set([]byte(types.TrustedSignersStoreKey), newSignersBz)
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventAddTsp,
		sdk.NewAttribute("tsp", tsp.String()),
	))
}

// RemoveTrustedServiceProvider allows to remove the given tsp from trusted entity
// list that can sign transactions setting an accrediter for a user.
func (k Keeper) RemoveTrustedServiceProvider(ctx sdk.Context, tsp sdk.AccAddress) {
	store := ctx.KVStore(k.StoreKey)

	signers := k.GetTrustedServiceProviders(ctx)
	if signers, success := signers.RemoveIfExisting(tsp); success {
		newSignersBz := k.Cdc.MustMarshalBinaryBare(&signers)
		store.Set([]byte(types.TrustedSignersStoreKey), newSignersBz)
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventRemoveTsp,
		sdk.NewAttribute("tsp", tsp.String()),
	))
}

// GetTrustedServiceProviders returns the list of signers that are allowed to sign
// transactions setting a specific accrediter for a user.
// NOTE. Any user which is not present inside the returned list SHOULD NOT
// be allowed to send a transaction setting an accrediter for another user.
func (k Keeper) GetTrustedServiceProviders(ctx sdk.Context) (signers ctypes.Addresses) {
	store := ctx.KVStore(k.StoreKey)

	signersBz := store.Get([]byte(types.TrustedSignersStoreKey))
	k.Cdc.MustUnmarshalBinaryBare(signersBz, &signers)
	// Cannot use add govAddress: trust service provider doesn't work proprerly
	//signers = append(signers, k.governmentKeeper.GetGovernmentAddress(ctx))
	return
}

// IsTrustedServiceProvider tells if the given signer is a trusted one or not
func (k Keeper) IsTrustedServiceProvider(ctx sdk.Context, signer sdk.Address) bool {

	signers := k.GetTrustedServiceProviders(ctx)
	return signers.Contains(signer) || signer.Equals(k.governmentKeeper.GetGovernmentAddress(ctx))
}

// TspIterator returns an Iterator for all the tsps stored.
func (k Keeper) TspIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.StoreKey)
	return sdk.KVStorePrefixIterator(store, []byte(types.TrustedSignersStoreKey))
}
