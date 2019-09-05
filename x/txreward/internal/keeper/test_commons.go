package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/txreward/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

type testInput struct {
	Cdc       *codec.Codec
	Ctx       sdk.Context
	TBRKeeper Keeper
}

var (
	distrAcc = supply.NewEmptyModuleAccount(types.ModuleName)
)

func setupTestInput() testInput {
	memDB := db.NewMemDB()
	cdc := testCodec()

	authKey := sdk.NewKVStoreKey("authCapKey")
	ibcKey := sdk.NewKVStoreKey("ibcCapKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	distrKey := sdk.NewKVStoreKey("distrKey")

	//TBR
	tbrStoreKey := sdk.NewKVStoreKey(types.BlockRewardsPoolPrefix)

	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, memDB)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(distrKey, sdk.StoreTypeIAVL, memDB)

	ms.MountStoreWithDB(tbrStoreKey, sdk.StoreTypeIAVL, memDB)

	_ = ms.LoadLatestVersion()
	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.GetAddress().String()] = true
	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
	blacklistedAddrs[bondPool.GetAddress().String()] = true
	blacklistedAddrs[distrAcc.GetAddress().String()] = true

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(cdc, authKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, blacklistedAddrs)

	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		types.ModuleName:          nil,
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
	}

	suk := supply.NewKeeper(cdc, keySupply, ak, bk, maccPerms)
	// set module accounts
	suk.SetModuleAccount(ctx, feeCollectorAcc)
	suk.SetModuleAccount(ctx, notBondedPool)
	suk.SetModuleAccount(ctx, bondPool)
	suk.SetModuleAccount(ctx, distrAcc)

	sk := staking.NewKeeper(cdc, keyStaking, tkeyStaking, suk, pk.Subspace(staking.DefaultParamspace), staking.DefaultCodespace)
	sk.SetParams(ctx, staking.DefaultParams())

	dk := distr.NewKeeper(cdc, distrKey, pk.Subspace(distr.DefaultParamspace), sk, suk, distr.DefaultCodespace, auth.FeeCollectorName, blacklistedAddrs)

	// set the distribution hooks on staking
	sk.SetHooks(dk.Hooks())

	tbrKeeper := NewKeeper(tbrStoreKey, bk, sk, dk, cdc)

	return testInput{Cdc: cdc, Ctx: ctx, TBRKeeper: tbrKeeper}
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterInterface((*auth.Account)(nil), nil)

	cdc.Seal()

	return cdc
}

var TestUtils = setupTestInput()

var TestFunder, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var TestAmount = sdk.Coin{
	Denom:  "ucommercio",
	Amount: sdk.NewInt(10),
}
