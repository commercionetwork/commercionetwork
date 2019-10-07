package ante_test

import (
	"github.com/commercionetwork/commercionetwork/x/pricefeed"
	"github.com/cosmos/cosmos-sdk/store"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// returns context and app with params set on account keeper
func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context, pricefeed.Keeper) {
	app := simapp.Setup(isCheckTx)

	// Setup the pricefeed keeper
	memDB := db.NewMemDB()

	govKey := sdk.NewKVStoreKey("government")
	pricefeedKey := sdk.NewKVStoreKey("pricefeed")

	ms := store.NewCommitMultiStore(memDB)
	app.MountStoreWithDB(govKey, sdk.StoreTypeIAVL, memDB)
	app.MountStoreWithDB(pricefeedKey, sdk.StoreTypeIAVL, memDB)

	_ = ms.LoadLatestVersion()
	pfk := pricefeed.NewKeeper(app.Codec(), pricefeedKey)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, isCheckTx, log.NewNopLogger())

	return app, ctx, pfk
}
