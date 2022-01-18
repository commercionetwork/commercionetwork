package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"

	//commsimapp "github.com/commercionetwork/commercionetwork/testutil/simapp"

	sdk "github.com/cosmos/cosmos-sdk/types"

	//"github.com/cosmos/cosmos-sdk/x/auth/types"
	//authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/keeper"
	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

type kycSimApp struct {
	*simapp.SimApp
	CommercioKycKeeper keeper.Keeper
}

type KeeperTestSuite struct {
	suite.Suite

	app *kycSimApp
	ctx sdk.Context

	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(true)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.CommercioKycKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
}

func createTestApp(isCheckTx bool) (*kycSimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)
	ctx, _, _, k := SetupTestInput()
	outApp := &kycSimApp{
		app,
		k,
	}

	//ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	//app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	return outApp, ctx
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
