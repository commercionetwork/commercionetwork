package keeper
import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	distrTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	accountTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govTypes "github.com/commercionetwork/commercionetwork/x/government/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	epochsTypes "github.com/commercionetwork/commercionetwork/x/epochs/types"

	distrKeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	accountKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	govKeeper "github.com/commercionetwork/commercionetwork/x/government/keeper"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	paramsKeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	epochsKeeper "github.com/commercionetwork/commercionetwork/x/epochs/keeper"
)

var (
	//distrAcc  = accountTypes.NewEmptyModuleAccount(types.ModuleName)
)

func setupKeeper(t testing.TB) (*Keeper, sdk.Context) {
	storeKeys := sdk.NewKVStoreKeys(
		types.StoreKey,
		paramsTypes.StoreKey,
		distrTypes.StoreKey,
		bankTypes.StoreKey,
		accountTypes.StoreKey,
		govTypes.StoreKey,
		epochsTypes.StoreKey,
		stakingTypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramsTypes.TStoreKey)
	
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	memStoreKeyGov := storetypes.NewMemoryStoreKey(govTypes.MemStoreKey)


	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	for _, storeKey := range storeKeys {
		stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	}
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(memStoreKeyGov, sdk.StoreTypeMemory, nil)
	
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	
	feeCollectorAcc := accountTypes.NewEmptyModuleAccount(accountTypes.FeeCollectorName)
	notBondedPool := accountTypes.NewEmptyModuleAccount(stakingTypes.NotBondedPoolName, accountTypes.Burner, accountTypes.Staking)
	bondPool := accountTypes.NewEmptyModuleAccount(stakingTypes.BondedPoolName, accountTypes.Burner, accountTypes.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.GetAddress().String()] = true
	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
	blacklistedAddrs[bondPool.GetAddress().String()] = true
	blacklistedAddrs[distrAcc.GetAddress().String()] = true

	ctx := sdk.NewContext(stateStore, tmproto.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	maccPerms := map[string][]string{
		accountTypes.FeeCollectorName:     nil,
		distrTypes.ModuleName:          nil,
		stakingTypes.BondedPoolName:    {accountTypes.Burner, accountTypes.Staking},
		stakingTypes.NotBondedPoolName: {accountTypes.Burner, accountTypes.Staking},
		types.ModuleName:          {accountTypes.Minter},
	}

	pk := paramsKeeper.NewKeeper(codec.NewProtoCodec(registry), codec.NewLegacyAmino(), storeKeys[paramsTypes.StoreKey], tkeys[paramsTypes.TStoreKey])
	ak := accountKeeper.NewAccountKeeper(codec.NewProtoCodec(registry), storeKeys[accountTypes.StoreKey], pk.Subspace("auth"), accountTypes.ProtoBaseAccount, maccPerms)
	bk := bankKeeper.NewBaseKeeper(codec.NewProtoCodec(registry), storeKeys[bankTypes.StoreKey], ak, pk.Subspace("bank"), blacklistedAddrs)
	sk := stakingKeeper.NewKeeper(codec.NewProtoCodec(registry), storeKeys[stakingTypes.StoreKey], ak, bk, pk.Subspace("staking"))
	//sk.SetParams(ctx, stakingTypes.DefaultParams())
	gk := govKeeper.NewKeeper(codec.NewProtoCodec(registry), storeKeys[govTypes.StoreKey], memStoreKeyGov)
	dk := distrKeeper.NewKeeper(codec.NewProtoCodec(registry), storeKeys[distrTypes.StoreKey], pk.Subspace("distribution"),ak, bk, sk, accountTypes.FeeCollectorName, blacklistedAddrs)
	sk.SetHooks(stakingTypes.NewMultiStakingHooks(dk.Hooks()))
	ek := epochsKeeper.NewKeeper(codec.NewProtoCodec(registry), storeKeys[epochsTypes.StoreKey]) 
	keeper := NewKeeper(
		codec.NewProtoCodec(registry),
		storeKeys[types.StoreKey],
		memStoreKey,
		dk,
		bk,
		ak,
		*gk,
		*ek,
		pk.Subspace("vbr"),
		sk,
	)
	ek.SetHooks(epochsTypes.NewMultiEpochHooks(keeper.Hooks()))
	/*params := types.Params{
					DistrEpochIdentifier: types.EpochDay,
					EarnRate: sdk.NewDecWithPrec(5,1),
				}
	keeper.SetParams(ctx, params)*/

	return keeper, ctx
}

// ---------------------------
// --- Reward distribution
// ---------------------------
var Params_test = types.Params{
					DistrEpochIdentifier: types.EpochDay,
					EarnRate: sdk.NewDecWithPrec(5,1),
				}
func TestKeeper_ComputeProposerReward(t *testing.T) {
	tests := []struct {
		name           string
		bonded         sdk.Int
		vNumber        int64
		expectedReward string
		params		   types.Params
	}{
		{
			"Compute reward with 100 validators",
			sdk.NewInt(100000000),
			100,
			"136986.301369863013698630",
			Params_test,
		},
		{
			"Compute reward with 50 validators",
			sdk.NewInt(100000000),
			50,
			"68493.150684931506849315",
			Params_test,
		},
		{
			"Compute reward with small bonded",
			sdk.NewInt(1),
			100,
			"0.001369863013698630",
			Params_test,
		},
		{
			"Compute reward per minute",
			sdk.NewInt(100000000),
			50,
			"47.564687975646879756",
			types.Params{
				DistrEpochIdentifier: types.EpochMinute,
				EarnRate: sdk.NewDecWithPrec(5,1),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			testVal := TestValidator.UpdateStatus(stakingTypes.Bonded)
			testVal, _ = testVal.AddTokensFromDel(tt.bonded)
			//k.SetRewardRate(ctx, TestRewarRate)
			//params := k.GetParams(ctx)
			params := tt.params
			reward := k.ComputeProposerReward(ctx, tt.vNumber, testVal, "ucommercio", params)

			expectedDecReward, _ := sdk.NewDecFromStr(tt.expectedReward)

			expected := sdk.DecCoins{sdk.NewDecCoinFromDec("ucommercio", expectedDecReward)}

			require.Equal(t, expected, reward)

		})
	}
}

//Do not work
func TestKeeper_DistributeBlockRewards(t *testing.T) {
	tests := []struct {
		name              string
		pool              sdk.DecCoins
		expectedValidator sdk.DecCoins
		expectedRemaining sdk.DecCoins
		expectedError     error
		bonded            sdk.Int
	}{
		{
			name:              "Reward with enough pool",
			pool:              sdk.DecCoins{sdk.NewInt64DecCoin("ucommercio", 10000)},
			expectedRemaining: sdk.DecCoins{sdk.NewInt64DecCoin("ucommercio", 9991)},
			expectedValidator: sdk.DecCoins{sdk.NewInt64DecCoin("ucommercio", 9)},
			bonded:            sdk.NewInt(1000000000),
		},
		{
			name:              "Reward with empty pool",
			pool:              sdk.DecCoins{sdk.NewInt64DecCoin("ucommercio", 0)},
			expectedRemaining: sdk.DecCoins{},
			expectedValidator: sdk.DecCoins(nil),
			bonded:            sdk.NewInt(1000000000),
		},
		{
			name:              "Reward not enough funds into pool",
			pool:              sdk.DecCoins{sdk.NewInt64DecCoin("ucommercio", 1)},
			expectedRemaining: sdk.DecCoins{sdk.NewInt64DecCoin("ucommercio", 1)},
			expectedError:     sdkErr.Wrap(sdkErr.ErrInsufficientFunds, "Pool hasn't got enough funds to supply validator's rewards"),
			expectedValidator: sdk.DecCoins(nil),
			bonded:            sdk.NewInt(1000000000),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			testVal := TestValidator.UpdateStatus(stakingTypes.Bonded)
			testVal, _ = testVal.AddTokensFromDel(tt.bonded)

			//k.SetRewardRate(ctx, TestRewarRate)
			k.SetTotalRewardPool(ctx, tt.pool)

			macc := k.VbrAccount(ctx)
			suppl, _ := tt.pool.TruncateDecimal()
			//_ = macc.SetCoins(sdk.NewCoins(suppl...))
			k.bankKeeper.SetBalances(ctx, macc.GetAddress(), sdk.NewCoins(suppl...))
			k.accountKeeper.SetModuleAccount(ctx, macc)

			validatorRewards := distrTypes.ValidatorCurrentRewards{Rewards: sdk.DecCoins{}}
			k.distKeeper.SetValidatorCurrentRewards(ctx, testVal.GetOperator(), validatorRewards)

			validatorOutstandingRewards := distrTypes.ValidatorOutstandingRewards{}
			k.distKeeper.SetValidatorOutstandingRewards(ctx, testVal.GetOperator(), validatorOutstandingRewards)
			//params := k.GetParams(ctx)
			params := types.Params{
				DistrEpochIdentifier: types.EpochDay,
				EarnRate: sdk.NewDecWithPrec(5,1),
			}
			reward := k.ComputeProposerReward(ctx, 1, testVal, "ucommercio", params)

			err := k.DistributeBlockRewards(ctx, testVal, reward)
			if tt.expectedError != nil {
				require.Equal(t, err.Error(), tt.expectedError.Error())
			}

			valCurReward := k.distKeeper.GetValidatorCurrentRewards(ctx, testVal.GetOperator())
			//rewardedVal := valCurReward.Rewards
			rewardPool := k.GetTotalRewardPool(ctx)

			require.Equal(t, tt.expectedRemaining, rewardPool)
			require.Equal(t, tt.expectedValidator, valCurReward.Rewards)

		})
	}
}

func TestKeeper_VbrAccount(t *testing.T) {
	tests := []struct {
		name              string
		wantModName       string
		wantModAccBalance sdk.Coins
		emptyPool         bool
	}{
		{
			"an empty vbr account",
			"vbr",
			sdk.NewCoins(),
			true,
		},
		{
			"a vbr account with coins in it",
			"vbr",
			sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(100000), Denom: "ucommercio"}),
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)
			macc := k.VbrAccount(ctx)

			require.Equal(t, macc.GetName(), tt.wantModName)
			
			if !tt.emptyPool{
				coins := sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(100000), Denom: "ucommercio"})
				k.bankKeeper.SetBalances(ctx, macc.GetAddress(), coins)
			}
			
			require.True(t, k.bankKeeper.GetAllBalances(ctx, macc.GetAddress()).IsEqual(tt.wantModAccBalance))
		})
	}
}

//not working
func TestKeeper_MintVBRTokens(t *testing.T) {
	tests := []struct {
		name       string
		wantAmount sdk.Coins
	}{
		{
			"add 10ucommercio",
			sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(10), Denom: "ucommercio"}),
		},
		{
			"add no ucommercio",
			sdk.NewCoins(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)
			k.bankKeeper.SetSupply(ctx, bankTypes.NewSupply(sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(10), Denom: "ucommercio"})))
			k.MintVBRTokens(ctx, tt.wantAmount)
			macc := k.VbrAccount(ctx)
			//require.True(t, macc.GetCoins().IsEqual(tt.wantAmount))
			require.True(t, k.bankKeeper.GetAllBalances(ctx, macc.GetAddress()).IsEqual(tt.wantAmount))
		})
	}
}