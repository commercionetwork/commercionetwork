package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	distrTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	accountTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govTypes "github.com/commercionetwork/commercionetwork/x/government/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	epochsTypes "github.com/commercionetwork/commercionetwork/x/epochs/types"

	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"

	govT "github.com/commercionetwork/commercionetwork/x/government/types"

	distrKeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	accountKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	govKeeper "github.com/commercionetwork/commercionetwork/x/government/keeper"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	paramsKeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	epochsKeeper "github.com/commercionetwork/commercionetwork/x/epochs/keeper"
)

var (
	distrAcc  = accountTypes.NewEmptyModuleAccount(types.ModuleName)
	TestFunder, _    = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
	TestDelegator, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
	valAddr, _       = sdk.ValAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")
	valAddrVal, _    = sdk.ValAddressFromBech32("cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3")
	valDelAddr, _    = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")
	PKs =  simapp.CreateTestPubKeys(10)
	TestValidator, _        = stakingTypes.NewValidator(valAddrVal, PKs[0], stakingTypes.Description{})
	TestAmount           = sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(100)))
	TestBlockRewardsPool = sdk.NewDecCoinsFromCoins(sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(100000), Denom: "stake"})...)
	TestRewarRate        = sdk.NewDecWithPrec(12, 3)
)

func SetupTestInput(emptyPool bool) (ctx sdk.Context, keeper Keeper) {
	memDB := db.NewMemDB()
	stateStore := store.NewCommitMultiStore(memDB)
	registry := codectypes.NewInterfaceRegistry()

	keys := sdk.NewKVStoreKeys(
		paramsTypes.StoreKey,
		accountTypes.StoreKey,
		bankTypes.StoreKey,
		stakingTypes.StoreKey,
		types.StoreKey,
		distrTypes.StoreKey,
		govT.StoreKey,
	)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	memStoreKeyGov := storetypes.NewMemoryStoreKey(govTypes.MemStoreKey)

	tStakingKey := sdk.NewTransientStoreKey("transient_test")
	tkeys := sdk.NewTransientStoreKeys(paramsTypes.TStoreKey, tStakingKey.String())

	ms := store.NewCommitMultiStore(memDB)
	for _, key := range keys {
		ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, memDB)
	}

	for _, tkey := range tkeys {
		ms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, memDB)
	}

	_ = ms.LoadLatestVersion()
	feeCollectorAcc := accountTypes.NewEmptyModuleAccount(accountTypes.FeeCollectorName)
	notBondedPool := accountTypes.NewEmptyModuleAccount(stakingTypes.NotBondedPoolName, accountTypes.Burner, accountTypes.Staking)
	bondPool := accountTypes.NewEmptyModuleAccount(stakingTypes.BondedPoolName, accountTypes.Burner, accountTypes.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.GetAddress().String()] = true
	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
	blacklistedAddrs[bondPool.GetAddress().String()] = true
	blacklistedAddrs[distrAcc.GetAddress().String()] = true

	ctx = sdk.NewContext(stateStore, tmproto.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	// add module accounts to supply keeper
	maccPerms := map[string][]string{
		accountTypes.FeeCollectorName:     nil,
		distrTypes.ModuleName:          nil,
		stakingTypes.BondedPoolName:    {accountTypes.Burner, accountTypes.Staking},
		stakingTypes.NotBondedPoolName: {accountTypes.Burner, accountTypes.Staking},
		types.ModuleName:          {accountTypes.Minter},
	}

	pk := paramsKeeper.NewKeeper(codec.NewProtoCodec(registry), codec.NewLegacyAmino(), keys[paramsTypes.StoreKey], tkeys[paramsTypes.TStoreKey])
	ak := accountKeeper.NewAccountKeeper(codec.NewProtoCodec(registry), keys[accountTypes.StoreKey], pk.Subspace("auth"), accountTypes.ProtoBaseAccount, maccPerms)
	bk := bankKeeper.NewBaseKeeper(codec.NewProtoCodec(registry), keys[bankTypes.StoreKey], ak, pk.Subspace("bank"), blacklistedAddrs)
	sk := stakingKeeper.NewKeeper(codec.NewProtoCodec(registry), keys[stakingTypes.StoreKey], ak, bk, pk.Subspace("staking"))
	gk := govKeeper.NewKeeper(codec.NewProtoCodec(registry), keys[govTypes.StoreKey], memStoreKeyGov)
	dk := distrKeeper.NewKeeper(codec.NewProtoCodec(registry), keys[distrTypes.StoreKey], pk.Subspace("distribution"),ak, bk, sk, accountTypes.FeeCollectorName, blacklistedAddrs)
	sk.SetHooks(dk.Hooks())
	ek := epochsKeeper.NewKeeper(codec.NewProtoCodec(registry), keys[epochsTypes.StoreKey])
	subspace, _ := pk.GetSubspace(types.ModuleName)
	k := NewKeeper(
		codec.NewProtoCodec(registry),
		keys[types.StoreKey],
		memStoreKey,
		dk,
		bk,
		ak,
		*gk,
		*ek,
		subspace,
		sk,
	)

	if !emptyPool {
		pool, _ := TestBlockRewardsPool.TruncateDecimal()
		macc := k.VbrAccount(ctx)

		k.bankKeeper.SetBalances(ctx, macc.GetAddress(), sdk.NewCoins(pool...))
		k.accountKeeper.SetModuleAccount(ctx, macc)
	}

	return ctx, *k
}