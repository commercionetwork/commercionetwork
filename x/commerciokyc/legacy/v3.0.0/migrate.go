package v3_0_0

import (
	"time"

	v220commerciokyc "github.com/commercionetwork/commercionetwork/x/commerciokyc/legacy/v2.2.0"
	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

// Migrate accepts exported genesis state from v2.2.0 and migrates it to v3.0.0
func Migrate(oldGenState v220commerciokyc.GenesisState) *types.GenesisState {

	var memberships []*types.Membership
	for _, oldMembership := range oldGenState.Memberships {
		var membership types.Membership
		var expiryAt time.Time
		expiryAt = oldMembership.ExpiryAt

		membership.Owner = oldMembership.Owner.String()
		membership.TspAddress = oldMembership.TspAddress.String()
		membership.MembershipType = oldMembership.MembershipType
		membership.ExpiryAt = &expiryAt

		memberships = append(memberships, &membership)
	}

	var invites []*types.Invite
	for _, oldInvite := range oldGenState.Invites {
		var invite types.Invite

		invite.Sender = oldInvite.Sender.String()
		invite.SenderMembership = oldInvite.SenderMembership
		invite.Status = uint64(oldInvite.Status)
		invite.User = oldInvite.User.String()

		invites = append(invites, &invite)
	}

	var tsps []string
	for _, oldTsp := range oldGenState.TrustedServiceProviders {
		tsps = append(tsps, oldTsp.String())
	}

	return &types.GenesisState{
		LiquidityPoolAmount:     oldGenState.LiquidityPoolAmount,
		TrustedServiceProviders: tsps,
		Invites:                 invites,
		Memberships:             memberships,
	}

}
