package keeper_test

import (
	gocontext "context"
	"time"

	"github.com/commercionetwork/commercionetwork/x/epochs/types"
)

func (suite *KeeperTestSuite) TestQueryEpochInfos() {
	suite.SetupTest()
	queryClient := suite.queryClient

	chainStartTime := suite.ctx.BlockTime()

	// Invalid param
	epochInfosResponse, err := queryClient.EpochInfos(gocontext.Background(), &types.QueryEpochsInfoRequest{})
	suite.Require().NoError(err)
	suite.Require().Len(epochInfosResponse.Epochs, 5)

	// check if EpochInfos are correct
	suite.Require().Equal(epochInfosResponse.Epochs[0].Identifier, "day")
	suite.Require().Equal(epochInfosResponse.Epochs[0].StartTime, chainStartTime)
	suite.Require().Equal(epochInfosResponse.Epochs[0].Duration, time.Hour*24)
	suite.Require().Equal(epochInfosResponse.Epochs[0].CurrentEpoch, int64(0))
	suite.Require().Equal(epochInfosResponse.Epochs[0].CurrentEpochStartTime, chainStartTime)
	suite.Require().Equal(epochInfosResponse.Epochs[0].EpochCountingStarted, false)
	suite.Require().Equal(epochInfosResponse.Epochs[0].CurrentEpochEnded, true)
	suite.Require().Equal(epochInfosResponse.Epochs[4].Identifier, "week")
	suite.Require().Equal(epochInfosResponse.Epochs[4].StartTime, chainStartTime)
	suite.Require().Equal(epochInfosResponse.Epochs[4].Duration, time.Hour*24*7)
	suite.Require().Equal(epochInfosResponse.Epochs[4].CurrentEpoch, int64(0))
	suite.Require().Equal(epochInfosResponse.Epochs[4].CurrentEpochStartTime, chainStartTime)
	suite.Require().Equal(epochInfosResponse.Epochs[4].EpochCountingStarted, false)
	suite.Require().Equal(epochInfosResponse.Epochs[4].CurrentEpochEnded, true)
}
