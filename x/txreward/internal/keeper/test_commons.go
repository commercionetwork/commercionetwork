package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/txreward/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

var (
	distrAcc = supply.NewEmptyModuleAccount(types.ModuleName)
)

var (
	multiPerm  = "multiple permissions account"
	randomPerm = "random permission"
	holder     = "holder"
)

var addr, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var valAddr, _ = sdk.ValAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")
var pubKey = ed25519.GenPrivKey().PubKey()
var TestValidator = staking.NewValidator(valAddr, pubKey, staking.Description{})

var TestFunder = types.Funder{Address: addr}

var TestAmount = sdk.Coin{
	Denom:  "ucommercio",
	Amount: sdk.NewInt(100),
}

var coin = sdk.Coin{Amount: sdk.NewInt(10000000000000000), Denom: types.DefaultBondDenom}
var coins = sdk.NewCoins(coin)
var TestBlockRewardsPool = types.BlockRewardsPool{
	Funds: sdk.NewDecCoins(coins),
}

var TestFunders = types.Funders{TestFunder}

func SetupTestInput() (cdc *codec.Codec, ctx sdk.Context, keeper Keeper) {
	memDB := db.NewMemDB()
	cdc = testCodec()

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

	ctx = sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(cdc, authKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, blacklistedAddrs)

	// add module accounts to supply keeper
	maccPerms := simapp.GetMaccPerms()
	maccPerms[holder] = nil
	maccPerms[supply.Burner] = []string{supply.Burner}
	maccPerms[supply.Minter] = []string{supply.Minter}
	maccPerms[multiPerm] = []string{supply.Burner, supply.Minter, supply.Staking}
	maccPerms[randomPerm] = []string{"random"}

	suk := supply.NewKeeper(cdc, keySupply, ak, bk, maccPerms)
	sk := staking.NewKeeper(cdc, keyStaking, tkeyStaking, suk, pk.Subspace(staking.DefaultParamspace), staking.DefaultCodespace)
	sk.SetParams(ctx, staking.DefaultParams())

	dk := distr.NewKeeper(cdc, distrKey, pk.Subspace(distr.DefaultParamspace), sk, suk, distr.DefaultCodespace, auth.FeeCollectorName, blacklistedAddrs)

	// set the distribution hooks on staking
	sk.SetHooks(dk.Hooks())

	tbrKeeper := NewKeeper(tbrStoreKey, bk, sk, dk, cdc)

	return cdc, ctx, tbrKeeper
}

func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})

	app.AccountKeeper.SetParams(ctx, auth.DefaultParams())
	app.BankKeeper.SetSendEnabled(ctx, true)

	return app, ctx
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterInterface((*auth.Account)(nil), nil)

	cdc.Seal()

	return cdc
}
