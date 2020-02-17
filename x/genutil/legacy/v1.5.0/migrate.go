package v1_5_0

import (
	"github.com/commercionetwork/commercionetwork/x/memberships"
	v134memberships "github.com/commercionetwork/commercionetwork/x/memberships/legacy/v1.3.4"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/genutil"
)

// Migrate migrates exported state from v1.3.3 to a v1.3.4 genesis state
func Migrate(appState genutil.AppMap) genutil.AppMap {
	v134Codec := codec.New()
	codec.RegisterCrypto(v134Codec)

	v150Codec := codec.New()
	codec.RegisterCrypto(v150Codec)

	// Migrate vbr state
	if appState[v134memberships.ModuleName] != nil {
		var genMemberships v134memberships.GenesisState
		v134Codec.MustUnmarshalJSON(appState[v134memberships.ModuleName], &genMemberships)

		delete(appState, v134memberships.ModuleName) //delete old key in case the name changed
		appState[v134memberships.ModuleName] = v150Codec.MustMarshalJSON(
			migrateMemberships(genMemberships),
		)
	}

	return appState
}

func migrateMemberships(oldState v134memberships.GenesisState) memberships.GenesisState {
	ng := memberships.GenesisState{
		LiquidityPoolAmount:     oldState.LiquidityPoolAmount,
		TrustedServiceProviders: oldState.TrustedServiceProviders,
		Credentials:             oldState.Credentials,
		StableCreditsDenom:      oldState.StableCreditsDenom,
		Memberships:             oldState.Memberships,
	}

	mutateStatus := func(status bool) memberships.InviteStatus {
		if status {
			return memberships.InviteStatusRewarded
		}

		return memberships.InviteStatusPending
	}

	for _, invite := range oldState.Invites {
		ng.Invites = append(ng.Invites, memberships.Invite{
			Sender: invite.Sender,
			User:   invite.User,
			Status: mutateStatus(invite.Rewarded),
		})
	}

	return ng
}
