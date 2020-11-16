package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
	"github.com/cosmos/cosmos-sdk/x/supply"
	supplyExported "github.com/cosmos/cosmos-sdk/x/supply/exported"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/vbr/types"
)

const (
	// 0.001
	dailyBlocks = 12960
)

type Keeper struct {
	cdc          *codec.Codec
	storeKey     sdk.StoreKey
	distKeeper   distribution.Keeper
	supplyKeeper supply.Keeper
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, dk distribution.Keeper, sk supply.Keeper) Keeper {
	return Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		distKeeper:   dk,
		supplyKeeper: sk,
	}
}

// -------------
// --- Pool
// -------------

// SetTotalRewardPool allows to set the value of the total rewards pool that has left
func (k Keeper) SetTotalRewardPool(ctx sdk.Context, updatedPool sdk.DecCoins) {
	store := ctx.KVStore(k.storeKey)
	if !updatedPool.Empty() {
		store.Set([]byte(types.PoolStoreKey), k.cdc.MustMarshalBinaryBare(&updatedPool))
	} else {
		store.Delete([]byte(types.PoolStoreKey))
	}
}

// GetTotalRewardPool returns the current total rewards pool amount
func (k Keeper) GetTotalRewardPool(ctx sdk.Context) sdk.DecCoins {
	macc := k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	mcoins := macc.GetCoins()

	return sdk.NewDecCoinsFromCoins(mcoins...)
}

// --------------------
// --- Year number
// --------------------

var (
	// DPY is Days Per Year
	DPY = sdk.NewDecWithPrec(36525, 2)
	// HPD is Hours Per Day
	HPD = sdk.NewDecWithPrec(24, 0)
	// MPH  is Minutes Per Hour
	MPH = sdk.NewDecWithPrec(60, 0)
	// BPM is Blocks Per Minutes
	BPM = sdk.NewDecWithPrec(9, 0)

	// BPD is Blocks Per Day
	BPD = HPD.Mul(MPH).Mul(BPM)

	// BPY is Blocks Per Year
	BPY = DPY.Mul(BPD)

	rewardRate = sdk.NewDecWithPrec(1, 3)
)

// ---------------------------
// --- Reward distribution
// ---------------------------

// ComputeProposerReward computes the final reward for the validator block's proposer
func (k Keeper) ComputeProposerReward(ctx sdk.Context, validatorsCount int64,
	proposer exported.ValidatorI, totalStakedTokens sdk.Int) sdk.DecCoins {

	// Get total bonded token
	proposerBonded := proposer.GetBondedTokens()

	// Get rewarded rate with
	// rewardRate

	rewardRateVal := rewardRate.Mul(sdk.NewDec(validatorsCount)).Quo(sdk.NewDec(100))

	// Compute the voting power for this validator at the current block
	VotingPower := proposerBonded.ToDec().Quo(totalStakedTokens.ToDec())
	exptedDailyBlocks := sdk.NewDec(dailyBlocks).Mul(VotingPower)

	Rnb := sdk.NewDecCoinsFromCoins(sdk.NewCoin("ucommercio", proposerBonded.ToDec().Mul(rewardRateVal).Quo(exptedDailyBlocks).TruncateInt()))

	return Rnb
}

// DistributeBlockRewards distributes the computed reward to the block proposer
func (k Keeper) DistributeBlockRewards(ctx sdk.Context, validator exported.ValidatorI, reward sdk.DecCoins) error {
	rewardPool := k.GetTotalRewardPool(ctx)
	// Check if the yearly pool and the total pool have enough funds
	if ctypes.IsAllGTE(rewardPool, reward) {
		// truncate fractional part and only take the integer part into account
		rewardInt, _ := reward.TruncateDecimal()
		k.SetTotalRewardPool(ctx, rewardPool.Sub(sdk.NewDecCoinsFromCoins(rewardInt...)))

		err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, distribution.ModuleName, rewardInt)
		if err != nil {
			return fmt.Errorf("could not send tokens from vbr to distribution module accounts: %w", err)
		}
		k.distKeeper.AllocateTokensToValidator(ctx, validator, sdk.NewDecCoinsFromCoins(rewardInt...))
	} else {
		return sdkErr.Wrap(sdkErr.ErrInsufficientFunds, "Pool hasn't got enough funds to supply validator's rewards")
	}

	return nil
}

// VbrAccount returns vbr's ModuleAccount
func (k Keeper) VbrAccount(ctx sdk.Context) supplyExported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// MintVBRTokens mints coins into the vbr's ModuleAccount
func (k Keeper) MintVBRTokens(ctx sdk.Context, coins sdk.Coins) error {
	if err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return fmt.Errorf("could not mint requested coins: %w", err)
	}

	return nil
}
