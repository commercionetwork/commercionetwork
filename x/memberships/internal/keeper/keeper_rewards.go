package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
func (keeper Keeper) DepositIntoPool(ctx sdk.Context, depositor sdk.AccAddress, amount sdk.Coins) sdk.Error {
	if !amount.IsValid() || amount.IsAnyNegative() {
		return sdk.ErrInvalidCoins(amount.String())
	}

	store := ctx.KVStore(keeper.StoreKey)

	// Remove the coins from the user wallet
	if _, err := keeper.BankKeeper.SubtractCoins(ctx, depositor, amount); err != nil {
		return err
	}

	// Add the amount to the pool
	var pool sdk.Coins
	keeper.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.LiquidityPoolStoreKey)), &pool)
	pool = pool.Add(amount)
	store.Set([]byte(types.LiquidityPoolStoreKey), keeper.cdc.MustMarshalBinaryBare(&pool))

	return nil
}

// SetPoolFunds allows to set the current pool funds amount
func (keeper Keeper) SetPoolFunds(ctx sdk.Context, pool sdk.Coins) {
	store := ctx.KVStore(keeper.StoreKey)

	if pool == nil {
		store.Delete([]byte(types.LiquidityPoolStoreKey))
	} else {
		store.Set([]byte(types.LiquidityPoolStoreKey), keeper.cdc.MustMarshalBinaryBare(&pool))
	}
}

// GetPoolFunds return the current pool funds for the given context
func (keeper Keeper) GetPoolFunds(ctx sdk.Context) sdk.Coins {
	store := ctx.KVStore(keeper.StoreKey)
	var pool sdk.Coins
	keeper.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.LiquidityPoolStoreKey)), &pool)
	return pool
}

// DistributeReward allows to distribute the rewards to the sender of the specified invite upon the receiver has
// properly bought a membership of the given membershipType
func (keeper Keeper) DistributeReward(ctx sdk.Context, invite types.Invite, membershipType string) sdk.Error {
	senderMembership, found := keeper.GetMembership(ctx, invite.Sender)
	if !found {
		return sdk.ErrUnauthorized("Invite sender does not have a membership")
	}

	recipientMembership, found := keeper.GetMembership(ctx, invite.User)
	if !found {
		return sdk.ErrUnauthorized("Invite recipient does not have a membership")
	}

	senderMembershipType := keeper.GetMembershipType(senderMembership)
	recipientMembershipType := keeper.GetMembershipType(recipientMembership)
	if recipientMembershipType != membershipType {
		return sdk.ErrUnknownRequest("Invite recipient membership is not the same as the bought membership")
	}

	// Get the reward amount by searching up inside the matrix.
	// Multiply the found amount by 1.000.000 as coins are represented as millionth of units, and make it an int
	rewardAmount := membershipRewards[senderMembershipType][recipientMembershipType].MulInt64(1000000).TruncateInt()

	// Create the coins that represent the reward
	stableCreditsDenom := keeper.GetStableCreditsDenom(ctx)

	// Get the pool amount
	poolFunds := keeper.GetPoolFunds(ctx)
	poolAmount := poolFunds.AmountOf(stableCreditsDenom)

	// Distribute the reward taking it from the pool amount
	if poolAmount.GT(sdk.ZeroInt()) {

		// If the reward is more than the current pool amount, set the reward as the total pool amount
		if rewardAmount.GT(poolAmount) {
			rewardAmount = poolAmount
		}
		rewardCoins := sdk.NewCoins(sdk.NewCoin(stableCreditsDenom, rewardAmount))

		// Subtract the amount from the pool
		poolFunds = poolFunds.Sub(rewardCoins)
		keeper.SetPoolFunds(ctx, poolFunds)

		// Send the reward to the invite sender
		if _, err := keeper.BankKeeper.AddCoins(ctx, invite.Sender, rewardCoins); err != nil {
			return err
		}
	}

	// Set the invitation as rewarded
	newInvite := types.Invite{Sender: invite.Sender, User: invite.User, Rewarded: true}
	keeper.SaveInvite(ctx, newInvite)

	return nil
}
