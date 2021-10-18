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

	distrTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	accountTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govTypes "github.com/commercionetwork/commercionetwork/x/government/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"

	distrKeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	accountKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	govKeeper "github.com/commercionetwork/commercionetwork/x/government/keeper"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	paramsKeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
)

var (
	distrAcc  = accountTypes.NewEmptyModuleAccount(types.ModuleName)
)

func setupKeeper(t testing.TB) (*Keeper, sdk.Context) {
	storeKeys := sdk.NewKVStoreKeys(
		types.StoreKey,
		distrTypes.StoreKey,
		bankTypes.StoreKey,
		accountTypes.StoreKey,
		govTypes.StoreKey,
	)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	memStoreKeyGov := storetypes.NewMemoryStoreKey(govTypes.MemStoreKey)

	tkeys := sdk.NewTransientStoreKeys(paramsTypes.TStoreKey)

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
	sk.SetParams(ctx, stakingTypes.DefaultParams())
	gk := govKeeper.NewKeeper(codec.NewProtoCodec(registry), storeKeys[govTypes.StoreKey], memStoreKeyGov)
	dk := distrKeeper.NewKeeper(codec.NewProtoCodec(registry), storeKeys[distrTypes.StoreKey], pk.Subspace("distribution"),ak, bk, sk, accountTypes.FeeCollectorName, blacklistedAddrs)
	
	keeper := NewKeeper(
		codec.NewProtoCodec(registry),
		storeKeys[types.StoreKey],
		memStoreKey,
		dk,
		bk,
		ak,
		*gk,
	)

	return keeper, ctx
}
