package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
)

var (
	distrAcc             = supply.NewEmptyModuleAccount(types.ModuleName)
	TestFunder, _        = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
	TestDelegator, _     = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
	valAddr, _           = sdk.ValAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")
	pubKey               = ed25519.GenPrivKey().PubKey()
	TestValidator        = staking.NewValidator(valAddr, pubKey, staking.Description{})
	TestAmount           = sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(100)))
	TestBlockRewardsPool = sdk.NewDecCoinsFromCoins(sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(100000), Denom: "stake"})...)
	TestRewarRate        = sdk.NewDecWithPrec(12, 3)
)

func SetupTestInput(emptyPool bool) (cdc *codec.Codec, ctx sdk.Context, keeper Keeper, accKeeper auth.AccountKeeper, bankKeeper bank.BaseKeeper, stakeKeeper staking.Keeper) {
	memDB := db.NewMemDB()
	cdc = testCodec()

	keys := sdk.NewKVStoreKeys(
		params.StoreKey,
		auth.StoreKey,
		supply.StoreKey,
		staking.StoreKey,
		types.StoreKey,
		distr.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(params.TStoreKey, staking.TStoreKey)

	ms := store.NewCommitMultiStore(memDB)
	for _, key := range keys {
		ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, memDB)
	}

	for _, tkey := range tkeys {
		ms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, memDB)
	}

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

	pk := params.NewKeeper(cdc, keys[params.StoreKey], tkeys[params.TStoreKey])
	ak := auth.NewAccountKeeper(cdc, keys[auth.StoreKey], pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), blacklistedAddrs)

	// add module accounts to supply keeper
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		types.ModuleName:          {supply.Minter},
	}

	suk := supply.NewKeeper(cdc, keys[supply.StoreKey], ak, bk, maccPerms)
	suk.SetSupply(ctx, supply.NewSupply(sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(100000), Denom: "stake"})))
	sk := staking.NewKeeper(cdc, keys[staking.StoreKey], suk, pk.Subspace(staking.DefaultParamspace))
	sk.SetParams(ctx, staking.DefaultParams())

	dk := distr.NewKeeper(cdc, keys[distr.StoreKey], pk.Subspace(distr.DefaultParamspace), sk, suk, auth.FeeCollectorName, blacklistedAddrs)

	// set the distribution hooks on staking
	sk.SetHooks(dk.Hooks())

	k := NewKeeper(cdc, keys[types.StoreKey], dk, suk, keeper.govKeeper)

	if !emptyPool {
		pool, _ := TestBlockRewardsPool.TruncateDecimal()
		macc := k.VbrAccount(ctx)
		_ = macc.SetCoins(sdk.NewCoins(pool...))
		suk.SetModuleAccount(ctx, macc)
	}

	return cdc, ctx, k, ak, bk, sk
}

func testCodec() *codec.Codec {
	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	types.RegisterCodec(cdc) // distr

	cdc.Seal()

	return cdc
}
