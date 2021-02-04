package keeper

import (
	governmentKeeper "github.com/commercionetwork/commercionetwork/x/government/keeper"
	governmentTypes "github.com/commercionetwork/commercionetwork/x/government/types"
	"github.com/commercionetwork/commercionetwork/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
	"time"
)

// This function creates an environment to test modules
func SetupTestInput(setGovAddress bool) (cdc *codec.Codec, ctx sdk.Context, keeper Keeper) {

	memDB := db.NewMemDB()
	cdc = testCodec()

	// Store keys
	keyGovernment := sdk.NewKVStoreKey(governmentTypes.GovernmentStoreKey)
	keyUpgrade := sdk.NewKVStoreKey(upgrade.StoreKey)

	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(keyGovernment, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyUpgrade, sdk.StoreTypeIAVL, memDB)
	_ = ms.LoadLatestVersion()

	ctx = sdk.NewContext(ms, abci.Header{Time: time.Now(), Height: 1, ChainID: "test-chain-id"}, false, log.NewNopLogger())

	govk := governmentKeeper.NewKeeper(cdc, keyGovernment)
	upgk := upgrade.NewKeeper(map[int64]bool{}, keyUpgrade, cdc)
	customUpgradeKeeper := NewKeeper(cdc, govk, upgk)

	if setGovAddress {
		_ = govk.SetGovernmentAddress(ctx, governmentTestAddress)
	}
	return cdc, ctx, customUpgradeKeeper
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	types.RegisterCodec(cdc)
	cdc.Seal()

	return cdc
}

// Testing variables
var governmentTestAddress, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var notGovernmentAddress, _ = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")
