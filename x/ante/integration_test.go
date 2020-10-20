package ante_test

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/commercionetwork/commercionetwork/simapp"
)

// returns context and app with params set on account keeper
func createTestApp(isCheckTx bool, isBlockZero bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)

	header := abci.Header{}

	if !isBlockZero {
		header.Height = 1
	}

	ctx := app.BaseApp.NewContext(isCheckTx, header)
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())

	return app, ctx
}
