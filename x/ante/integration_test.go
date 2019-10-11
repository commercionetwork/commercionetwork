package ante_test

import (
	"github.com/commercionetwork/commercionetwork/x/pricefeed"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// returns context and app with params set on account keeper
func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context, pricefeed.Keeper) {
	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	ctx = ctx.WithBlockHeight(1)

	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	app.MountKVStores(sdk.NewKVStoreKeys(pricefeed.StoreKey))

	pfk := pricefeed.NewKeeper(app.Codec(), app.GetKey(pricefeed.StoreKey))

	return app, ctx, pfk
}
