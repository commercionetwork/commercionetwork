package keeper

import (
	"testing"

	"cosmossdk.io/log"
	//"cosmossdk.io/simapp"
	"cosmossdk.io/store"
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	storemetrics "cosmossdk.io/store/metrics"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	epochsKeeper "github.com/commercionetwork/commercionetwork/x/epochs/keeper"
	epochsTypes "github.com/commercionetwork/commercionetwork/x/epochs/types"
	govKeeper "github.com/commercionetwork/commercionetwork/x/government/keeper"
	govTypes "github.com/commercionetwork/commercionetwork/x/government/types"
	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	//"github.com/cosmos/ibc-go/testing/simapp"
	"github.com/stretchr/testify/require"
)

var (
	valAddrVal, _    = sdk.ValAddressFromBech32("cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3")
	PKs              = simapp.CreateTestPubKeys(10)
	TestValidator, _ = stakingTypes.NewValidator(valAddrVal, PKs[0], stakingTypes.Description{})
)

func SetupKeeper(t testing.TB) (*Keeper, sdk.Context) {
	distrAcc := accountTypes.NewEmptyModuleAccount(types.ModuleName)

	storeKeys := storetypes.NewKVStoreKeys(
		types.StoreKey,
		paramsTypes.StoreKey,
		distrTypes.StoreKey,
		bankTypes.StoreKey,
		accountTypes.StoreKey,
		govTypes.StoreKey,
		epochsTypes.StoreKey,
		stakingTypes.StoreKey,
	)
	tkeys := storetypes.NewTransientStoreKeys(paramsTypes.TStoreKey)

	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	memStoreKeyGov := storetypes.NewMemoryStoreKey(govTypes.MemStoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), storemetrics.NewNoOpMetrics())
	for _, storeKey := range storeKeys {
		stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	}
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(memStoreKeyGov, storetypes.StoreTypeMemory, nil)

	stateStore.MountStoreWithDB(tkeys[paramsTypes.TStoreKey], storetypes.StoreTypeTransient, db)
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
	//bk.SetSupply(ctx, bankTypes.NewSupply(sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(100000), Denom: "stake"})))
	//bk.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(100000), Denom: "stake"}))

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
		*sk,
	)
	ek.SetHooks(epochsTypes.NewMultiEpochHooks(keeper.Hooks()))

	government, err := sdk.AccAddressFromBech32(types.ValidMsgSetParams.Government)
	if err != nil {
		panic(err)
	}

	keeper.govKeeper.SetGovernmentAddress(ctx, government)

	return keeper, ctx
}
