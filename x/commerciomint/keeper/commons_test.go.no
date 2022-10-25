package keeper

import (
	"context"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramsKeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	db "github.com/tendermint/tm-db"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"

	governmentKeeper "github.com/commercionetwork/commercionetwork/x/government/keeper"
	governmentTypes "github.com/commercionetwork/commercionetwork/x/government/types"
)

func SetupTestInput() (sdk.Context, bankKeeper.Keeper, governmentKeeper.Keeper, Keeper) {
	memDB := db.NewMemDB()
	app := simapp.Setup(false)
	cdc := app.AppCodec()

	keys := sdk.NewKVStoreKeys(
		authTypes.StoreKey,
		bankTypes.StoreKey,
		paramsTypes.StoreKey,
		governmentTypes.StoreKey,
		types.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramsTypes.TStoreKey)

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

	// legacyCodec := codec.NewLegacyAmino()
	maccPerms := map[string][]string{
		types.ModuleName: {authTypes.Minter, authTypes.Burner},
	}

	pk := paramsKeeper.NewKeeper(cdc, codec.NewLegacyAmino(), keys[paramsTypes.StoreKey], tkeys[paramsTypes.TStoreKey])

	ak := authKeeper.NewAccountKeeper(cdc, keys[authTypes.StoreKey], pk.Subspace(authTypes.DefaultParams().String()), authTypes.ProtoBaseAccount, maccPerms)
	bk := bankKeeper.NewBaseKeeper(cdc, keys[bankTypes.StoreKey], ak, pk.Subspace(bankTypes.DefaultParams().String()), nil)

	govkeeper := governmentKeeper.NewKeeper(cdc, keys[governmentTypes.StoreKey], keys[governmentTypes.StoreKey])

	//bk.SetSupply(ctx, bankTypes.NewSupply(sdk.NewCoins(*testEtp.Credits)))

	memAcc := authTypes.NewEmptyModuleAccount(types.ModuleName, authTypes.Minter, authTypes.Burner)
	ak.SetModuleAccount(ctx, memAcc)

	mintK := NewKeeper(
		cdc,
		keys[types.StoreKey],
		keys[types.MemStoreKey],
		bk, ak, *govkeeper,
		pk.Subspace(types.ModuleName),
	)

	mintK.UpdateParams(ctx, validParams)

	return ctx, bk, *govkeeper, *mintK
}

func SetupMsgServer() (context.Context, bankKeeper.Keeper, governmentKeeper.Keeper, Keeper, types.MsgServer) {
	ctx, bk, gk, k := SetupTestInput()

	wctx := sdk.WrapSDKContext(ctx)

	return wctx, bk, gk, k, NewMsgServerImpl(k)
}

// ----------------------
// --- Test variables
// ----------------------

var testEtpOwner, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var ownerAnother, _ = sdk.AccAddressFromBech32("cosmos14lultfckehtszvzw4ehu0apvsr77afvyhgqhwh")
var government = ownerAnother

var testID = "2908006A-93D4-4517-A8F5-393EEEBDDB61"
var halfCoinSub = sdk.NewCoin(types.CreditsDenom, sdk.NewInt(10))

var testEtp = types.NewPosition(
	testEtpOwner,
	sdk.NewInt(100),
	validDepositCoin,
	testID,
	time.Now().UTC(),
	sdk.NewDec(2),
)

var validParams = types.Params{
	ConversionRate: validConversionRate,
	FreezePeriod:   validFreezePeriod,
}

var validDepositCoin = sdk.NewCoin(types.CreditsDenom, sdk.NewInt(50))
var inValidDepositCoin = sdk.NewCoin(types.BondDenom, sdk.NewInt(10))
var validBurnCoin = inValidDepositCoin
var inValidBurnCoin = validDepositCoin
var validConversionRate = sdk.NewDec(2)
var invalidConversionRate = sdk.NewDec(-1)

var zeroUCCC = sdk.NewCoin(types.CreditsDenom, sdk.ZeroInt())

var validFreezePeriod time.Duration = 0
var invalidFreezePeriod = -time.Minute

var testLiquidityPool = sdk.NewCoins(sdk.NewInt64Coin(types.BondDenom, 10000))

var testEtp1, testEtp2, testEtpAnotherOwner types.Position

func init() {
	testEtp1 = testEtp
	testEtp1.ID = strings.Replace(testEtp1.ID, "0", "A", 1)
	testEtp2 = testEtp
	testEtp2.ID = strings.Replace(testEtp1.ID, "0", "B", 1)
	testEtpAnotherOwner = testEtp
	testEtpAnotherOwner.ID = strings.Replace(testEtp1.ID, "0", "C", 1)
	testEtpAnotherOwner.Owner = ownerAnother.String()
}
