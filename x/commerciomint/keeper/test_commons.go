package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	params "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramsType "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	db "github.com/tendermint/tm-db"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	government "github.com/commercionetwork/commercionetwork/x/government/keeper"

	governmentTypes "github.com/commercionetwork/commercionetwork/x/government/types"
)

func SetupTestInput() (sdk.Context, bankKeeper.Keeper, government.Keeper, Keeper) {
	memDB := db.NewMemDB()
	app := simapp.Setup(false)
	cdc := app.AppCodec()

	keys := sdk.NewKVStoreKeys(
		authTypes.StoreKey,
		bankTypes.StoreKey,
		paramsType.StoreKey,
		governmentTypes.StoreKey,
		types.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramsType.TStoreKey)

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
	maccPerms := map[string][]string{
		types.ModuleName: {authTypes.Minter, authTypes.Burner},
	}

	pk := params.NewKeeper(cdc, legacyCodec, keys[paramsType.StoreKey], tkeys[paramsType.TStoreKey])
	ak := authKeeper.NewAccountKeeper(cdc, keys[authTypes.StoreKey], pk.Subspace(authTypes.DefaultParams().String()), authTypes.ProtoBaseAccount, maccPerms)
	bk := bankKeeper.NewBaseKeeper(cdc, keys[bankTypes.StoreKey], ak, pk.Subspace(bankTypes.DefaultParams().String()), nil)

	govkeeper := government.NewKeeper(cdc, keys[governmentTypes.StoreKey], keys[governmentTypes.StoreKey])

	bk.SetSupply(ctx, bankTypes.NewSupply(sdk.NewCoins(*testEtp.Credits)))

	memAcc := authTypes.NewEmptyModuleAccount(types.ModuleName, authTypes.Minter, authTypes.Burner)
	ak.SetModuleAccount(ctx, memAcc)

	mintK := NewKeeper(
		cdc,
		keys[types.StoreKey],
		keys[types.MemStoreKey],
		bk, ak, *govkeeper)

	err := mintK.SetConversionRate(ctx, sdk.NewDec(2))
	if err != nil {
		panic(err)
	}
	err = mintK.SetFreezePeriod(ctx, 0)
	if err != nil {
		panic(err)
	}

	return ctx, bk, *govkeeper, *mintK
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
	time.Now().UTC(),
	sdk.NewDec(2),
)

var fakeEtp = types.NewPosition(
	testEtpOwner,
	sdk.NewInt(100),
	sdk.NewCoin("uccc", sdk.NewInt(50)),
	fakeID,
	time.Now().UTC(),
	sdk.NewDec(2),
)

var testLiquidityPool = sdk.NewCoins(sdk.NewInt64Coin(testLiquidityDenom, 10000))
