package v3_0_0

import (
	"time"

	v220commerciokyc "github.com/commercionetwork/commercionetwork/x/commerciokyc/legacy/v2.2.0"
	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

// Migrate accepts exported genesis state from v2.2.0 and migrates it to v3.0.0
func Migrate(v220GenState v220commerciokyc.GenesisState) *types.GenesisState {

	var memberships []*types.Membership
	for _, v220Membership := range v220GenState.Memberships {
		var expiryAt time.Time
		expiryAt = v220Membership.ExpiryAt
		memberships = append(memberships, &types.Membership{
			Owner:          v220Membership.Owner.String(),
			TspAddress:     v220Membership.TspAddress.String(),
			MembershipType: v220Membership.MembershipType,
			ExpiryAt:       &expiryAt,
		})
	}

	var invites []*types.Invite
	for _, v220Invite := range v220GenState.Invites {
		invites = append(invites, &types.Invite{
			Sender:           v220Invite.Sender.String(),
			SenderMembership: v220Invite.SenderMembership,
			Status:           uint64(v220Invite.Status),
			User:             v220Invite.User.String(),
		})
	}

	var tsps []string
	for _, oldTsp := range v220GenState.TrustedServiceProviders {
		tsps = append(tsps, oldTsp.String())
	}

	return &types.GenesisState{
		LiquidityPoolAmount:     v220GenState.LiquidityPoolAmount,
		TrustedServiceProviders: tsps,
		Invites:                 invites,
		Memberships:             memberships,
	}

}
