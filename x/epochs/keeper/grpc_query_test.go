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

	suite.Require().Equal(epochInfosResponse.Epochs[1].Identifier, "hour")
	suite.Require().Equal(epochInfosResponse.Epochs[1].StartTime, chainStartTime)
	suite.Require().Equal(epochInfosResponse.Epochs[1].Duration, time.Hour)
	suite.Require().Equal(epochInfosResponse.Epochs[1].CurrentEpoch, int64(0))
	suite.Require().Equal(epochInfosResponse.Epochs[1].CurrentEpochStartTime, chainStartTime)
	suite.Require().Equal(epochInfosResponse.Epochs[1].EpochCountingStarted, false)
	suite.Require().Equal(epochInfosResponse.Epochs[1].CurrentEpochEnded, true)

	suite.Require().Equal(epochInfosResponse.Epochs[2].Identifier, "minute")
	suite.Require().Equal(epochInfosResponse.Epochs[2].StartTime, chainStartTime)
	suite.Require().Equal(epochInfosResponse.Epochs[2].Duration, time.Minute)
	suite.Require().Equal(epochInfosResponse.Epochs[2].CurrentEpoch, int64(0))
	suite.Require().Equal(epochInfosResponse.Epochs[2].CurrentEpochStartTime, chainStartTime)
	suite.Require().Equal(epochInfosResponse.Epochs[2].EpochCountingStarted, false)
	suite.Require().Equal(epochInfosResponse.Epochs[2].CurrentEpochEnded, true)

	suite.Require().Equal(epochInfosResponse.Epochs[3].Identifier, "month")
	suite.Require().Equal(epochInfosResponse.Epochs[3].StartTime, chainStartTime)
	suite.Require().Equal(epochInfosResponse.Epochs[3].Duration, time.Hour*24*30)
	suite.Require().Equal(epochInfosResponse.Epochs[3].CurrentEpoch, int64(0))
	suite.Require().Equal(epochInfosResponse.Epochs[3].CurrentEpochStartTime, chainStartTime)
	suite.Require().Equal(epochInfosResponse.Epochs[3].EpochCountingStarted, false)
	suite.Require().Equal(epochInfosResponse.Epochs[3].CurrentEpochEnded, true)

	suite.Require().Equal(epochInfosResponse.Epochs[4].Identifier, "week")
	suite.Require().Equal(epochInfosResponse.Epochs[4].StartTime, chainStartTime)
	suite.Require().Equal(epochInfosResponse.Epochs[4].Duration, time.Hour*24*7)
	suite.Require().Equal(epochInfosResponse.Epochs[4].CurrentEpoch, int64(0))
	suite.Require().Equal(epochInfosResponse.Epochs[4].CurrentEpochStartTime, chainStartTime)
	suite.Require().Equal(epochInfosResponse.Epochs[4].EpochCountingStarted, false)
	suite.Require().Equal(epochInfosResponse.Epochs[4].CurrentEpochEnded, true)
}

func (suite *KeeperTestSuite) TestQueryCurrentEpoch() {
	suite.SetupTest()
	queryClient := suite.queryClient

	_, err := queryClient.CurrentEpoch(gocontext.Background(), &types.QueryCurrentEpochRequest{Identifier: "invalid"})
	suite.Require().Error(err)

	currentEpochDay, err := queryClient.CurrentEpoch(gocontext.Background(), &types.QueryCurrentEpochRequest{Identifier: "day"})
	suite.Require().NoError(err)
	suite.Require().Equal(currentEpochDay.CurrentEpoch, int64(0))

	currentEpochHour, err := queryClient.CurrentEpoch(gocontext.Background(), &types.QueryCurrentEpochRequest{Identifier: "hour"})
	suite.Require().NoError(err)
	suite.Require().Equal(currentEpochHour.CurrentEpoch, int64(0))

	currentEpochMinute, err := queryClient.CurrentEpoch(gocontext.Background(), &types.QueryCurrentEpochRequest{Identifier: "minute"})
	suite.Require().NoError(err)
	suite.Require().Equal(currentEpochMinute.CurrentEpoch, int64(0))

	currentEpochmonth, err := queryClient.CurrentEpoch(gocontext.Background(), &types.QueryCurrentEpochRequest{Identifier: "month"})
	suite.Require().NoError(err)
	suite.Require().Equal(currentEpochmonth.CurrentEpoch, int64(0))

	currentEpochWeek, err := queryClient.CurrentEpoch(gocontext.Background(), &types.QueryCurrentEpochRequest{Identifier: "week"})
	suite.Require().NoError(err)
	suite.Require().Equal(currentEpochWeek.CurrentEpoch, int64(0))
}
