package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/accreditation/internal/types"
	utypes "github.com/commercionetwork/commercionetwork/x/utils/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type Keeper struct {
	StoreKey sdk.StoreKey

	// Bank keeper to send tokens
	BankKeeper bank.Keeper

	// Pointer to the codec that is used by Amino to encode and decode binary structs.
	Cdc *codec.Codec
}

// NewKeeper creates new instances of the accreditation module Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, bankKeeper bank.Keeper) Keeper {
	return Keeper{
		StoreKey:   storeKey,
		BankKeeper: bankKeeper,
		Cdc:        cdc,
	}
}

// SetAccrediter allows to set a given user as being accreditated by the given accrediter.
func (keeper Keeper) SetAccrediter(ctx sdk.Context, accrediter sdk.AccAddress, user sdk.AccAddress) {
	store := ctx.KVStore(keeper.StoreKey)

	// Save the accrediter
	store.Set(user, accrediter)
}

// GetAccrediter allows to get the accrediter a given user
func (keeper Keeper) GetAccrediter(ctx sdk.Context, user sdk.AccAddress) (accrediter sdk.AccAddress, found bool) {
	store := ctx.KVStore(keeper.StoreKey)

	found = store.Has(user)
	if found {
		accrediter = store.Get(user)
	}

	return
}

// AddTrustworthySigner allows to add the given signer as a trustworthy entity
// that can sign transactions setting an accrediter for a user.
func (keeper Keeper) AddTrustworthySigner(ctx sdk.Context, signer sdk.AccAddress) {
	store := ctx.KVStore(keeper.StoreKey)

	signers := utypes.Addresses(keeper.GetTrustworthySigners(ctx))
	signers = signers.AppendIfMissing(signer)

	newSignersBz := keeper.Cdc.MustMarshalBinaryBare(&signers)
	store.Set([]byte(types.TrustworthySignersKey), newSignersBz)
}

// IsTrustworthySigner tells if the given signer is a trustworthy one or not
func (keeper Keeper) IsTrustworthySigner(ctx sdk.Context, signer sdk.AccAddress) bool {
	signers := keeper.GetTrustworthySigners(ctx)
	for _, s := range signers {
		if s.Equals(signer) {
			return true
		}
	}
	return false
}

// -----------------------
// --- Genesis utils
// -----------------------

// GetAccreditations returns all the accreditations that have been
func (keeper Keeper) GetAccreditations(ctx sdk.Context) (accreditations []types.Accreditation) {
	store := ctx.KVStore(keeper.StoreKey)
	iterator := store.Iterator(nil, nil)

	for ; iterator.Valid() && string(iterator.Key()) != types.TrustworthySignersKey; iterator.Next() {
		accreditation := types.Accreditation{
			Accrediter: iterator.Value(),
			User:       iterator.Value(),
		}
		accreditations = append(accreditations, accreditation)
	}

	return
}

// GetTrustworthySigners returns the list of signers that are allowed to sign
// transactions setting a specific accrediter for a user.
// NOTE. Any user which is not present inside the returned list SHOULD NOT
// be allowed to send a transaction setting an accrediter for another user.
func (keeper Keeper) GetTrustworthySigners(ctx sdk.Context) (signers []sdk.AccAddress) {
	store := ctx.KVStore(keeper.StoreKey)

	signersBz := store.Get([]byte(types.TrustworthySignersKey))
	keeper.Cdc.MustUnmarshalBinaryBare(signersBz, &signers)

	return
}
