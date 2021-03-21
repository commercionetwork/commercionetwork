package keeper

import (
	"fmt"
	"time"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	government "github.com/commercionetwork/commercionetwork/x/government/keeper"
)

const (
	//eventBuyMembership    = "buy_membership"
	eventAssignMembership = "assign_membership"
	eventRemoveMembership = "remove_membership"

	secondsPerYear time.Duration = time.Hour * 24 * 365
)

var membershipCosts = map[string]int64{
	types.MembershipTypeGreen:  5,
	types.MembershipTypeBronze: 25,
	types.MembershipTypeSilver: 250,
	types.MembershipTypeGold:   2500,
	types.MembershipTypeBlack:  50000,
}

type Keeper struct {
	Cdc              *codec.Codec
	StoreKey         sdk.StoreKey
	SupplyKeeper     supply.Keeper
	BankKeeper       bank.Keeper
	governmentKeeper government.Keeper
	accountKeeper    auth.AccountKeeper
}

// NewKeeper creates new instances of the accreditation module Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, supplyKeeper supply.Keeper, sendKeeper bank.Keeper, governmentKeeper government.Keeper, accountKeeper auth.AccountKeeper) Keeper {

	// ensure commerciokyc module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		Cdc:              cdc,
		StoreKey:         storeKey,
		SupplyKeeper:     supplyKeeper,
		BankKeeper:       sendKeeper,
		governmentKeeper: governmentKeeper,
		accountKeeper:    accountKeeper,
	}
}

// storageForAddr returns a string representing the KVStore storage key for an addr.
func (k Keeper) storageForAddr(addr sdk.AccAddress) []byte {
	return append([]byte(types.MembershipsStorageKey), k.Cdc.MustMarshalBinaryBare(addr)...)
}

// BuyMembership allow to commerciokyc and assign a membership of the given membershipType to the specified user.
// If the user already has a membership assigned, deletes the current one and assigns to it the new one.
func (k Keeper) BuyMembership(ctx sdk.Context, buyer sdk.AccAddress, membershipType string, tsp sdk.AccAddress, expited_at time.Time) error {
	if membershipType == types.MembershipTypeBlack {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, "cannot buy black membership")
	}

	// Transfer the tsp tokens to government
	membershipPrice := membershipCosts[membershipType] * 1000000 // Always multiply by one million
	membershipCost := sdk.NewCoins(sdk.NewInt64Coin("uccc", membershipPrice))
	govAddr := k.governmentKeeper.GetGovernmentAddress(ctx)
	if err := k.BankKeeper.SendCoins(ctx, tsp, govAddr, membershipCost); err != nil {
		return err
	}

	// Assign the membership
	return k.AssignMembership(ctx, buyer, membershipType, tsp, expited_at)
}

// AssignMembership allow to assign a membership of the given membershipType to the specified user with tsp and expired height.
// TODO maybe it's better to pass membership object to function
func (k Keeper) AssignMembership(ctx sdk.Context, user sdk.AccAddress, membershipType string, tsp sdk.AccAddress, expited_at time.Time) error {
	// Check the membership type validity.
	if !types.IsMembershipTypeValid(membershipType) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid membership type: %s", membershipType))
	}

	// TODO resolve problems in init genesis to remove membershipType != types.MembershipTypeBlack
	if k.IsTrustedServiceProvider(ctx, user) && membershipType != types.MembershipTypeBlack {
		return sdkErr.Wrap(sdkErr.ErrUnauthorized,
			fmt.Sprintf("account \"%s\" is a Trust Service Provider: remove from tsps list before", user.String()),
		)
	}

	// Check if height is greater then zero
	if expited_at.Before(time.Now()) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid expiry date: %s", expited_at))
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

	//expited_at := ctx.BlockTime() + (365 * 24 * 60 * 60) seconds
	membership := types.NewMembership(membershipType, user, tsp, expited_at.UTC())

	// Save membership
	store.Set(staddr, k.Cdc.MustMarshalBinaryBare(&membership))
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventAssignMembership,
		sdk.NewAttribute("owner", membership.Owner.String()),
		sdk.NewAttribute("membership_type", membership.MembershipType),
		sdk.NewAttribute("tsp_address", membership.TspAddress.String()),
		sdk.NewAttribute("expiry_at", membership.ExpiryAt.String()),
	))

	return nil
}

// ComputeExpiryHeight compute expiry height of membership.
func (k Keeper) ComputeExpiryHeight(blockTime time.Time) time.Time {
	expirationAt := blockTime.Add(secondsPerYear)
	return expirationAt
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
	k.Cdc.MustUnmarshalBinaryBare(membershipRaw, &ms)
	return ms, nil
}

// RemoveMembership allows to remove any existing membership associated with the given user.
func (k Keeper) RemoveMembership(ctx sdk.Context, user sdk.AccAddress) error {
	store := ctx.KVStore(k.StoreKey)

	if k.IsTrustedServiceProvider(ctx, user) {
		return sdkErr.Wrap(sdkErr.ErrUnauthorized,
			fmt.Sprintf("account \"%s\" is a Trust Service Provider: remove from tsps list before", user.String()),
		)
	}

	if !store.Has(k.storageForAddr(user)) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest,
			fmt.Sprintf("account \"%s\" does not have any membership", user.String()),
		)
	}

	store.Delete(k.storageForAddr(user))

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventRemoveMembership,
		sdk.NewAttribute("subscriber", user.String()),
	))

	return nil
}

// MembershipIterator returns an Iterator for all the memberships stored.
func (k Keeper) MembershipIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.StoreKey)
	return sdk.KVStorePrefixIterator(store, []byte(types.MembershipsStorageKey))
}

// ExtractMembership extracts a Membership struct from retrieved value
/*func (k Keeper) ExtractMembership(value []byte) types.Membership {
	var m types.Membership
	k.Cdc.MustUnmarshalBinaryBare(value, &m)
	return m
}*/

// GetMemberships extracts all memerships
func (k Keeper) GetMemberships(ctx sdk.Context) types.Memberships {
	im := k.MembershipIterator(ctx)
	ms := types.Memberships{}
	defer im.Close()
	for ; im.Valid(); im.Next() {
		var m types.Membership
		k.Cdc.MustUnmarshalBinaryBare(im.Value(), &m)
		ms = append(ms, m)
	}

	return ms
}

// GetTspMemberships extracts all memerships
func (k Keeper) GetTspMemberships(ctx sdk.Context, tsp sdk.Address) types.Memberships {
	im := k.MembershipIterator(ctx)
	m := types.Membership{}
	ms := types.Memberships{}
	defer im.Close()
	for ; im.Valid(); im.Next() {
		k.Cdc.MustUnmarshalBinaryBare(im.Value(), &m)
		if !m.TspAddress.Equals(tsp) {
			continue
		}
		ms = append(ms, m)
	}

	return ms
}

// ExportMemberships extracts all memberships for export
func (k Keeper) ExportMemberships(ctx sdk.Context) types.Memberships {
	im := k.MembershipIterator(ctx)
	m := types.Membership{}
	ms := types.Memberships{}
	defer im.Close()
	for ; im.Valid(); im.Next() {
		k.Cdc.MustUnmarshalBinaryBare(im.Value(), &m)
		ms = append(ms, m)
	}
	return ms
}

// RemoveExpiredMemberships delete all expired memberships
func (k Keeper) RemoveExpiredMemberships(ctx sdk.Context) error {
	blockTime := ctx.BlockTime()
	for _, m := range k.GetMemberships(ctx) {
		if blockTime.After(m.ExpiryAt) {
			if m.MembershipType == types.MembershipTypeBlack {
				expiredAt := k.ComputeExpiryHeight(ctx.BlockTime())
				membership := types.NewMembership(types.MembershipTypeBlack, m.Owner, m.TspAddress, expiredAt)
				store := ctx.KVStore(k.StoreKey)
				staddr := k.storageForAddr(m.Owner)
				store.Set(staddr, k.Cdc.MustMarshalBinaryBare(&membership))
			} else {
				err := k.RemoveMembership(ctx, m.Owner)
				if err != nil {
					panic(err)
				}
			}
		}
	}
	return nil
}
