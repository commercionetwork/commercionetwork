package keeper

import (
	"time"

	commerciomintKeeper "github.com/commercionetwork/commercionetwork/x/commerciomint/keeper"
	commerciomintTypes "github.com/commercionetwork/commercionetwork/x/commerciomint/types"
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

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	government "github.com/commercionetwork/commercionetwork/x/government/keeper"
	governmentTypes "github.com/commercionetwork/commercionetwork/x/government/types"
)

// SetupTestInput function create an environment to test modules
func SetupTestInput() (sdk.Context, bankKeeper.Keeper, government.Keeper, Keeper) {

	memDB := db.NewMemDB()
	legacyAmino := codec.NewLegacyAmino()

	keys := sdk.NewKVStoreKeys(
		authTypes.StoreKey,
		bankTypes.StoreKey,
		paramsTypes.StoreKey,
		governmentTypes.StoreKey,
		types.StoreKey,
		commerciomintTypes.StoreKey,
	)
	tKeys := sdk.NewTransientStoreKeys(paramsTypes.TStoreKey)

	ms := store.NewCommitMultiStore(memDB)
	for _, key := range keys {
		ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, memDB)
	}
	for _, tkey := range tKeys {
		ms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, memDB)
	}
	_ = ms.LoadLatestVersion()

	app := simapp.Setup(false)
	cdc := app.AppCodec()

	ctx := sdk.NewContext(ms, tmproto.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	ctx = ctx.WithBlockTime(time.Now())
	//ctx = ctx.WithBlockTime(time.Now().Add(time.Hour * 1))

	maccPerms := map[string][]string{
		types.ModuleName:              {authTypes.Minter, authTypes.Burner},
		commerciomintTypes.ModuleName: {authTypes.Minter, authTypes.Burner},
	}

	pk := paramsKeeper.NewKeeper(cdc, legacyAmino, keys[paramsTypes.StoreKey], tKeys[paramsTypes.TStoreKey])
	ak := authKeeper.NewAccountKeeper(cdc, keys[authTypes.StoreKey], pk.Subspace(authTypes.DefaultParams().String()), authTypes.ProtoBaseAccount, maccPerms)
	bk := bankKeeper.NewBaseKeeper(cdc, keys[bankTypes.StoreKey], ak, pk.Subspace(bankTypes.DefaultParams().String()), nil)

	//bk.SetSupply(ctx, bankTypes.NewSupply(sdk.NewCoins(sdk.NewInt64Coin(stakeDenom, 1))))
	//ak.SetModuleAccount(ctx, authTypes.NewEmptyModuleAccount(types.ModuleName))
	govk := government.NewKeeper(cdc, keys[governmentTypes.StoreKey], keys[governmentTypes.StoreKey])

	mintAcc := authTypes.NewEmptyModuleAccount(commerciomintTypes.ModuleName, authTypes.Minter, authTypes.Burner)
	ak.SetModuleAccount(ctx, mintAcc)

	mk := commerciomintKeeper.NewKeeper(cdc, keys[commerciomintTypes.StoreKey], keys[commerciomintTypes.StoreKey], bk, ak, *govk, pk.Subspace(commerciomintTypes.ModuleName))
	memAcc := authTypes.NewEmptyModuleAccount(types.ModuleName, authTypes.Minter, authTypes.Burner)
	ak.SetModuleAccount(ctx, memAcc)
	bk.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(stakeDenom, 1)))

	k := NewKeeper(
		cdc,
		keys[types.StoreKey],
		keys[types.MemStoreKey],
		bk, *govk, ak, *mk)

	// TODO shall we drop the following?
	k.mintKeeper.UpdateParams(ctx, validCommercioMintParams)
	// k.MintKeeper.UpdateConversionRate(ctx, sdk.NewDecWithPrec(7, 1))

	k.GovKeeper.SetGovernmentAddress(ctx, testUser3)
	return ctx, bk, *govk, *k
}

/*func testCodec() *codec.Codec {
	var cdc = codec.NewLegacyAmino()

	authTypes.RegisterLegacyAminoCodec(cdc)
	bankTypes.RegisterLegacyAminoCodec(cdc)
	sdk.RegisterLegacyAminoCodec(cdc)

	types.RegisterCodec(cdc)

	cdc.Seal()

	return cdc
}*/

// Testing variables
var testUser, _ = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")
var testUser2, _ = sdk.AccAddressFromBech32("cosmos1h7tw92a66gr58pxgmf6cc336lgxadpjz5d5psf")
var testUser3, _ = sdk.AccAddressFromBech32("cosmos14lultfckehtszvzw4ehu0apvsr77afvyhgqhwh")
var testTsp, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")

var testInvalidMembership = "gren"

//var stableCreditDenom = "uccc"
var testExpiration = time.Now().Add(time.Hour * 24 * 365).UTC()
var testExpirationNegative = time.Now()
var depositStableCoin = sdk.NewCoins(sdk.NewInt64Coin(stableCreditDenom, 50000000))
var depositTestCoin = sdk.NewCoins(sdk.NewInt64Coin(stakeDenom, 50000000))
var yearBlocks = int64(4733640)

var testInviteSender, _ = sdk.AccAddressFromBech32("cosmos1005d6lt2wcfuulfpegz656ychljt3k3u4hn5my")

var validConversionRate = sdk.NewDecWithPrec(7, 1)
var validFreezePeriod time.Duration = 0
var validCommercioMintParams = commerciomintTypes.Params{
	ConversionRate: validConversionRate,
	FreezePeriod:   validFreezePeriod,
}
