package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	distKeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
 	accountTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	accountKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	govKeeper "github.com/commercionetwork/commercionetwork/x/government/keeper"
	// this line is used by starport scaffolding # ibc/keeper/import
	epochsKeeper "github.com/commercionetwork/commercionetwork/x/epochs/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

type (
	Keeper struct {
		cdc      			codec.Marshaler
		storeKey			sdk.StoreKey
		memKey   			sdk.StoreKey
		distKeeper  		distKeeper.Keeper
		bankKeeper			bankKeeper.Keeper
		accountKeeper 		accountKeeper.AccountKeeper
		govKeeper   		govKeeper.Keeper
		epochsKeeper 		epochsKeeper.Keeper
		paramSpace       	paramtypes.Subspace
		stakingKeeper		stakingKeeper.Keeper
	}
)

func NewKeeper(
	cdc 				codec.Marshaler,
	storeKey 			sdk.StoreKey,
	memKey 				sdk.StoreKey,
	distKeeper   		distKeeper.Keeper,
	bankKeeper			bankKeeper.Keeper,
	accountKeeper 		accountKeeper.AccountKeeper,
	govKeeper    		govKeeper.Keeper,
	epochsKeeper	 	epochsKeeper.Keeper,
	paramSpace 			paramtypes.Subspace,
	stakingKeeper 		stakingKeeper.Keeper,

) *Keeper {

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
		distKeeper: distKeeper,
		bankKeeper: bankKeeper,
		accountKeeper: accountKeeper,
		govKeeper: govKeeper,
		epochsKeeper: epochsKeeper,
		paramSpace: paramSpace,
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
	if !updatedPool.Empty() {
		store.Set([]byte(types.PoolStoreKey), k.cdc.MustMarshalBinaryBare(&pool))
	} else {
		store.Delete([]byte(types.PoolStoreKey))
	}
}
// GetTotalRewardPool returns the current total rewards pool amount
func (k Keeper) GetTotalRewardPool(ctx sdk.Context) sdk.DecCoins {
	macc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	//mcoins := macc.GetCoins()
	coins := GetCoins(k, ctx, macc)

	return sdk.NewDecCoinsFromCoins(coins...)
}

// VbrAccount returns vbr's ModuleAccount
func (k Keeper) VbrAccount(ctx sdk.Context) accountTypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
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
	/*for _, coin := range k.bankKeeper.GetAllBalances(ctx, macc.GetAddress()) {
		coins = append(coins, coin)
	}*/
	coins = append(coins, k.bankKeeper.GetAllBalances(ctx, macc.GetAddress())...)
	
	return coins
}
// ComputeProposerReward computes the final reward for the validator block's proposer
func (k Keeper) ComputeProposerReward(ctx sdk.Context, vCount int64, validator stakingTypes.ValidatorI, denom string, params types.Params) sdk.DecCoins {
	// Get total bonded token of validator
	validatorBonded := validator.GetBondedTokens()

	validatorBondedPerc := sdk.NewDecCoinFromDec(denom, validatorBonded.ToDec().Mul(params.EarnRate))
	validatorsPerc := sdk.NewDec(vCount).QuoInt64(int64(100)) 
	
	//compute the annual distribution ((validator's token * 0.5)*(total_validators/100))
	annualDistribution := sdk.NewDecCoinFromDec(denom, validatorBondedPerc.Amount.Mul(validatorsPerc))
	var epochDuration sdk.Dec
	switch (params.DistrEpochIdentifier){
		case types.EpochDay: 
			epochDuration = sdk.NewDec(365)
		case types.EpochWeek:
			epochDuration = sdk.NewDec(365).Quo(sdk.NewDec(7))
		case types.EpochMinute:
			epochDuration = sdk.NewDec(365*24*60)
		case types.EpochMonthly:
			epochDuration = sdk.NewDec(12)
		default:
			return nil
	}
	
	// Compute reward
	return sdk.NewDecCoins(sdk.NewDecCoinFromDec(denom, annualDistribution.Amount.Quo(epochDuration)))
}

// DistributeBlockRewards distributes the computed reward to the block proposer
func (k Keeper) DistributeBlockRewards(ctx sdk.Context, validator stakingTypes.ValidatorI, reward sdk.DecCoins) error {
	rewardPool := k.GetTotalRewardPool(ctx)
	// Check if the yearly pool and the total pool have enough funds
	if ctypes.IsAllGTE(rewardPool, reward) {
		// truncate fractional part and only take the integer part into account
		rewardInt, _ := reward.TruncateDecimal()

		k.SetTotalRewardPool(ctx, rewardPool.Sub(sdk.NewDecCoinsFromCoins(rewardInt...)))

		err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, distributionTypes.ModuleName, rewardInt)
		if err != nil {
			return nil
		}
		k.distKeeper.AllocateTokensToValidator(ctx, validator, sdk.NewDecCoinsFromCoins(rewardInt...))
	} else {
		return sdkErr.Wrap(sdkErr.ErrInsufficientFunds, "Pool hasn't got enough funds to supply validator's rewards")
	}

	return nil
}
