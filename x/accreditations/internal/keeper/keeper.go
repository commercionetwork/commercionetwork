package keeper

import (
	"errors"

	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/types"
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
	cdc *codec.Codec
}

// NewKeeper creates new instances of the accreditation module Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, bankKeeper bank.Keeper) Keeper {
	return Keeper{
		StoreKey:   storeKey,
		BankKeeper: bankKeeper,
		cdc:        cdc,
	}
}

// SetAccrediter allows to set a given user as being accreditated by the given accrediter.
func (keeper Keeper) SetAccrediter(ctx sdk.Context, accrediter sdk.AccAddress, user sdk.AccAddress) error {
	store := ctx.KVStore(keeper.StoreKey)
	if store.Has(user) {
		return errors.New("user already have an accrediter")
	}

	// Build the accreditation
	accreditation := types.Accreditation{
		Accrediter: accrediter,
		User:       user,
		Rewarded:   false,
	}

	// Save the accreditation
	accreditationBz := keeper.cdc.MustMarshalBinaryBare(&accreditation)
	store.Set(user, accreditationBz)
	return nil
}

// GetAccrediter allows to get the accrediter a given user
func (keeper Keeper) GetAccrediter(ctx sdk.Context, user sdk.AccAddress) sdk.AccAddress {
	store := ctx.KVStore(keeper.StoreKey)

	// Check existence
	if !store.Has(user) {
		return nil
	}

	var accreditation types.Accreditation
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(user), &accreditation)
	return accreditation.Accrediter
}

// DistributeReward allows to give the specified accrediter the specified amount of reward related
// to the accreditation of the specified user
func (keeper Keeper) DistributeReward(ctx sdk.Context, accrediter sdk.AccAddress, reward sdk.Coins, user sdk.AccAddress) error {
	store := ctx.KVStore(keeper.StoreKey)
	if !store.Has(user) {
		return errors.New("user does not have an accrediter")
	}

	if reward.IsAnyNegative() {
		return errors.New("rewards cannot be negative")
	}

	// Get the liquidity pool value
	var liquidity sdk.Coins
	keeper.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.LiquidityPoolKey)), &liquidity)
	if reward == nil || reward.IsAnyGT(liquidity) {
		return errors.New("liquidity pool has not a sufficient amount of tokens for this reward")
	}

	// Decrement pool and send the amount to the accrediter
	liquidity = liquidity.Sub(reward)
	if _, err := keeper.BankKeeper.AddCoins(ctx, accrediter, reward); err != nil {
		return err
	}

	// Save the updated pool
	store.Set([]byte(types.LiquidityPoolKey), keeper.cdc.MustMarshalBinaryBare(&liquidity))

	// Update the accreditation
	var accreditation types.Accreditation
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(user), &accreditation)
	accreditation.Rewarded = true
	store.Set(user, keeper.cdc.MustMarshalBinaryBare(&accreditation))

	return nil
}

// DepositIntoPool allows anyone to deposit into the liquidity pool that
// will be used when giving out rewards for accreditations.
func (keeper Keeper) DepositIntoPool(ctx sdk.Context, amount sdk.Coins) error {
	if amount.IsAnyNegative() {
		return errors.New("amount cannot be negative")
	}

	store := ctx.KVStore(keeper.StoreKey)

	// Add the amount to the pool
	var pool sdk.Coins
	keeper.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.LiquidityPoolKey)), &pool)
	pool = pool.Add(amount)
	store.Set([]byte(types.LiquidityPoolKey), keeper.cdc.MustMarshalBinaryBare(&pool))

	return nil
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

// AddTrustworthySigner allows to add the given signer as a trustworthy entity
// that can sign transactions setting an accrediter for a user.
func (keeper Keeper) AddTrustworthySigner(ctx sdk.Context, signer sdk.AccAddress) {
	store := ctx.KVStore(keeper.StoreKey)

	signers := utypes.Addresses(keeper.GetTrustworthySigners(ctx))
	signers = signers.AppendIfMissing(signer)

	newSignersBz := keeper.cdc.MustMarshalBinaryBare(&signers)
	store.Set([]byte(types.TrustworthySignersKey), newSignersBz)
}

// GetAccreditations returns all the accreditations that have been
func (keeper Keeper) GetAccreditations(ctx sdk.Context) (accreditations []types.Accreditation) {
	store := ctx.KVStore(keeper.StoreKey)
	iterator := store.Iterator(nil, nil)

	for ; iterator.Valid() && string(iterator.Key()) != types.TrustworthySignersKey; iterator.Next() {
		var accreditation types.Accreditation
		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &accreditation)
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
	keeper.cdc.MustUnmarshalBinaryBare(signersBz, &signers)

	return
}
