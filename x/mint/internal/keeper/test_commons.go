package keeper

import (
	"time"

	"github.com/commercionetwork/commercionetwork/x/government"
	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
	"github.com/commercionetwork/commercionetwork/x/pricefeed"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

var TestCreditsDenom = "stake"
var TestLiquidityDenom = "ucommercio"
var TestOwner, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var timezone, _ = time.LoadLocation("UTC")
var TestTimestamp = time.Date(1990, 01, 01, 20, 20, 00, 0, timezone)

var TestCdpRequest = types.CdpRequest{
	Signer:          TestOwner,
	DepositedAmount: sdk.NewCoins(sdk.NewCoin(TestLiquidityDenom, sdk.NewInt(100))),
	Timestamp:       TestTimestamp,
}

var TestCdp = types.Cdp{
	Owner:           TestOwner,
	DepositedAmount: sdk.NewCoins(sdk.NewCoin(TestLiquidityDenom, sdk.NewInt(100))),
	CreditsAmount:   sdk.NewCoins(sdk.NewCoin(TestCreditsDenom, sdk.NewInt(50))),
	Timestamp:       TestTimestamp,
}

var TestLiquidityPool = sdk.Coins{sdk.NewInt64Coin(TestLiquidityDenom, 10000)}

func SetupTestInput() (*codec.Codec, sdk.Context, bank.Keeper, pricefeed.Keeper, Keeper) {
	memDB := db.NewMemDB()
	cdc := testCodec()

	keys := sdk.NewKVStoreKeys(
		auth.StoreKey,
		params.StoreKey,
		supply.StoreKey,
		pricefeed.StoreKey,

		types.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(
		staking.TStoreKey,
		params.TStoreKey,
	)

	ms := store.NewCommitMultiStore(memDB)
	for _, key := range keys {
		ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, memDB)
	}
	for _, tkey := range tkeys {
		ms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, nil)
	}

	_ = ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	pk := params.NewKeeper(cdc, keys[params.StoreKey], tkeys[params.TStoreKey], params.DefaultCodespace)
	ak := auth.NewAccountKeeper(cdc, keys[auth.StoreKey], pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, nil)
	maccPerms := map[string][]string{
		types.ModuleName: {supply.Minter, supply.Burner},
	}
	sk := supply.NewKeeper(cdc, keys[supply.StoreKey], ak, bk, maccPerms)
	pfk := pricefeed.NewKeeper(cdc, keys[pricefeed.StoreKey])

	mintK := NewKeeper(cdc, keys[types.StoreKey], sk, pfk)

	// Set initial supply
	sk.SetSupply(ctx, supply.NewSupply(TestCdp.CreditsAmount))

	// Set module accounts
	mintAcc := supply.NewEmptyModuleAccount(types.ModuleName, supply.Minter, supply.Burner)
	mintK.supplyKeeper.SetModuleAccount(ctx, mintAcc)

	// Set the credits denom
	mintK.SetCreditsDenom(ctx, TestCreditsDenom)

	return cdc, ctx, bk, pfk, mintK
}

func testCodec() *codec.Codec {
	var cdc = codec.New()
	sdk.RegisterCodec(cdc)

	bank.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	pricefeed.RegisterCodec(cdc)
	government.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	cdc.Seal()
	return cdc
}
