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

// --------------------
// --- Blocks per day
// --- Blocks per year
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
	//50% 
	vbr_earn_rate = sdk.NewDecWithPrec(50, 2);
)

type (
	Keeper struct {
		cdc      	codec.Marshaler
		storeKey	sdk.StoreKey
		memKey   sdk.StoreKey
		distKeeper  distKeeper.Keeper
		bankKeeper	bankKeeper.Keeper
		accountKeeper accountKeeper.AccountKeeper
		govKeeper   govKeeper.Keeper
		epochsKeeper epochsKeeper.Keeper
		paramSpace       paramtypes.Subspace
		stakingKeeper stakingKeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.Marshaler,
	storeKey sdk.StoreKey,
	memKey sdk.StoreKey,
	distKeeper   distKeeper.Keeper,
	bankKeeper	bankKeeper.Keeper,
	accountKeeper accountKeeper.AccountKeeper,
	govKeeper    govKeeper.Keeper,
	epochsKeeper epochsKeeper.Keeper,
	paramSpace paramtypes.Subspace,
	stakingKeeper stakingKeeper.Keeper,

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

// ---------------------------
// --- Reward distribution
// ---------------------------
// SetRewardRate store the vbr reward rate.
func (k Keeper) SetRewardRateKeeper(ctx sdk.Context, rate sdk.Dec) error {
	if err := types.ValidateRewardRate(rate); err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	rewardRate := types.VbrRewardrate{RewardRate: rate}
	store.Set([]byte(types.RewardRateKey), k.cdc.MustMarshalBinaryBare(&rewardRate))
	return nil
}

// GetRewardRate retrieve the vbr reward rate.
func (k Keeper) GetRewardRateKeeper(ctx sdk.Context) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	var rate types.VbrRewardrate
	k.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.RewardRateKey)), &rate)
	return rate.RewardRate
}

// SetAutomaticWithdraw store the automatic withdraw flag.
func (k Keeper) SetAutomaticWithdrawKeeper(ctx sdk.Context, autoW bool) error {
	store := ctx.KVStore(k.storeKey)
	autoWithdraw := types.VbrAutoW{AutoW: autoW}
	store.Set([]byte(types.AutomaticWithdraw), k.cdc.MustMarshalBinaryBare(&autoWithdraw))
	return nil
}

// GetAutomaticWithdraw retrieve automatic withdraw flag.
func (k Keeper) GetAutomaticWithdrawKeeper(ctx sdk.Context) bool {
	store := ctx.KVStore(k.storeKey)
	var autoW types.VbrAutoW
	k.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.AutomaticWithdraw)), &autoW)
	return autoW.AutoW
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
func (k Keeper) ComputeProposerReward(ctx sdk.Context, vCount int64, validator stakingTypes.ValidatorI, denom string, epochIdentifier string) sdk.DecCoins {

	// Get rewarded rate
	//rewardRate := k.GetRewardRateKeeper(ctx)

	// Calculate rewarded rate with validator percentage
	//rewardRateVal := rewardRate.Mul(sdk.NewDec(vCount)).Quo(sdk.NewDec(100))

	// Get total bonded token of validator
	validatorBonded := validator.GetBondedTokens()

	validatorBondedPerc := sdk.NewDecCoinFromDec(denom, validatorBonded.ToDec().Mul(vbr_earn_rate))
	validatorsPerc := sdk.NewDec(vCount).QuoInt64(int64(100)) 
	
	//compute the annual distribution ((validator's token * 0.5)*(total_validators/100))
	annualDistribution := sdk.NewDecCoinFromDec(denom, validatorBondedPerc.Amount.Mul(validatorsPerc))
	var epochDuration sdk.Dec
	switch (epochIdentifier){
		case "day": 
			epochDuration = sdk.NewDec(365)
		case "week":
			epochDuration = sdk.NewDec(365).Quo(sdk.NewDec(7))
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
		// TODO this error continue when pool hasn't enough funds for all rewards. Find a method to avoid this
		return sdkErr.Wrap(sdkErr.ErrInsufficientFunds, "Pool hasn't got enough funds to supply validator's rewards")
	}

	return nil
}
