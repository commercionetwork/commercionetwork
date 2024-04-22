package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	errorsmod "cosmossdk.io/errors"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	storetypes "cosmossdk.io/store/types"

	govKeeper "github.com/commercionetwork/commercionetwork/x/government/keeper"
	accountKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	accountTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distKeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	epochsKeeper "github.com/commercionetwork/commercionetwork/x/epochs/keeper"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type (
	Keeper struct {
		cdc           codec.Codec
		storeKey      storetypes.StoreKey
		memKey        storetypes.StoreKey
		distKeeper    distKeeper.Keeper
		bankKeeper    bankKeeper.Keeper
		accountKeeper accountKeeper.AccountKeeper
		govKeeper     govKeeper.Keeper
		epochsKeeper  epochsKeeper.Keeper
		paramSpace    paramtypes.Subspace
		stakingKeeper stakingKeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	memKey storetypes.StoreKey,
	distKeeper distKeeper.Keeper,
	bankKeeper bankKeeper.Keeper,
	accountKeeper accountKeeper.AccountKeeper,
	govKeeper govKeeper.Keeper,
	epochsKeeper epochsKeeper.Keeper,
	paramSpace paramtypes.Subspace,
	stakingKeeper stakingKeeper.Keeper,

) *Keeper {
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		distKeeper:    distKeeper,
		bankKeeper:    bankKeeper,
		accountKeeper: accountKeeper,
		govKeeper:     govKeeper,
		epochsKeeper:  epochsKeeper,
		paramSpace:    paramSpace,
		stakingKeeper: stakingKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// -------------
// --- Pool
// -------------

// SetTotalRewardPool allows to set the value of the total rewards pool that has left
func (k Keeper) SetTotalRewardPool(ctx sdk.Context, updatedPool sdk.DecCoins) {
	store := ctx.KVStore(k.storeKey)
	pool := types.VbrPool{Amount: updatedPool}
	// TODO: can use standard KeyPrefix function
	// types.KeyPrefix(types.PoolStoreKey)
	poolKeyPrefix := []byte(types.PoolStoreKey)
	if !updatedPool.Empty() {
		store.Set(poolKeyPrefix, k.cdc.MustMarshal(&pool))
	} else {
		store.Delete(poolKeyPrefix)
	}
}

// VbrAccount returns vbr's ModuleAccount
func (k Keeper) VbrAccount(ctx sdk.Context) accountTypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// GetTotalRewardPool returns the current total rewards pool amount
func (k Keeper) GetTotalRewardPool(ctx sdk.Context) sdk.DecCoins {
	macc := k.VbrAccount(ctx)
	coins := GetCoins(k, ctx, macc)

	return sdk.NewDecCoinsFromCoins(coins...)
}

// MintVBRTokens mints coins into the vbr's ModuleAccount
func (k Keeper) MintVBRTokens(ctx sdk.Context, coins sdk.Coins) error {
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return fmt.Errorf("could not mint requested coins: %w", err)
	}

	return nil
}

func GetCoins(k Keeper, ctx sdk.Context, macc accountTypes.ModuleAccountI) sdk.Coins {
	var coins sdk.Coins
	coins = append(coins, k.bankKeeper.GetAllBalances(ctx, macc.GetAddress())...)
	return coins
}

// ComputeProposerReward computes the final reward for the validator block's proposer
func (k Keeper) ComputeProposerReward(ctx sdk.Context, vCount int64, validator stakingTypes.ValidatorI, denom string, params types.Params) sdk.DecCoins {
	// Get total bonded token of validator

	validatorBonded := validator.GetBondedTokens()

	validatorBondedPerc := sdk.NewDecCoinFromDec(denom, validatorBonded.ToDec().Mul(params.EarnRate))
	// TODO: number of validator should be get from staking module
	// paramsVal := k.stakingKeeper.GetParams(ctx)
	// paramsVal.MaxValidators
	validatorsPerc := sdk.NewDec(vCount).QuoInt64(int64(100))

	//compute the annual distribution ((validator's token * 0.5)*(total_validators/100))
	annualDistribution := sdk.NewDecCoinFromDec(denom, validatorBondedPerc.Amount.Mul(validatorsPerc))
	var epochDuration sdk.Dec
	switch params.DistrEpochIdentifier {
	case types.EpochDay:
		epochDuration = sdk.NewDec(365)
	case types.EpochWeek:
		epochDuration = sdk.NewDec(365).Quo(sdk.NewDec(7))
	case types.EpochMinute:
		epochDuration = sdk.NewDec(365 * 24 * 60)
	case types.EpochHour:
		epochDuration = sdk.NewDec(365 * 24)
	case types.EpochMonth:
		epochDuration = sdk.NewDec(12)
	default:
		return nil
	}

	// Compute reward
	return sdk.NewDecCoins(sdk.NewDecCoinFromDec(denom, annualDistribution.Amount.Quo(epochDuration)))
}

// DistributeBlockRewards distributes the computed reward to the block proposer
func (k Keeper) DistributeBlockRewards(ctx sdk.Context, validator stakingTypes.ValidatorI, reward sdk.DecCoins) error {
	// TODO: mybe it's better to get from out of DistributeBlockRewards method the rewardPool amount
	// so you can don't call the method
	rewardPool := k.GetTotalRewardPool(ctx)
	// Check if the yearly pool and the total pool have enough funds
	if ctypes.IsAllGTE(rewardPool, reward) {
		// truncate fractional part and only take the integer part into account
		rewardInt, _ := reward.TruncateDecimal()

		k.SetTotalRewardPool(ctx, rewardPool.Sub(sdk.NewDecCoinsFromCoins(rewardInt...)))

		err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, distributionTypes.ModuleName, rewardInt)
		if err != nil {
			// TODO: if sent token module to module fails the pool should be return to previous value
			// k.SetTotalRewardPool(ctx, rewardPool)
			return nil
		}
		k.distKeeper.AllocateTokensToValidator(ctx, validator, sdk.NewDecCoinsFromCoins(rewardInt...))
	} else {
		return errorsmod.Wrap(sdkErr.ErrInsufficientFunds, "Pool hasn't got enough funds to supply validator's rewards")
	}

	return nil
}
