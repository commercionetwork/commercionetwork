package keeper

import (
	"errors"

	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/types"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
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
func NewKeeper(storeKey sdk.StoreKey, bankKeeper bank.Keeper, cdc *codec.Codec) Keeper {
	return Keeper{
		StoreKey:   storeKey,
		BankKeeper: bankKeeper,
		cdc:        cdc,
	}
}

// SetAccrediter allows to set a given user as being accreditated by the given accrediter.
func (keeper Keeper) SetAccrediter(ctx sdk.Context, user, accrediter sdk.AccAddress) error {
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
	accreditationBz := keeper.cdc.MustMarshalBinaryBare(accreditation)
	store.Set(user, accreditationBz)
	return nil
}

// GetAccreditation allows to get the accrediter a given user
func (keeper Keeper) GetAccreditation(ctx sdk.Context, user sdk.AccAddress) types.Accreditation {
	store := ctx.KVStore(keeper.StoreKey)

	var accreditation types.Accreditation
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(user), &accreditation)
	return accreditation
}

// GetAccreditations returns all the accreditations that have been
func (keeper Keeper) GetAccreditations(ctx sdk.Context) (accreditations []types.Accreditation) {
	store := ctx.KVStore(keeper.StoreKey)
	iterator := store.Iterator(nil, nil)

	for ; iterator.Valid() && string(iterator.Key()) != types.TrustedSignersStoreKey; iterator.Next() {
		var accreditation types.Accreditation
		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &accreditation)
		accreditations = append(accreditations, accreditation)
	}

	return
}

// DepositIntoPool allows anyone to deposit into the liquidity pool that
// will be used when giving out rewards for accreditations.
func (keeper Keeper) DepositIntoPool(ctx sdk.Context, depositor sdk.AccAddress, amount sdk.Coins) error {
	if amount.IsAnyNegative() {
		return errors.New("amount cannot be negative")
	}

	store := ctx.KVStore(keeper.StoreKey)

	// Remove the coins from the user wallet
	_, err := keeper.BankKeeper.SubtractCoins(ctx, depositor, amount)
	if err != nil {
		return err
	}

	// Add the amount to the pool
	var pool sdk.Coins
	keeper.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.LiquidityPoolKey)), &pool)
	pool = pool.Add(amount)
	store.Set([]byte(types.LiquidityPoolKey), keeper.cdc.MustMarshalBinaryBare(&pool))

	return nil
}

// SetPoolFunds allows to set the current pool funds amount
func (keeper Keeper) SetPoolFunds(ctx sdk.Context, pool sdk.Coins) {
	store := ctx.KVStore(keeper.StoreKey)
	store.Set([]byte(types.LiquidityPoolKey), keeper.cdc.MustMarshalBinaryBare(&pool))
}

// GetPoolFunds return the current pool funds for the given context
func (keeper Keeper) GetPoolFunds(ctx sdk.Context) sdk.Coins {
	store := ctx.KVStore(keeper.StoreKey)
	var pool sdk.Coins
	keeper.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.LiquidityPoolKey)), &pool)
	return pool
}

// DistributeReward allows to give the specified accrediter the specified amount of reward related
// to the accreditation of the specified user
func (keeper Keeper) DistributeReward(ctx sdk.Context, accrediter sdk.AccAddress, reward sdk.Coins, user sdk.AccAddress) error {
	if reward == nil {
		return errors.New("reward cannot be empty")
	}

	if reward.IsAnyNegative() {
		return errors.New("rewards cannot be negative")
	}

	store := ctx.KVStore(keeper.StoreKey)
	if !store.Has(user) {
		return errors.New("user does not have an accrediter")
	}

	var accreditation types.Accreditation
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(user), &accreditation)

	if accreditation.Rewarded {
		return errors.New("the accrediter has already been rewarded for this user")
	}

	// Get the liquidity pool value
	var liquidity sdk.Coins
	keeper.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.LiquidityPoolKey)), &liquidity)
	if liquidity == nil || reward.IsAnyGT(liquidity) {
		return errors.New("liquidity pool has not a sufficient amount of tokens for this reward")
	}

	// Decrement pool and send the amount to the accrediter
	liquidity = liquidity.Sub(reward)
	if _, err := keeper.BankKeeper.AddCoins(ctx, accrediter, reward); err != nil {
		return err
	}

	// Save the updated pool
	if liquidity.Empty() {

	} else {
		store.Set([]byte(types.LiquidityPoolKey), keeper.cdc.MustMarshalBinaryBare(&liquidity))
	}

	// Update the accreditation
	accreditation.Rewarded = true
	store.Set(user, keeper.cdc.MustMarshalBinaryBare(&accreditation))

	return nil
}

// AddTrustedSigner allows to add the given signer as a trusted entity
// that can sign transactions setting an accrediter for a user.
func (keeper Keeper) AddTrustedSigner(ctx sdk.Context, signer sdk.AccAddress) {
	store := ctx.KVStore(keeper.StoreKey)

	signers := ctypes.Addresses(keeper.GetTrustedSigners(ctx))
	signers = signers.AppendIfMissing(signer)

	newSignersBz := keeper.cdc.MustMarshalBinaryBare(&signers)
	store.Set([]byte(types.TrustedSignersStoreKey), newSignersBz)
}

// GetTrustedSigners returns the list of signers that are allowed to sign
// transactions setting a specific accrediter for a user.
// NOTE. Any user which is not present inside the returned list SHOULD NOT
// be allowed to send a transaction setting an accrediter for another user.
func (keeper Keeper) GetTrustedSigners(ctx sdk.Context) (signers ctypes.Addresses) {
	store := ctx.KVStore(keeper.StoreKey)

	signersBz := store.Get([]byte(types.TrustedSignersStoreKey))
	keeper.cdc.MustUnmarshalBinaryBare(signersBz, &signers)

	return
}

// IsTrustedSigner tells if the given signer is a trusted one or not
func (keeper Keeper) IsTrustedSigner(ctx sdk.Context, signer sdk.AccAddress) bool {
	signers := keeper.GetTrustedSigners(ctx)
	return signers.Contains(signer)
}
