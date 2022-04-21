package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	accTypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

const (
	stakeDenom            = "ucommercio"
	stableCreditDenom     = "uccc"
	eventBuyMembership    = "buy_membership"
	eventAssignMembership = "assign_membership"
	eventRemoveMembership = "remove_membership"
	eventDistributeReward = "distribute_reward"
	eventDeposit          = "deposit_into_pool"
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

	// Check if the expired at is greater then current time
	// Blocktime maybe better with ctx.BlockHeader().Time
	if expited_at.Before(ctx.BlockTime()) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid expiry date: %s is before current block time", expited_at))
	}

	// Delete membership if exists
	_ = k.DeleteMembership(ctx, user)

	// Check if user already has a membership.
	// TODO: this check wont pass if DeleteMembership doesn't work.
	//       Maybe it's better to check error from DeleteMembership method
	store := ctx.KVStore(k.storeKey)
	staddr := k.storageForAddr(user)
	if store.Has(staddr) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest,
			fmt.Sprintf(
				"cannot add membership \"%s\" for address %s: user already has a membership",
				membershipType,
				user,
			),
		)
	}

	// Save membership
	membership := types.NewMembership(membershipType, user, tsp, expited_at.UTC())
	store.Set(staddr, k.cdc.MustMarshal(&membership))

	// TODO: add event to distinguish assign from buy, or add specific event to eventManager in buy method
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventAssignMembership,
		sdk.NewAttribute("owner", membership.Owner),
		sdk.NewAttribute("membership_type", membership.MembershipType),
		sdk.NewAttribute("tsp_address", membership.TspAddress),
		sdk.NewAttribute("expiry_at", membership.ExpiryAt.String()),
	))

	logger := k.Logger(ctx)
	logger.Debug("Membership successfully set up")

	return nil
}

// DeleteMembership allows to remove any existing membership associated with the given user.
func (k Keeper) DeleteMembership(ctx sdk.Context, user sdk.AccAddress) error {
	store := ctx.KVStore(k.storeKey)

	// Check if membership must be deleted is owned user by a trust service provider
	if k.IsTrustedServiceProvider(ctx, user) {
		return sdkErr.Wrap(sdkErr.ErrUnauthorized,
			fmt.Sprintf("account \"%s\" is a Trust Service Provider: remove from tsps list before", user.String()),
		)
	}

	// Check if user has a membership
	if !store.Has(k.storageForAddr(user)) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest,
			fmt.Sprintf("account \"%s\" does not have any membership", user.String()),
		)
	}

	// Delete membership user
	store.Delete(k.storageForAddr(user))
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventRemoveMembership,
		sdk.NewAttribute("subscriber", user.String()),
	))
	logger := k.Logger(ctx)
	logger.Debug("Membership successfully removed")

	return nil
}

// DistributeReward allows to distribute the rewards to the sender of the specified invite upon the receiver has
// properly bought a membership of the given membershipType
// TODO: method returns an error even if the membership has been purchased. Maybe need returns a boolean or evalutes error in different way
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
	var rewardCrossValue sdk.Dec
	var ok bool
	if rewardCrossValue, ok = membershipRewards[senderMembershipType][recipientMembershipType]; !ok {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, "Invalid reward options")
	}
	rewardAmount := rewardCrossValue.MulInt64(1000000).TruncateInt()
	//rewardAmount := membershipRewards[senderMembershipType][recipientMembershipType].MulInt64(1000000).TruncateInt()

	// Get the pool amount
	poolAmount := k.GetPoolFunds(ctx).AmountOf(stakeDenom)

	// Distribute the reward taking it from the pool amount
	// TODO: return immediatly if there is no funds
	var returnMethod error
	returnMethod = nil
	if poolAmount.GT(sdk.ZeroInt()) {

		// If the reward is more than the current pool amount, set the reward as the total pool amount
		if rewardAmount.GT(poolAmount) {
			rewardAmount = poolAmount
		}
		// Calcute equivalent distribution in uccc
		ucccConversionRate := k.mintKeeper.GetConversionRate(ctx)
		//kmintTypes.GetConv

		rewardCoins := sdk.NewCoins(sdk.NewCoin(stableCreditDenom, rewardAmount))
		// TODO check calculation mint amount. See calculation of mint
		// TODO check conversion rate not zero
		rewardStakeCoinAmount := sdk.NewDecFromInt(rewardAmount).Mul(ucccConversionRate).Ceil().TruncateInt()
		stakeEquivCoins := sdk.NewCoins(sdk.NewCoin(stakeDenom, rewardStakeCoinAmount))

		if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, stakeEquivCoins); err != nil {
			return sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
		}

		// Need Mint tokens to module and send to account
		//
		// k.bankKeeper.AddCoins(ctx, inviteSender, rewardCoins)
		// WARNING: some problems here. MintCoins inflate the Supply
		// Before AddCoins didn't that
		if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, rewardCoins); err != nil {
			return sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
		}

		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, inviteSender, rewardCoins); err != nil {
			return sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
		}

		// Emits events
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			eventDistributeReward,
			sdk.NewAttribute("invite_sender", invite.Sender),
			sdk.NewAttribute("reward_coins", rewardCoins.String()),
			sdk.NewAttribute("sender_membership_type", senderMembershipType),
			sdk.NewAttribute("recipient_membership_type", recipientMembership.MembershipType),
			sdk.NewAttribute("invite_recipient", invite.User),
		))

	} else {
		returnMethod = sdkErr.Wrap(sdkErr.ErrUnauthorized, "ABR pool has zero tokens")
	}

	// Set the invitation as rewarded
	newInvite := types.Invite{
		Sender:           invite.Sender,
		User:             invite.User,
		SenderMembership: invite.SenderMembership,
		Status:           uint64(types.InviteStatusRewarded), // TODO control conversion
	}

	k.SaveInvite(ctx, newInvite)

	return returnMethod
}

func IsValidMembership(ctx sdk.Context, expiredAt time.Time, mt string) bool {
	if expiredAt.Before(ctx.BlockTime()) && mt != types.MembershipTypeBlack {
		return false
	}
	return true
}

// GetMembership allows to retrieve any existent membership for the specified user.
func (k Keeper) GetMembership(ctx sdk.Context, user sdk.AccAddress) (types.Membership, error) {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(k.storageForAddr(user)) {
		return types.Membership{}, sdkErr.Wrap(sdkErr.ErrUnknownRequest,
			fmt.Sprintf("membership not found for user \"%s\"", user.String()),
		)
	}

	membershipRaw := store.Get(k.storageForAddr(user))
	var ms types.Membership
	k.cdc.MustUnmarshal(membershipRaw, &ms)
	if !IsValidMembership(ctx, *ms.ExpiryAt, ms.MembershipType) {
		return types.Membership{}, sdkErr.Wrap(sdkErr.ErrUnknownRequest,
			fmt.Sprintf("membership for user \"%s\" has expired", user.String()),
		)
	}
	return ms, nil
}

// GetMemberships extracts all memerships
func (k Keeper) GetMemberships(ctx sdk.Context) []*types.Membership {
	im := k.MembershipIterator(ctx)
	ms := []*types.Membership{}
	defer im.Close()
	for ; im.Valid(); im.Next() {
		var m types.Membership
		k.cdc.MustUnmarshal(im.Value(), &m)
		// Returns only valid memberships
		if !IsValidMembership(ctx, *m.ExpiryAt, m.MembershipType) {
			continue
		}
		ms = append(ms, &m)
	}

	return ms
}

// MembershipIterator returns an Iterator for all the memberships stored.
func (k Keeper) MembershipIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte(types.MembershipsStorageKey))
}

// ComputeExpiryHeight compute expiry height of membership.
func (k Keeper) ComputeExpiryHeight(blockTime time.Time) time.Time {
	expirationAt := blockTime.Add(secondsPerYear)
	return expirationAt
}

// GetTspMemberships extracts all memerships
func (k Keeper) GetTspMemberships(ctx sdk.Context, tsp sdk.Address) types.Memberships {
	im := k.MembershipIterator(ctx)
	m := types.Membership{}
	ms := types.Memberships{}
	defer im.Close()
	for ; im.Valid(); im.Next() {
		k.cdc.MustUnmarshal(im.Value(), &m)
		if m.TspAddress != tsp.String() {
			continue
		}
		ms = append(ms, m)
	}

	return ms
}

// GetMembershipModuleAccount returns the module account for the commerciokyc module
func (k Keeper) GetMembershipModuleAccount(ctx sdk.Context) accTypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// storageForAddr returns a string representing the KVStore storage key for an addr.
func (k Keeper) storageForAddr(addr sdk.AccAddress) []byte {
	//return append([]byte(types.MembershipsStorageKey), k.cdc.MustMarshal(&addr)...)
	return append([]byte(types.MembershipsStorageKey), addr.Bytes()...)
}
