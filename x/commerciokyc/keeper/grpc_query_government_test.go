package keeper

// import (
// 	gocontext "context"
// 	"fmt"

// 	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// )

// // Tsps query tests
// func (suite *KeeperTestSuite) TestGRPCTsps() {

// 	queryClient := suite.queryClient
// 	app := suite.app
// 	ctx := suite.ctx
// 	app.CommercioKycKeeper.AddTrustedServiceProvider(ctx, testTsp)
// 	var expectedRes []string
// 	expectedRes = append(expectedRes, testTsp.String())

// 	var req *types.QueryTspsRequest

// 	testCases := []struct {
// 		msg      string
// 		malleate func()
// 		expPass  bool
// 	}{
// 		{
// 			"valid request",
// 			func() {
// 				req = &types.QueryTspsRequest{}
// 			},
// 			true,
// 		},
// 	}

// 	for _, testCase := range testCases {
// 		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
// 			testCase.malleate()

// 			res, err := queryClient.Tsps(gocontext.Background(), req)

// 			if testCase.expPass {
// 				suite.Require().NoError(err)
// 				suite.Require().NotNil(res)
// 				suite.Require().Equal(expectedRes, res.Tsps)
// 			} else {
// 				suite.Require().Error(err)
// 				suite.Require().Nil(res)
// 			}
// 		})
// 	}
// }

// // Funds query tests
// func (suite *KeeperTestSuite) TestGRPCFunds() {

// 	queryClient := suite.queryClient
// 	app := suite.app
// 	ctx := suite.ctx

// 	coins := sdk.NewCoins(sdk.NewCoin(stakeDenom, math.NewInt(1000)))
// 	app.CommercioKycKeeper.SetLiquidityPoolToAccount(ctx, coins)

// 	var req *types.QueryFundsRequest

// 	testCases := []struct {
// 		msg      string
// 		malleate func()
// 		expPass  bool
// 	}{
// 		{
// 			"valid request",
// 			func() {
// 				req = &types.QueryFundsRequest{}
// 			},
// 			true,
// 		},
// 	}

// 	for _, testCase := range testCases {
// 		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
// 			testCase.malleate()

// 			res, err := queryClient.Funds(gocontext.Background(), req)

// 			if testCase.expPass {
// 				suite.Require().NoError(err)
// 				suite.Require().NotNil(res)
// 				coins = coins.Add(sdk.NewCoin(stakeDenom, math.NewInt(1)))
// 				suite.Require().Equal(coins, res.Funds)
// 			} else {
// 				suite.Require().Error(err)
// 				suite.Require().Nil(res)
// 			}
// 		})
// 	}
// }
