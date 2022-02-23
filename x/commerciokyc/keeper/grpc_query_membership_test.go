package keeper

import (
	gocontext "context"

	"fmt"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

func (suite *KeeperTestSuite) TestGRPCMemberships() {

	queryClient := suite.queryClient
	app := suite.app
	ctx := suite.ctx

	app.CommercioKycKeeper.AssignMembership(ctx, testInviteSender, types.MembershipTypeGold, testTsp, testExpiration)
	membership := types.Membership{
		Owner:          testInviteSender.String(),
		TspAddress:     testTsp.String(),
		MembershipType: types.MembershipTypeGold,
		ExpiryAt:       &testExpiration,
	}
	var expectedRes []*types.Membership
	expectedRes = append(expectedRes, &membership)

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
				suite.Require().Equal(expectedRes, res.Memberships)

			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryMembership() {
	queryClient := suite.queryClient
	app := suite.app
	ctx := suite.ctx

	app.CommercioKycKeeper.AssignMembership(ctx, testInviteSender, types.MembershipTypeGold, testTsp, testExpiration)
	membership := types.Membership{
		Owner:          testInviteSender.String(),
		TspAddress:     testTsp.String(),
		MembershipType: types.MembershipTypeGold,
		ExpiryAt:       &testExpiration,
	}
	var expectedRes []*types.Membership
	expectedRes = append(expectedRes, &membership)

	var req *types.QueryMembershipRequest

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
		{
			"success",
			func() {
				req = &types.QueryMembershipRequest{Address: testInviteSender.String()}
			},
			true,
			func(res *types.QueryMembershipResponse) {
				suite.Require().True(res.Membership.Equals(membership))
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()

			res, err := queryClient.Membership(gocontext.Background(), req)

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
