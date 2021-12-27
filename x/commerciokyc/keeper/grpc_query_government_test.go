package keeper_test

import (
	gocontext "context"
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Tsps query tests
func (suite *KeeperTestSuite) TestGRPCTsps() {

	ctx, _, _, k := SetupTestInput()
	queryClient := suite.queryClient
	k.AddTrustedServiceProvider(ctx, testTsp)
	var tspsRes ctypes.Strings
	tspsRes.AppendIfMissing(testTsp.String())

	var req *types.QueryTspsRequest

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryTspsRequest{}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			tspRes, err := queryClient.Tsps(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(tspRes)
				suite.Require().Equal(tspRes.Tsps, tspsRes)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(tspRes)
			}
		})
	}
}

// Funds query tests
func (suite *KeeperTestSuite) TestGRPCFunds() {

	ctx, _, _, k := SetupTestInput()
	queryClient := suite.queryClient
	coins := sdk.NewCoins(sdk.NewCoin("somecoin", sdk.NewInt(1000)))
	k.DepositIntoPool(ctx, testTsp, coins)

	var tspsRes ctypes.Strings
	tspsRes.AppendIfMissing(testTsp.String())

	var req *types.QueryFundsRequest

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryFundsRequest{}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			fundsRes, err := queryClient.Funds(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(fundsRes)
				suite.Require().Equal(fundsRes.Funds, coins)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(fundsRes)
			}
		})
	}
}
