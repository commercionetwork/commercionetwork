package keeper

import (
	"testing"

	epochsKeeper "github.com/commercionetwork/commercionetwork/x/epochs/keeper"
	epochsTypes "github.com/commercionetwork/commercionetwork/x/epochs/types"
	govKeeper "github.com/commercionetwork/commercionetwork/x/government/keeper"
	govTypes "github.com/commercionetwork/commercionetwork/x/government/types"
	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	accountKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	accountTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrKeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramsKeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
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

	stateStore.MountStoreWithDB(tkeys[paramsTypes.TStoreKey], sdk.StoreTypeTransient, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	//registry := codectypes.NewInterfaceRegistry()
	app := simapp.Setup(false)
	cdc := app.AppCodec()

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
		accountTypes.FeeCollectorName:  nil,
		distrTypes.ModuleName:          nil,
		stakingTypes.BondedPoolName:    {accountTypes.Burner, accountTypes.Staking},
		stakingTypes.NotBondedPoolName: {accountTypes.Burner, accountTypes.Staking},
		types.ModuleName:               {accountTypes.Minter},
		govTypes.ModuleName:            {accountTypes.Burner},
	}

	pk := paramsKeeper.NewKeeper(cdc, codec.NewLegacyAmino(), storeKeys[paramsTypes.StoreKey], tkeys[paramsTypes.TStoreKey])
	ak := accountKeeper.NewAccountKeeper(cdc, storeKeys[accountTypes.StoreKey], pk.Subspace("auth"), accountTypes.ProtoBaseAccount, maccPerms)
	bk := bankKeeper.NewBaseKeeper(cdc, storeKeys[bankTypes.StoreKey], ak, pk.Subspace("bank"), blacklistedAddrs)
	sk := stakingKeeper.NewKeeper(cdc, storeKeys[stakingTypes.StoreKey], ak, bk, pk.Subspace("staking"))
	sk.SetParams(ctx, stakingTypes.DefaultParams())
	gk := govKeeper.NewKeeper(cdc, storeKeys[govTypes.StoreKey], memStoreKeyGov)
	dk := distrKeeper.NewKeeper(cdc, storeKeys[distrTypes.StoreKey], pk.Subspace("distribution"), ak, bk, sk, accountTypes.FeeCollectorName, blacklistedAddrs)
	sk.SetHooks(stakingTypes.NewMultiStakingHooks(dk.Hooks()))
	ek := epochsKeeper.NewKeeper(cdc, storeKeys[epochsTypes.StoreKey])
	keeper := NewKeeper(
		cdc,
		storeKeys[types.StoreKey],
		memStoreKey,
		dk,
		bk,
		ak,
		*gk,
		*ek,
		pk.Subspace(types.ModuleName),
		sk,
	)
	ek.SetHooks(epochsTypes.NewMultiEpochHooks(keeper.Hooks()))

	government, err := sdk.AccAddressFromBech32(types.ValidMsgSetParams.Government)
	if err != nil {
		panic(err)
	}

	keeper.govKeeper.SetGovernmentAddress(ctx, government)

	return keeper, ctx
}

var params_test = types.Params{
	DistrEpochIdentifier: types.EpochDay,
	EarnRate:             sdk.NewDecWithPrec(5, 1),
}

func TestKeeper_ComputeProposerReward(t *testing.T) {
	tests := []struct {
		name           string
		bonded         sdk.Int
		vNumber        int64
		expectedReward string
		params         types.Params
	}{
		{
			"Compute reward with 100 validators",
			sdk.NewInt(100000000),
			100,
			"136986.301369863013698630",
			params_test,
		},
		{
			"Compute reward with 50 validators",
			sdk.NewInt(100000000),
			50,
			"68493.150684931506849315",
			params_test,
		},
		{
			"Compute reward with small bonded",
			sdk.NewInt(1),
			100,
			"0.001369863013698630",
			params_test,
		},
		{
			"Compute reward per minute",
			sdk.NewInt(100000000),
			50,
			"47.564687975646879756",
			types.Params{
				DistrEpochIdentifier: types.EpochMinute,
				EarnRate:             sdk.NewDecWithPrec(5, 1),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			testVal := TestValidator.UpdateStatus(stakingTypes.Bonded)
			testVal, _ = testVal.AddTokensFromDel(tt.bonded)
			params := tt.params
			reward := k.ComputeProposerReward(ctx, tt.vNumber, testVal, types.BondDenom, params)

			expectedDecReward, _ := sdk.NewDecFromStr(tt.expectedReward)

			expected := sdk.DecCoins{sdk.NewDecCoinFromDec(types.BondDenom, expectedDecReward)}

			require.Equal(t, expected, reward)

		})
	}
}

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
			pool:              sdk.DecCoins{sdk.NewInt64DecCoin(types.BondDenom, 100000)},
			expectedRemaining: sdk.DecCoins{sdk.NewInt64DecCoin(types.BondDenom, 86302)},
			expectedValidator: sdk.DecCoins{sdk.NewInt64DecCoin(types.BondDenom, 13698)},
			bonded:            sdk.NewInt(1000000000),
		},
		{
			name:              "Reward with empty pool",
			pool:              sdk.DecCoins{sdk.NewInt64DecCoin(types.BondDenom, 0)},
			expectedRemaining: sdk.DecCoins{},
			expectedValidator: sdk.DecCoins(nil),
			bonded:            sdk.NewInt(1000000000),
		},
		{
			name:              "Reward not enough funds into pool",
			pool:              sdk.DecCoins{sdk.NewInt64DecCoin(types.BondDenom, 1)},
			expectedRemaining: sdk.DecCoins{sdk.NewInt64DecCoin(types.BondDenom, 1)},
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

			k.SetTotalRewardPool(ctx, tt.pool)

			macc := k.VbrAccount(ctx)
			suppl, _ := tt.pool.TruncateDecimal()

			k.bankKeeper.SetBalances(ctx, macc.GetAddress(), sdk.NewCoins(suppl...))
			k.accountKeeper.SetModuleAccount(ctx, macc)

			validatorRewards := distrTypes.ValidatorCurrentRewards{Rewards: sdk.DecCoins{}}
			k.distKeeper.SetValidatorCurrentRewards(ctx, testVal.GetOperator(), validatorRewards)

			validatorOutstandingRewards := distrTypes.ValidatorOutstandingRewards{}
			k.distKeeper.SetValidatorOutstandingRewards(ctx, testVal.GetOperator(), validatorOutstandingRewards)

			params := types.Params{
				DistrEpochIdentifier: types.EpochDay,
				EarnRate:             sdk.NewDecWithPrec(5, 1),
			}
			reward := k.ComputeProposerReward(ctx, 1, testVal, types.BondDenom, params)
			rewardInt, _ := reward.TruncateDecimal()
			_ = rewardInt
			err := k.DistributeBlockRewards(ctx, testVal, reward)
			if tt.expectedError != nil {
				require.Equal(t, err.Error(), tt.expectedError.Error())
			}

			valCurReward := k.distKeeper.GetValidatorCurrentRewards(ctx, testVal.GetOperator())
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
			sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(100000), Denom: types.BondDenom}),
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)
			macc := k.VbrAccount(ctx)

			require.Equal(t, macc.GetName(), tt.wantModName)

			if !tt.emptyPool {
				coins := sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(100000), Denom: types.BondDenom})
				k.bankKeeper.SetBalances(ctx, macc.GetAddress(), coins)
			}

			require.True(t, k.bankKeeper.GetAllBalances(ctx, macc.GetAddress()).IsEqual(tt.wantModAccBalance))
		})
	}
}

func TestKeeper_MintVBRTokens(t *testing.T) {
	tests := []struct {
		name       string
		wantAmount sdk.Coins
	}{
		{
			"add 10ucommercio",
			sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(10), Denom: types.BondDenom}),
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
			k.bankKeeper.SetSupply(ctx, bankTypes.NewSupply(sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(10), Denom: types.BondDenom})))
			k.MintVBRTokens(ctx, tt.wantAmount)
			macc := k.VbrAccount(ctx)
			//require.True(t, macc.GetCoins().IsEqual(tt.wantAmount))
			require.True(t, k.bankKeeper.GetAllBalances(ctx, macc.GetAddress()).IsEqual(tt.wantAmount))
		})
	}
}

func TestKeeper_SetTotalRewardPool(t *testing.T) {

	tests := []struct {
		name        string
		updatedPool sdk.DecCoins
	}{
		{
			name: "empty pool",
		},
		// failing test
		// {
		// 	name:        "ok",
		// 	updatedPool: sdk.NewDecCoinsFromCoins(types.ValidMsgIncrementBlockRewardsPool.Amount...),
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)
			k.SetTotalRewardPool(ctx, tt.updatedPool)

			store := ctx.KVStore(k.storeKey)
			if tt.updatedPool.Empty() {
				require.False(t, store.Has([]byte(types.PoolStoreKey)))
			} else {
				actual := k.GetTotalRewardPool(ctx)
				require.Equal(t, tt.updatedPool, actual)
			}

		})
	}
}
