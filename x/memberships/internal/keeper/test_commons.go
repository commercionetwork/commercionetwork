package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/accreditations"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

//This function create an environment to test modules
func SetupTestInput() (*codec.Codec, sdk.Context, accreditations.Keeper, bank.Keeper, Keeper) {

	memDB := db.NewMemDB()
	cdc := testCodec()

	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	authKey := sdk.NewKVStoreKey("authCapKey")
	ibcKey := sdk.NewKVStoreKey("ibcCapKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	nftKey := sdk.NewKVStoreKey("nft")

	accreditationKey := sdk.NewKVStoreKey("accreditations")
	storeKey := sdk.NewKVStoreKey("memberships")

	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, memDB)
	ms.MountStoreWithDB(accreditationKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(nftKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, memDB)

	_ = ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, map[string]bool{})
	accreditationKeeper := accreditations.NewKeeper(accreditationKey, bankKeeper, cdc)

	nftk := nft.NewKeeper(cdc, nftKey)
	memk := NewKeeper(cdc, storeKey, nftk)

	return cdc, ctx, accreditationKeeper, bankKeeper, memk
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	nft.RegisterCodec(cdc)
	accreditations.RegisterCodec(cdc)

	cdc.Seal()

	return cdc
}

// Test variables
var TestUserAddress, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var TestMembershipType = "bronze"
