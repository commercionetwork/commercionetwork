package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

var membershipRewards = map[string]map[string]sdk.Dec{
	types.MembershipTypeBronze: {
		types.MembershipTypeBronze: sdk.NewDecWithPrec(125, 2),  // 5% of 1.25
		types.MembershipTypeSilver: sdk.NewDecWithPrec(25, 0),   // 10% of 250
		types.MembershipTypeGold:   sdk.NewDecWithPrec(375, 0),  // 15% of 2500
		types.MembershipTypeBlack:  sdk.NewDecWithPrec(5000, 0), // 20% of 25000
	},
	types.MembershipTypeSilver: {
		types.MembershipTypeBronze: sdk.NewDecWithPrec(5, 0),     // 20% of 1.25
		types.MembershipTypeSilver: sdk.NewDecWithPrec(75, 0),    // 30% of 250
		types.MembershipTypeGold:   sdk.NewDecWithPrec(1000, 0),  // 40% of 2500
		types.MembershipTypeBlack:  sdk.NewDecWithPrec(12500, 0), // 50% of 25000
	},
	types.MembershipTypeGold: {
		types.MembershipTypeBronze: sdk.NewDecWithPrec(125, 1),   // 50% of 1.25
		types.MembershipTypeSilver: sdk.NewDecWithPrec(150, 0),   // 60% of 250
		types.MembershipTypeGold:   sdk.NewDecWithPrec(1750, 0),  // 70% of 2500
		types.MembershipTypeBlack:  sdk.NewDecWithPrec(20000, 0), // 80% of 25000
	},
	types.MembershipTypeBlack: {
		types.MembershipTypeBronze: sdk.NewDecWithPrec(175, 2),   // 70% of 1.25
		types.MembershipTypeSilver: sdk.NewDecWithPrec(200, 0),   // 80% of 250
		types.MembershipTypeGold:   sdk.NewDecWithPrec(2250, 0),  // 90% of 2500
		types.MembershipTypeBlack:  sdk.NewDecWithPrec(25000, 0), // 100% of 25000
	},
}

// DepositIntoPool allows the depositor to deposit the specified amount inside the rewards pool
func (k Keeper) DepositIntoPool(ctx sdk.Context, depositor sdk.AccAddress, amount sdk.Coins) error {
	// Send the coins from the user wallet to the pool
	for _, coin := range amount {
		if coin.Denom != "ucommercio" {
			return sdkErr.Wrap(sdkErr.ErrInsufficientFunds, "deposit into membership pool can only be expressed in ucommercio")
		}
	}

	if err := k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, amount); err != nil {
		return err
	}
	return nil
}

// DistributeReward allows to distribute the rewards to the sender of the specified invite upon the receiver has
// properly bought a membership of the given membershipType
func (k Keeper) DistributeReward(ctx sdk.Context, invite types.Invite) error {
	// the invite we got is either invalid or already rewarded, get out!
	if invite.Status == types.InviteStatusRewarded || invite.Status == types.InviteStatusInvalid {
		return nil
	}
	// Calculate reward for invite
	_, err := k.GetMembership(ctx, invite.Sender)
	if err != nil || invite.SenderMembership == "" {
		return sdkErr.Wrap(sdkErr.ErrUnauthorized, "Invite sender does not have a membership")
	}

	recipientMembership, err := k.GetMembership(ctx, invite.User)
	if err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnauthorized, "Invite recipient does not have a membership")
	}

	senderMembershipType := invite.SenderMembership
	recipientMembershipType := recipientMembership.MembershipType

	// Get the reward amount by searching up inside the matrix.
	// Multiply the found amount by 1.000.000 as coins are represented as millionth of units, and make it an int
	rewardAmount := membershipRewards[senderMembershipType][recipientMembershipType].MulInt64(1000000).TruncateInt()

	// Get the pool amount
	poolAmount := k.GetPoolFunds(ctx).AmountOf("ucommercio")

	// Distribute the reward taking it from the pool amount
	if poolAmount.GT(sdk.ZeroInt()) {

		// If the reward is more than the current pool amount, set the reward as the total pool amount
		if rewardAmount.GT(poolAmount) {
			rewardAmount = poolAmount
		}
		rewardCoins := sdk.NewCoins(sdk.NewCoin("ucommercio", rewardAmount))

		// Send the reward to the invite sender
		if err := k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, invite.Sender, rewardCoins); err != nil {
			return err
		}
	}

	// Set the invitation as rewarded
	newInvite := types.Invite{
		Sender:           invite.Sender,
		User:             invite.User,
		SenderMembership: invite.SenderMembership,
		Status:           types.InviteStatusRewarded,
	}

	k.SaveInvite(ctx, newInvite)

	return nil
}
