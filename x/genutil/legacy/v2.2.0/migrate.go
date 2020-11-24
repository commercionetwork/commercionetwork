package v2_2_0

import (
	commKyc "github.com/commercionetwork/commercionetwork/x/commerciokyc"
	v212memberships "github.com/commercionetwork/commercionetwork/x/commerciokyc/legacy/v2.1.2"
	commKycTypes "github.com/commercionetwork/commercionetwork/x/commerciokyc/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	govern "github.com/commercionetwork/commercionetwork/x/government"
	governTypes "github.com/commercionetwork/commercionetwork/x/government/types"

	"github.com/cosmos/cosmos-sdk/x/genutil"
)

func Migrate(appState genutil.AppMap) genutil.AppMap {
	oldAccountsCdc := codec.New()
	codec.RegisterCrypto(oldAccountsCdc)

	v039Codec := codec.New()
	codec.RegisterCrypto(v039Codec)

	if appState[governTypes.ModuleName] == nil {
		panic("Government module not set: invalid genesis file")
	}

	var govState govern.GenesisState
	v039Codec.MustUnmarshalJSON(appState[governTypes.ModuleName], &govState)

	if govState.GovernmentAddress == nil {
		panic("Government address not set: invalid genesis file")
	}

	govAddr := govState.GovernmentAddress
	expiredAt := int64(4733640)

	if appState[v212memberships.ModuleName] != nil {
		var genMemberships v212memberships.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v212memberships.ModuleName], &genMemberships)

		delete(appState, v212memberships.ModuleName) // delete old key in case the name changed
		appState[commKycTypes.ModuleName] = v039Codec.MustMarshalJSON(
			migrateMemberships(genMemberships, govAddr, expiredAt),
		)
	}

	return appState
}

func migrateMemberships(oldState v212memberships.GenesisState, govAddress sdk.AccAddress, expiredAt int64) commKyc.GenesisState {

	memberships := commKycTypes.Memberships{}

	for _, oldMembership := range oldState.Memberships {
		membership := migrateMembership(oldMembership, govAddress, expiredAt)
		memberships = append(memberships, membership)
	}

	return commKyc.GenesisState{
		LiquidityPoolAmount:     oldState.LiquidityPoolAmount,
		TrustedServiceProviders: oldState.TrustedServiceProviders,
		Memberships:             memberships,
	}
}

func migrateMembership(oldMembership v212memberships.Membership, govAddress sdk.AccAddress, expiredAt int64) commKycTypes.Membership {
	return commKycTypes.Membership{
		TspAddress:     govAddress,
		ExpiryAt:       expiredAt,
		Owner:          oldMembership.Owner,
		MembershipType: oldMembership.MembershipType,
	}
}
