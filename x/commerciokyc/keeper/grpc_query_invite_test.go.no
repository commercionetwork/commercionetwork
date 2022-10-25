package keeper

import (
	gocontext "context"
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

func (suite *KeeperTestSuite) TestGRPCInvites() {

	queryClient := suite.queryClient
	app := suite.app
	ctx := suite.ctx
	// Setup membership before invite user
	app.CommercioKycKeeper.AssignMembership(ctx, testUser, types.MembershipTypeGold, testTsp, testExpiration)
	app.CommercioKycKeeper.SetInvite(ctx, testInviteSender, testUser)
	app.CommercioKycKeeper.SetInvite(ctx, testUser3, testUser)

	invite := types.Invite{
		Sender:           testUser.String(),
		SenderMembership: types.MembershipTypeGold,
		User:             testInviteSender.String(),
		Status:           uint64(types.InviteStatusPending),
	}
	invite2 := types.Invite{
		Sender:           testUser.String(),
		SenderMembership: types.MembershipTypeGold,
		User:             testUser3.String(),
		Status:           uint64(types.InviteStatusPending),
	}

	var expectedRes []*types.Invite
	expectedRes = append(expectedRes, &invite)
	expectedRes = append(expectedRes, &invite2)

	var req *types.QueryInvitesRequest

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryInvitesRequest{}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			res, err := queryClient.Invites(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(expectedRes, res.Invites)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCInvite() {

	queryClient := suite.queryClient
	app := suite.app
	ctx := suite.ctx
	// Setup membership before invite user
	app.CommercioKycKeeper.AssignMembership(ctx, testUser, types.MembershipTypeGold, testTsp, testExpiration)
	app.CommercioKycKeeper.SetInvite(ctx, testInviteSender, testUser)

	invite := types.Invite{
		Sender:           testUser.String(),
		SenderMembership: types.MembershipTypeGold,
		User:             testInviteSender.String(),
		Status:           uint64(types.InviteStatusPending),
	}
	expectedRes := &invite

	var req *types.QueryInviteRequest

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				req = &types.QueryInviteRequest{Address: testInviteSender.String()}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			res, err := queryClient.Invite(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(expectedRes, res.Invite)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}
