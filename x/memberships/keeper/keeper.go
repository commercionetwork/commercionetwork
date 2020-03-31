package keeper

import (
	"fmt"

	governmentKeeper "github.com/commercionetwork/commercionetwork/x/government/keeper"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/x/auth"

	creditrisk "github.com/commercionetwork/commercionetwork/x/creditrisk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/commercionetwork/commercionetwork/x/memberships/types"
)

var membershipCosts = map[string]int64{
	types.MembershipTypeBronze: 25,
	types.MembershipTypeSilver: 250,
	types.MembershipTypeGold:   2500,
	types.MembershipTypeBlack:  50000,
}

type Keeper struct {
	Cdc              *codec.Codec
	StoreKey         sdk.StoreKey
	SupplyKeeper     supply.Keeper
	governmentKeeper governmentKeeper.Keeper
	accountKeeper    auth.AccountKeeper
}

// NewKeeper creates new instances of the accreditation module Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, supplyKeeper supply.Keeper, governmentKeeper governmentKeeper.Keeper, accountKeeper auth.AccountKeeper) Keeper {

	// ensure commerciomint module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		Cdc:              cdc,
		StoreKey:         storeKey,
		SupplyKeeper:     supplyKeeper,
		governmentKeeper: governmentKeeper,
		accountKeeper:    accountKeeper,
	}
}

// storageForAddr returns a string representing the KVStore storage key for an addr.
func (k Keeper) storageForAddr(addr sdk.AccAddress) []byte {
	return append([]byte(types.MembershipsStorageKey), k.Cdc.MustMarshalBinaryBare(addr)...)
}

// BuyMembership allow to commerciomint and assign a membership of the given membershipType to the specified user.
// If the user already has a membership assigned, deletes the current one and assigns to it the new one.
func (k Keeper) BuyMembership(ctx sdk.Context, buyer sdk.AccAddress, membershipType string) error {
	if membershipType == types.MembershipTypeBlack {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, "cannot buy black membership")
	}

	// Transfer the buyer's tokens to the credit risk pool
	membershipPrice := membershipCosts[membershipType] * 1000000 // Always multiply by one million
	membershipCost := sdk.NewCoins(sdk.NewInt64Coin(k.GetStableCreditsDenom(ctx), membershipPrice))
	if err := k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, buyer, creditrisk.ModuleName, membershipCost); err != nil {
		return err
	}

	// Assign the membership
	return k.AssignMembership(ctx, buyer, membershipType)
}

// AssignMembership allow to commerciomint and assign a membership of the given membershipType to the specified user.
// If the user already has a membership assigned, deletes the current one and assigns to it the new one.
// Returns the URI of the new minted token represented the assigned membership, or an error if something goes w
func (k Keeper) AssignMembership(ctx sdk.Context, user sdk.AccAddress, membershipType string) error {
	// Check the membership type validity
	if !types.IsMembershipTypeValid(membershipType) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid membership type: %s", membershipType))
	}

	_ = k.RemoveMembership(ctx, user)

	store := ctx.KVStore(k.StoreKey)

	staddr := k.storageForAddr(user)
	if store.Has(staddr) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest,
			fmt.Sprintf(
				"cannot add membership \"%s\" for address %s: user already have a membership",
				membershipType,
				user.String(),
			),
		)
	}

	store.Set(staddr, k.Cdc.MustMarshalBinaryBare(membershipType))

	return nil
}

// GetMembership allows to retrieve any existent membership for the specified user.
// The second returned false (the boolean one) tells if the NFT token representing the membership was found or not
func (k Keeper) GetMembership(ctx sdk.Context, user sdk.AccAddress) (types.Membership, error) {
	store := ctx.KVStore(k.StoreKey)

	if !store.Has(k.storageForAddr(user)) {
		return types.Membership{}, sdkErr.Wrap(sdkErr.ErrUnknownRequest,
			fmt.Sprintf("membership not found for user \"%s\"", user.String()),
		)
	}

	membershipRaw := store.Get(k.storageForAddr(user))
	var ms types.Membership
	k.Cdc.MustUnmarshalBinaryBare(membershipRaw, &ms.MembershipType)
	ms.Owner = user

	return ms, nil
}

// RemoveMembership allows to remove any existing membership associated with the given user.
func (k Keeper) RemoveMembership(ctx sdk.Context, user sdk.AccAddress) error {
	store := ctx.KVStore(k.StoreKey)

	if !store.Has(k.storageForAddr(user)) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest,
			fmt.Sprintf("account \"%s\" does not have any membership", user.String()),
		)
	}

	store.Delete(k.storageForAddr(user))

	return nil
}

// GetMembershipIterator returns an Iterator for all the memberships stored.
func (k Keeper) MembershipIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.StoreKey)

	return sdk.KVStorePrefixIterator(store, []byte(types.MembershipsStorageKey))
}

// ExtractMembership extracts a Membership struct from key and value retrieved
// from MembershipIterator
func (k Keeper) ExtractMembership(key []byte, value []byte) types.Membership {
	rawAddr := key[len(types.MembershipsStorageKey):]

	var addr sdk.AccAddress
	var m string

	k.Cdc.MustUnmarshalBinaryBare(rawAddr, &addr)
	k.Cdc.MustUnmarshalBinaryBare(value, &m)

	return types.Membership{
		Owner:          addr,
		MembershipType: m,
	}

}

// GetStableCreditsDenom returns the denom that must be used when referring to stable credits
// that can be used to purchase a membership
func (k Keeper) GetStableCreditsDenom(ctx sdk.Context) (denom string) {
	store := ctx.KVStore(k.StoreKey)
	return string(store.Get([]byte(types.StableCreditsStoreKey)))
}

// SetStableCreditsDenom allows to set the denom of the coins that must be used as stable credits
// when purchasing a membership.
func (k Keeper) SetStableCreditsDenom(ctx sdk.Context, denom string) {
	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.StableCreditsStoreKey), []byte(denom))
}
