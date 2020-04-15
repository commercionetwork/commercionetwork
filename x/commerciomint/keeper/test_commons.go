package keeper

import (
	governmentKeeper "github.com/commercionetwork/commercionetwork/x/government/keeper"
	governmentTypes "github.com/commercionetwork/commercionetwork/x/government/types"
	pricefeedKeeper "github.com/commercionetwork/commercionetwork/x/pricefeed/keeper"
	pricefeedTypes "github.com/commercionetwork/commercionetwork/x/pricefeed/types"
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

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	creditrisk "github.com/commercionetwork/commercionetwork/x/creditrisk/types"
)

func SetupTestInput() (sdk.Context, bank.Keeper, pricefeedKeeper.Keeper, governmentKeeper.Keeper, supply.Keeper, Keeper) {
	memDB := db.NewMemDB()
	cdc := testCodec()

	keys := sdk.NewKVStoreKeys(
		auth.StoreKey,
		params.StoreKey,
		supply.StoreKey,
		pricefeedTypes.StoreKey,
		governmentTypes.StoreKey,
		creditrisk.StoreKey,
		types.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(params.TStoreKey)

	ms := store.NewCommitMultiStore(memDB)
	for _, key := range keys {
		ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, memDB)
	}
	for _, tkey := range tkeys {
		ms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, nil)
	}
	_ = ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	pk := params.NewKeeper(cdc, keys[params.StoreKey], tkeys[params.TStoreKey])
	ak := auth.NewAccountKeeper(cdc, keys[auth.StoreKey], pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), nil)
	maccPerms := map[string][]string{
		types.ModuleName:      {supply.Minter, supply.Burner},
		creditrisk.ModuleName: {supply.Minter, supply.Burner},
	}
	sk := supply.NewKeeper(cdc, keys[supply.StoreKey], ak, bk, maccPerms)

	govkeeper := governmentKeeper.NewKeeper(cdc, keys[governmentTypes.StoreKey])
	pfk := pricefeedKeeper.NewKeeper(cdc, keys[pricefeedTypes.StoreKey], govkeeper)

	mintK := NewKeeper(cdc, keys[types.StoreKey], sk, pfk, govkeeper)

	// Set initial supply
	sk.SetSupply(ctx, supply.NewSupply(sdk.NewCoins(testCdp.Credits)))

	// Set module accounts
	mintAcc := supply.NewEmptyModuleAccount(types.ModuleName, supply.Minter, supply.Burner)
	mintK.supplyKeeper.SetModuleAccount(ctx, mintAcc)

	// Set the credits denom
	mintK.SetCreditsDenom(ctx, testCreditsDenom)
	// Set cdp collateral rate
	mintK.SetCollateralRate(ctx, sdk.NewDec(2))

	return ctx, bk, pfk, govkeeper, sk, mintK
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	bank.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	pricefeedTypes.RegisterCodec(cdc)
	governmentTypes.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	types.RegisterCodec(cdc)

	cdc.Seal()
	return cdc
}

// ----------------------
// --- Test variables
// ----------------------

var testCreditsDenom = "stake"
var testLiquidityDenom = "ucommercio"
var testCdpOwner, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")

var testCdp = types.NewPosition(
	testCdpOwner,
	sdk.NewCoins(sdk.NewCoin(testLiquidityDenom, sdk.NewInt(100))),
	sdk.NewCoin(testCreditsDenom, sdk.NewInt(50)),
	10,
)

var testLiquidityPool = sdk.NewCoins(sdk.NewInt64Coin(testLiquidityDenom, 10000))
