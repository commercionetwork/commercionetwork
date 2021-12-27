package keeper_test

import (
	gocontext "context"

	"fmt"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestGRPCMemberships() {

	ctx, _, _, k := SetupTestInput()
	queryClient := suite.queryClient

	k.AssignMembership(ctx, testInviteSender, types.MembershipTypeGold, testTsp, testExpiration)
	var membershipsRes types.Memberships
	_ = membershipsRes

	var req *types.QueryMembershipsRequest

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryMembershipsRequest{}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			res, err := queryClient.Memberships(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryMembership() {
	var (
		req *types.QueryMembershipRequest
	)
	//_, _, addr := testdata.KeyTestPubAddr()

	testCases := []struct {
		msg       string
		malleate  func()
		expPass   bool
		posttests func(res *types.QueryMembershipResponse)
	}{
		{
			"empty request",
			func() {
				req = &types.QueryMembershipRequest{}
			},
			false,
			func(res *types.QueryMembershipResponse) {},
		},
		{
			"invalid request",
			func() {
				req = &types.QueryMembershipRequest{Address: ""}
			},
			false,
			func(res *types.QueryMembershipResponse) {},
		},
		{
			"account not found",
			func() {
				req = &types.QueryMembershipRequest{Address: testUser.String()}
			},
			false,
			func(res *types.QueryMembershipResponse) {},
		},
		/*{
			"success",
			func() {
				suite.app.AccountKeeper.SetAccount(suite.ctx,
					suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr))
				req = &types.QueryMembershipResponse{Address: addr.String()}
			},
			true,
			func(res *types.QueryMembershipResponse) {
				var newAccount types.AccountI
				err := suite.app.InterfaceRegistry().UnpackAny(res.Account, &newAccount)
				suite.Require().NoError(err)
				suite.Require().NotNil(newAccount)
				suite.Require().True(addr.Equals(newAccount.GetAddress()))
			},
		},*/
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.ctx)

			res, err := suite.queryClient.Membership(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}

			tc.posttests(res)
		})
	}
}
