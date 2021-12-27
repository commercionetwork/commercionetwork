package keeper_test

import (
	gocontext "context"
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

func (suite *KeeperTestSuite) TestGRPCInvites() {

	ctx, _, _, k := SetupTestInput()
	queryClient := suite.queryClient
	k.Invite(ctx, testInviteSender, testUser)
	var inviteRes types.Invites
	_ = inviteRes

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
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}
