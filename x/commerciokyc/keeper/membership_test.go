package keeper

import (
	"fmt"
	"testing"
	"time"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

func TestKeeper_AssignMembership(t *testing.T) {
	tests := []struct {
		name               string
		existingMembership string
		membershipType     string
		expectedNotExists  bool
		userIsTsp          bool
		user               sdk.AccAddress
		tsp                sdk.AccAddress
		expiredAt          time.Time
		error              error
	}{
		{
			name:           "Invalid membership type returns error",
			membershipType: "grn",
			user:           testUser,
			tsp:            testTsp,
			expiredAt:      testExpiration,
			error:          sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid membership type: grn"),
		},
		{
			name:           "Membership with invalid expired date",
			user:           testUser,
			tsp:            testTsp,
			expiredAt:      testExpirationNegative,
			membershipType: types.MembershipTypeBronze,
			error:          sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid expiry date: %s", testExpirationNegative)),
		},
		/*{
			name:               "Invalid tsp",
			user:               testUser,
			tsp:                testUser,
			height:             testHeight,
			existingMembership: types.MembershipTypeBronze,
			error:              sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid expiry height: -1"),
		},*/
		{
			name:           "Non existing membership is properly saved",
			user:           testUser,
			tsp:            testTsp,
			expiredAt:      testExpiration,
			membershipType: types.MembershipTypeBronze,
		},
		{
			name:               "Existing membership is replaced",
			user:               testUser,
			tsp:                testTsp,
			expiredAt:          testExpiration,
			existingMembership: types.MembershipTypeBronze,
			membershipType:     types.MembershipTypeGold,
		},
		{
			name:               "Cannot assign tsp membership",
			userIsTsp:          true,
			user:               testUser,
			tsp:                testTsp,
			expiredAt:          testExpiration,
			existingMembership: types.MembershipTypeBlack,
			membershipType:     types.MembershipTypeGold,
			error:              sdkErr.Wrap(sdkErr.ErrUnauthorized, "account \""+testUser.String()+"\" is a Trust Service Provider: remove from tsps list before"),
		},
		/*{
			name:               "Assign \"none\" membership type to delete membership",
			user:               testUser,
			tsp:                testTsp,
			height:             testHeight,
			existingMembership: types.MembershipTypeBronze,
			membershipType:     types.MembershipTypeNone,
			expectedNotExists:  true,
			error:              sdkErr.Wrap(sdkErr.ErrUnknownRequest, "membership not found for user \"cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae\""),
		},*/
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			if len(test.existingMembership) != 0 {
				err := k.AssignMembership(ctx, test.user, test.existingMembership, test.tsp, test.expiredAt)
				require.NoError(t, err)
			}

			if test.userIsTsp {
				k.AddTrustedServiceProvider(ctx, test.user)
			}

			err := k.AssignMembership(ctx, test.user, test.membershipType, test.tsp, test.expiredAt)

			if err != nil {
				if test.expectedNotExists {
					_, err2 := k.GetMembership(ctx, test.user)
					require.Equal(t, test.error.Error(), err2.Error())
				}
				require.Equal(t, test.error.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
