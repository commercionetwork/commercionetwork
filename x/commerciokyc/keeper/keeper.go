package keeper

import (
	"fmt"
	"strconv"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	government "github.com/commercionetwork/commercionetwork/x/government/keeper"
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

var (
	// DPY is Days Per Year
	DPY = sdk.NewDecWithPrec(36525, 2)
	// HPD is Hours Per Day
	HPD = sdk.NewDecWithPrec(24, 0)
	// MPH is Minutes Per Hour
	MPH = sdk.NewDecWithPrec(60, 0)
	// BPM Blocks Per Minutes 7 secs x Block
	BPM = sdk.NewDecWithPrec(9, 0)
	// BPY is Blocks Per Year
	BPY = DPY.Mul(HPD).Mul(MPH).Mul(BPM)
)

// NewKeeper creates new instances of the accreditation module Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, supplyKeeper supply.Keeper, sendKeeper bank.Keeper, governmentKeeper government.Keeper, accountKeeper auth.AccountKeeper) Keeper {

	// ensure commerciomint module account is set
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

// BuyMembership allow to commerciomint and assign a membership of the given membershipType to the specified user.
// If the user already has a membership assigned, deletes the current one and assigns to it the new one.
func (k Keeper) BuyMembership(ctx sdk.Context, buyer sdk.AccAddress, membershipType string, tsp sdk.AccAddress, height int64) error {
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
	return k.AssignMembership(ctx, buyer, membershipType, tsp, height)
}

// AssignMembership allow to commerciomint and assign a membership of the given membershipType to the specified user.
// If the user already has a membership assigned, deletes the current one and assigns to it the new one.
// Returns the URI of the new minted token represented the assigned membership, or an error if something goes w
// THIS FUNCTION CAN TRANSFORM WITH AssignMembership(ctx sdk.Context, m types.Membership)
func (k Keeper) AssignMembership(ctx sdk.Context, user sdk.AccAddress, membershipType string, tsp sdk.AccAddress, height int64) error {
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
	if height <= 0 {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid expiry height: %s", strconv.FormatInt(height, 10)))
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

	//height := ctx.BlockHeight() + (365 * 24 * 60 * 60/7)
	membership := types.NewMembership(membershipType, user, tsp, height)

	// Save membership
	store.Set(staddr, k.Cdc.MustMarshalBinaryBare(&membership))

	return nil
}

// ComputeExpiryHeight compute expiry height of membership.
func (k Keeper) ComputeExpiryHeight(blockHeight int64) int64 {
	blocksOfYear := DPY.Mul(HPD).Mul(MPH).Mul(BPM)
	return sdk.NewDec(blockHeight).Add(blocksOfYear).TruncateInt64()
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
	return nil
}

// MembershipIterator returns an Iterator for all the memberships stored.
func (k Keeper) MembershipIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.StoreKey)
	return sdk.KVStorePrefixIterator(store, []byte(types.MembershipsStorageKey))
}

// ExtractMembership extracts a Membership struct from retrieved value
func (k Keeper) ExtractMembership(value []byte) types.Membership {
	var m types.Membership
	k.Cdc.MustUnmarshalBinaryBare(value, &m)
	return m
}

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
func (k Keeper) ExportMemberships(ctx sdk.Context, height int64) types.Memberships {
	im := k.MembershipIterator(ctx)
	m := types.Membership{}
	ms := types.Memberships{}
	defer im.Close()
	for ; im.Valid(); im.Next() {
		k.Cdc.MustUnmarshalBinaryBare(im.Value(), &m)
		m.ExpiryAt = m.ExpiryAt - height
		if m.ExpiryAt <= 0 && m.MembershipType != types.MembershipTypeBlack {
			continue
		}
		ms = append(ms, m)
	}

	return ms
}

// RemoveExpiredMemberships delete all expired memberships
func (k Keeper) RemoveExpiredMemberships(ctx sdk.Context) error {
	blockHeight := ctx.BlockHeight()
	if blockHeight == 0 {
		blockHeight = 1
	}
	for _, m := range k.GetMemberships(ctx) {
		h := m.ExpiryAt - blockHeight
		if h <= 0 {
			if m.MembershipType == types.MembershipTypeBlack {
				height := k.ComputeExpiryHeight(ctx.BlockHeight())
				membership := types.NewMembership(types.MembershipTypeBlack, m.Owner, m.TspAddress, height)
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
