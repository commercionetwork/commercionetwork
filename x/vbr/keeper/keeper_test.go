package keeper
import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
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
		accountTypes.FeeCollectorName:     nil,
		distrTypes.ModuleName:          nil,
		stakingTypes.BondedPoolName:    {accountTypes.Burner, accountTypes.Staking},
		stakingTypes.NotBondedPoolName: {accountTypes.Burner, accountTypes.Staking},
		types.ModuleName:          {accountTypes.Minter},
		govTypes.ModuleName:              {accountTypes.Burner},
	}

	pk := paramsKeeper.NewKeeper(cdc, codec.NewLegacyAmino(), storeKeys[paramsTypes.StoreKey], tkeys[paramsTypes.TStoreKey])
	ak := accountKeeper.NewAccountKeeper(cdc, storeKeys[accountTypes.StoreKey], pk.Subspace("auth"), accountTypes.ProtoBaseAccount, maccPerms)
	bk := bankKeeper.NewBaseKeeper(cdc, storeKeys[bankTypes.StoreKey], ak, pk.Subspace("bank"), blacklistedAddrs)
	sk := stakingKeeper.NewKeeper(cdc, storeKeys[stakingTypes.StoreKey], ak, bk, pk.Subspace("staking"))
	sk.SetParams(ctx, stakingTypes.DefaultParams())
	gk := govKeeper.NewKeeper(cdc, storeKeys[govTypes.StoreKey], memStoreKeyGov)
	dk := distrKeeper.NewKeeper(cdc, storeKeys[distrTypes.StoreKey], pk.Subspace("distribution"),ak, bk, sk, accountTypes.FeeCollectorName, blacklistedAddrs)
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
		pk.Subspace("vbr"),
		sk,
	)
	ek.SetHooks(epochsTypes.NewMultiEpochHooks(keeper.Hooks()))
	
	return keeper, ctx
}

var params_test = types.Params{
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
			params := tt.params
			reward := k.ComputeProposerReward(ctx, tt.vNumber, testVal, "ucommercio", params)

			expectedDecReward, _ := sdk.NewDecFromStr(tt.expectedReward)

			expected := sdk.DecCoins{sdk.NewDecCoinFromDec("ucommercio", expectedDecReward)}

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
			pool:              sdk.DecCoins{sdk.NewInt64DecCoin("ucommercio", 100000)},
			expectedRemaining: sdk.DecCoins{sdk.NewInt64DecCoin("ucommercio", 86302)},
			expectedValidator: sdk.DecCoins{sdk.NewInt64DecCoin("ucommercio", 13698)},
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
				EarnRate: sdk.NewDecWithPrec(5,1),
			}
			reward := k.ComputeProposerReward(ctx, 1, testVal, "ucommercio", params)
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