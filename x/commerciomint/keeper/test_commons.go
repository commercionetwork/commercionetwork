package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authType "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	params "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramsType "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	db "github.com/tendermint/tm-db"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	government "github.com/commercionetwork/commercionetwork/x/government/keeper"

	governmentTypes "github.com/commercionetwork/commercionetwork/x/government/types"
)

func SetupTestInput() (sdk.Context, bank.Keeper, government.Keeper, Keeper) {
	memDB := db.NewMemDB()
	cdc := testCodec()

	keys := sdk.NewKVStoreKeys(
		authType.StoreKey,
		paramsType.StoreKey,
		governmentTypes.StoreKey,
		types.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramsType.QuerierRoute)

	ms := store.NewCommitMultiStore(memDB)
	for _, key := range keys {
		ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, memDB)
	}
	for _, tkey := range tkeys {
		ms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, nil)
	}
	_ = ms.LoadLatestVersion()

	var header tmproto.Header
	header.ChainID = "test-chain-id"

	ctx := sdk.NewContext(ms, header, false, log.NewNopLogger()).WithBlockTime(time.Now())

	legacyCodec := codec.NewLegacyAmino()
	pk := params.NewKeeper(cdc, legacyCodec, keys[paramsType.StoreKey], tkeys[paramsType.TStoreKey])
	/*ak := auth.NewAccountKeeper(cdc, keys[authType.StoreKey], pk.Subspace(authType.Subspace.Name()), authType.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak, pk.Subspace(bank.DefaultParamspace), nil)
	maccPerms := map[string][]string{
		types.ModuleName: {supply.Minter, supply.Burner},
	}
	sk := supply.NewKeeper(cdc, keys[supply.StoreKey], ak, bk, maccPerms)

	govkeeper := government.NewKeeper(cdc, keys[governmentTypes.StoreKey])

	mintK := NewKeeper(cdc, keys[types.StoreKey], sk, govkeeper)

	// Set initial supply
	sk.SetSupply(ctx, supply.NewSupply(sdk.NewCoins(testEtp.Credits)))

	// Set module accounts
	mintAcc := supply.NewEmptyModuleAccount(types.ModuleName, supply.Minter, supply.Burner)
	mintK.supplyKeeper.SetModuleAccount(ctx, mintAcc)

	// Set etp collateral rate
	err := mintK.SetConversionRate(ctx, sdk.NewDec(2))
	if err != nil {
		panic(err)
	}

	// Set etp freeze period
	err = mintK.SetFreezePeriod(ctx, 0)
	if err != nil {
		panic(err)
	}*/

	return ctx, bk, govkeeper, sk, mintK
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	bank.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
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

var testLiquidityDenom = "ucommercio"
var testEtpOwner, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var testID = "2908006A-93D4-4517-A8F5-393EEEBDDB61"
var fakeID = "2908006A-93D4-4517-A8F5-393EEEBDDB61"
var halfCoinSub = sdk.NewCoin("uccc", sdk.NewInt(10))

var testEtp = types.NewPosition(
	testEtpOwner,
	sdk.NewInt(100),
	sdk.NewCoin("uccc", sdk.NewInt(50)),
	testID,
	time.Now(),
	sdk.NewDec(2),
)

var fakeEtp = types.NewPosition(
	testEtpOwner,
	sdk.NewInt(100),
	sdk.NewCoin("uccc", sdk.NewInt(50)),
	fakeID,
	time.Now(),
	sdk.NewDec(2),
)

var testLiquidityPool = sdk.NewCoins(sdk.NewInt64Coin(testLiquidityDenom, 10000))
