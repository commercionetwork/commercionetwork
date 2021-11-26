package keeper

import (
	"fmt"
	"time"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	accTypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	//kmint "github.com/commercionetwork/commercionetwork/x/commerciomint/keeper"
)

const (
	stakeDenom        = "ucommercio"
	stableCreditDenom = "uccc"
)

var membershipRewards = map[string]map[string]sdk.Dec{
	types.MembershipTypeGreen: {
		types.MembershipTypeGreen:  sdk.NewDecWithPrec(1, 2),   // 1% of 1
		types.MembershipTypeBronze: sdk.NewDecWithPrec(1, 1),   // 2% of 5
		types.MembershipTypeSilver: sdk.NewDecWithPrec(15, 1),  // 3% of 50
		types.MembershipTypeGold:   sdk.NewDecWithPrec(20, 0),  // 4% of 500
		types.MembershipTypeBlack:  sdk.NewDecWithPrec(250, 0), // 2.5% of 10000
	},
	types.MembershipTypeBronze: {
		types.MembershipTypeGreen:  sdk.NewDecWithPrec(25, 3),   // 2.5% of 1
		types.MembershipTypeBronze: sdk.NewDecWithPrec(125, 2),  // 5% of 25
		types.MembershipTypeSilver: sdk.NewDecWithPrec(25, 0),   // 10% of 250
		types.MembershipTypeGold:   sdk.NewDecWithPrec(375, 0),  // 15% of 2500
		types.MembershipTypeBlack:  sdk.NewDecWithPrec(5000, 0), // 10% of 50000
	},
	types.MembershipTypeSilver: {
		types.MembershipTypeGreen:  sdk.NewDecWithPrec(1, 1),     // 1% of 1
		types.MembershipTypeBronze: sdk.NewDecWithPrec(5, 0),     // 20% of 25
		types.MembershipTypeSilver: sdk.NewDecWithPrec(75, 0),    // 30% of 250
		types.MembershipTypeGold:   sdk.NewDecWithPrec(1000, 0),  // 40% of 2500
		types.MembershipTypeBlack:  sdk.NewDecWithPrec(12500, 0), // 12.5% of 50000
	},
	types.MembershipTypeGold: {
		types.MembershipTypeGreen:  sdk.NewDecWithPrec(4, 1),     // 40% of 1
		types.MembershipTypeBronze: sdk.NewDecWithPrec(125, 1),   // 50% of 25
		types.MembershipTypeSilver: sdk.NewDecWithPrec(150, 0),   // 60% of 250
		types.MembershipTypeGold:   sdk.NewDecWithPrec(1750, 0),  // 70% of 2500
		types.MembershipTypeBlack:  sdk.NewDecWithPrec(20000, 0), // 40% of 50000
	},
	types.MembershipTypeBlack: {
		types.MembershipTypeGreen:  sdk.NewDecWithPrec(5, 1),     // 50% of 1
		types.MembershipTypeBronze: sdk.NewDecWithPrec(175, 2),   // 70% of 25
		types.MembershipTypeSilver: sdk.NewDecWithPrec(200, 0),   // 80% of 250
		types.MembershipTypeGold:   sdk.NewDecWithPrec(2250, 0),  // 90% of 2500
		types.MembershipTypeBlack:  sdk.NewDecWithPrec(25000, 0), // 50% of 50000
	},
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
			fmt.Sprintf("account \"%s\" is a Trust Service Provider: remove from tsps list before", user),
		)
	}

	// Check if height is greater then zero
	if expited_at.Before(time.Now()) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid expiry date: %s", expited_at))
	}

	_ = k.DeleteMembership(ctx, user)

	store := ctx.KVStore(k.StoreKey)

	staddr := k.storageForAddr(user)
	if store.Has(staddr) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest,
			fmt.Sprintf(
				"cannot add membership \"%s\" for address %s: user already have a membership",
				membershipType,
				user,
			),
		)
	}

	//expited_at := ctx.BlockTime() + (365 * 24 * 60 * 60) seconds
	//membership := types.NewMembership(membershipType, user, tsp, expited_at.UTC())

	// Save membership
	membership := types.NewMembership(membershipType, user, tsp, expited_at.UTC())
	store.Set(staddr, k.Cdc.MustMarshalBinaryBare(&membership))
	// TODO emits events
	/*ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventAssignMembership,
		sdk.NewAttribute("owner", membership.Owner.String()),
		sdk.NewAttribute("membership_type", membership.MembershipType),
		sdk.NewAttribute("tsp_address", membership.TspAddress.String()),
		sdk.NewAttribute("expiry_at", membership.ExpiryAt.String()),
	))*/

	return nil
}

// DeleteMembership allows to remove any existing membership associated with the given user.
func (k Keeper) DeleteMembership(ctx sdk.Context, user sdk.AccAddress) error {
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
	/*
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			eventRemoveMembership,
			sdk.NewAttribute("subscriber", user.String()),
		))*/

	return nil
}

// DistributeReward allows to distribute the rewards to the sender of the specified invite upon the receiver has
// properly bought a membership of the given membershipType
func (k Keeper) DistributeReward(ctx sdk.Context, invite types.Invite) error {
	// the invite we got is either invalid or already rewarded, get out!
	inviteStatus := types.InviteStatus(invite.Status)

	if inviteStatus == types.InviteStatusRewarded || inviteStatus == types.InviteStatusInvalid {
		return nil
	}
	// Calculate reward for invite
	inviteSender, _ := sdk.AccAddressFromBech32(invite.Sender)
	_, err := k.GetMembership(ctx, inviteSender)
	if err != nil || invite.SenderMembership == "" {
		return sdkErr.Wrap(sdkErr.ErrUnauthorized, "Invite sender does not have a membership")
	}

	inviteUser, _ := sdk.AccAddressFromBech32(invite.User)
	recipientMembership, err := k.GetMembership(ctx, inviteUser)
	if err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnauthorized, "Invite recipient does not have a membership")
	}

	senderMembershipType := invite.SenderMembership
	recipientMembershipType := recipientMembership.MembershipType

	// Get the reward amount by searching up inside the matrix.
	// Multiply the found amount by 1.000.000 as coins are represented as millionth of units, and make it an int
	/*
		var rewardCrossValue sdk.Dec
		var ok bool
		if rewardCrossValue, ok = membershipRewards[senderMembershipType][recipientMembershipType]; !ok {
			return sdkErr.Wrap(sdkErr.ErrUnauthorized, "Invalid reward options")
		}
		rewardAmount := rewardCrossValue.MulInt64(1000000).TruncateInt()
	*/
	rewardAmount := membershipRewards[senderMembershipType][recipientMembershipType].MulInt64(1000000).TruncateInt()

	// Get the pool amount
	poolAmount := k.GetPoolFunds(ctx).AmountOf(stakeDenom)

	// Distribute the reward taking it from the pool amount
	if poolAmount.GT(sdk.ZeroInt()) {

		// If the reward is more than the current pool amount, set the reward as the total pool amount
		if rewardAmount.GT(poolAmount) {
			rewardAmount = poolAmount
		}
		// Calcute equivalent distribution in uccc
		//mintk := kmint.NewKeeper(k.Cdc, k.StoreKey,)

		//ucccConversionRate := mintk.GetConversionRate(ctx)
		//kmintTypes.GetConv

		rewardCoins := sdk.NewCoins(sdk.NewCoin(stakeDenom, rewardAmount))

		// Send the reward to the invite sender
		inviteSender, _ := sdk.AccAddressFromBech32(invite.Sender)
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, inviteSender, rewardCoins); err != nil {
			return err
		}
		// TODO  emits events
		/*
			ctx.EventManager().EmitEvent(sdk.NewEvent(
				eventDistributeReward,
				sdk.NewAttribute("invite_sender", invite.Sender.String()),
				sdk.NewAttribute("reward_coins", rewardCoins.String()),
				sdk.NewAttribute("sender_membership_type", senderMembershipType),
				sdk.NewAttribute("recipient_membership_type", recipientMembership.MembershipType),
				sdk.NewAttribute("invite_recipient", invite.User.String()),
			))
		*/

	}

	// Set the invitation as rewarded
	newInvite := types.Invite{
		Sender:           invite.Sender,
		User:             invite.User,
		SenderMembership: invite.SenderMembership,
		Status:           uint64(types.InviteStatusRewarded), // TODO control conversion
	}

	k.SaveInvite(ctx, newInvite)

	return nil
}

// GetTrustedServiceProviders returns the list of signers that are allowed to sign
// transactions setting a specific accrediter for a user.
// NOTE. Any user which is not present inside the returned list SHOULD NOT
// be allowed to send a transaction setting an accrediter for another user.
func (k Keeper) GetTrustedServiceProviders(ctx sdk.Context) (signers types.TrustedServiceProviders) {
	store := ctx.KVStore(k.StoreKey)

	signersBz := store.Get([]byte(types.TrustedSignersStoreKey))
	k.Cdc.UnmarshalBinaryBare(signersBz, &signers)

	//k.Cdc.MustUnmarshalBinaryBare(signersBz, &signers)
	// Cannot use add govAddress: trust service provider doesn't work proprerly
	//signers = append(signers, k.governmentKeeper.GetGovernmentAddress(ctx))
	return signers
}

// IsTrustedServiceProvider tells if the given signer is a trusted one or not
func (k Keeper) IsTrustedServiceProvider(ctx sdk.Context, signer sdk.Address) bool {
	var signers ctypes.Strings
	signers = k.GetTrustedServiceProviders(ctx).Addresses
	return signers.Contains(signer.String()) || signer.Equals(k.govKeeper.GetGovernmentAddress(ctx))
}

// TspIterator returns an Iterator for all the tsps stored.
func (k Keeper) TspIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.StoreKey)
	return sdk.KVStorePrefixIterator(store, []byte(types.TrustedSignersStoreKey))
}

// storageForAddr returns a string representing the KVStore storage key for an addr.
func (k Keeper) storageForAddr(addr sdk.AccAddress) []byte {
	//return append([]byte(types.MembershipsStorageKey), k.Cdc.MustMarshalBinaryBare(&addr)...)
	return append([]byte(types.MembershipsStorageKey), addr.Bytes()...)
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

// GetMemberships extracts all memerships
func (k Keeper) GetMemberships(ctx sdk.Context) []*types.Membership {
	im := k.MembershipIterator(ctx)
	ms := []*types.Membership{}
	defer im.Close()
	for ; im.Valid(); im.Next() {
		var m types.Membership
		k.Cdc.MustUnmarshalBinaryBare(im.Value(), &m)
		ms = append(ms, &m)
	}

	return ms
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

// GetMembershipModuleAccount returns the module account for the commerciokyc module
func (k Keeper) GetMembershipModuleAccount(ctx sdk.Context) accTypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// MembershipIterator returns an Iterator for all the memberships stored.
func (k Keeper) MembershipIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.StoreKey)
	return sdk.KVStorePrefixIterator(store, []byte(types.MembershipsStorageKey))
}

// ComputeExpiryHeight compute expiry height of membership.
func (k Keeper) ComputeExpiryHeight(blockTime time.Time) time.Time {
	expirationAt := blockTime.Add(SecondsPerYear)
	return expirationAt
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

	/*
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			eventRemoveMembership,
			sdk.NewAttribute("subscriber", user.String()),
		))
	*/

	return nil
}

// GetTspMemberships extracts all memerships
func (k Keeper) GetTspMemberships(ctx sdk.Context, tsp sdk.Address) types.Memberships {
	im := k.MembershipIterator(ctx)
	m := types.Membership{}
	ms := types.Memberships{}
	defer im.Close()
	for ; im.Valid(); im.Next() {
		k.Cdc.MustUnmarshalBinaryBare(im.Value(), &m)
		if m.TspAddress != tsp.String() {
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
		if blockTime.After(*m.ExpiryAt) {
			mOwner, _ := sdk.AccAddressFromBech32(m.Owner)
			mTspAddress, _ := sdk.AccAddressFromBech32(m.TspAddress)
			if m.MembershipType == types.MembershipTypeBlack {
				expiredAt := k.ComputeExpiryHeight(ctx.BlockTime())
				membership := types.NewMembership(types.MembershipTypeBlack, mOwner, mTspAddress, expiredAt)
				store := ctx.KVStore(k.StoreKey)
				staddr := k.storageForAddr(mOwner)
				store.Set(staddr, k.Cdc.MustMarshalBinaryBare(&membership))
			} else {
				err := k.RemoveMembership(ctx, mOwner)
				if err != nil {
					panic(err)
				}
			}
		}
	}
	return nil
}
