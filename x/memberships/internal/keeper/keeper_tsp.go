package keeper

import (
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AddTrustedServiceProvider allows to add the given signer as a trusted entity
// that can sign transactions setting an accrediter for a user.
func (k Keeper) AddTrustedServiceProvider(ctx sdk.Context, tsp sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	signers := k.GetTrustedServiceProviders(ctx)
	if signers, success := signers.AppendIfMissing(tsp); success {
		newSignersBz := k.cdc.MustMarshalBinaryBare(&signers)
		store.Set([]byte(types.TrustedSignersStoreKey), newSignersBz)
	}
}

// GetTrustedServiceProviders returns the list of signers that are allowed to sign
// transactions setting a specific accrediter for a user.
// NOTE. Any user which is not present inside the returned list SHOULD NOT
// be allowed to send a transaction setting an accrediter for another user.
func (k Keeper) GetTrustedServiceProviders(ctx sdk.Context) (signers ctypes.Addresses) {
	store := ctx.KVStore(k.storeKey)

	signersBz := store.Get([]byte(types.TrustedSignersStoreKey))
	k.cdc.MustUnmarshalBinaryBare(signersBz, &signers)

	return
}

// IsTrustedServiceProvider tells if the given signer is a trusted one or not
func (k Keeper) IsTrustedServiceProvider(ctx sdk.Context, signer sdk.Address) bool {
	signers := k.GetTrustedServiceProviders(ctx)
	return signers.Contains(signer)
}
