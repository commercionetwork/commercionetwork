package keeper_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/app"
	"github.com/commercionetwork/commercionetwork/testutil/simapp"
	"github.com/commercionetwork/commercionetwork/x/epochs"
	"github.com/commercionetwork/commercionetwork/x/epochs/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app         *app.App
	ctx         sdk.Context
	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = simapp.New("")
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})

	genState := types.DefaultGenesis()
	epochs.InitGenesis(suite.ctx, suite.app.EpochsKeeper, *genState)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.EpochsKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
