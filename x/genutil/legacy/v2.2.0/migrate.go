package v2_2_0

import (
	"time"

	commKyc "github.com/commercionetwork/commercionetwork/x/commerciokyc"
	v212memberships "github.com/commercionetwork/commercionetwork/x/commerciokyc/legacy/v2.1.2"
	commKycTypes "github.com/commercionetwork/commercionetwork/x/commerciokyc/types"

	"github.com/commercionetwork/commercionetwork/x/upgrade"
	upgradeTypes "github.com/commercionetwork/commercionetwork/x/upgrade/types"
	v212vbr "github.com/commercionetwork/commercionetwork/x/vbr/legacy/v2.1.2"
	v220vbr "github.com/commercionetwork/commercionetwork/x/vbr/legacy/v2.2.0"
	vbrTypes "github.com/commercionetwork/commercionetwork/x/vbr/types"

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

	currentTime := time.Now()
	// Use only day and add 1 year so we have YYYY-MM-DD 00:00:00 and date should be the same in all nodes
	expiredAt, _ := time.Parse("2006-01-02", currentTime.Format("2006-01-02"))
	expiredAt = expiredAt.Add(time.Hour * 24 * 365)

	if appState[v212memberships.ModuleName] != nil {
		var genMemberships v212memberships.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v212memberships.ModuleName], &genMemberships)

		delete(appState, v212memberships.ModuleName) // delete old key in case the name changed
		appState[commKycTypes.ModuleName] = v039Codec.MustMarshalJSON(
			migrateMemberships(genMemberships, govAddr, expiredAt),
		)
	}

	// Migrate vbr state
	if appState[v212vbr.ModuleName] != nil {
		var genVbr v212vbr.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v212vbr.ModuleName], &genVbr)

		delete(appState, v212vbr.ModuleName) // delete old key in case the name changed
		appState[vbrTypes.ModuleName] = v039Codec.MustMarshalJSON(
			v220vbr.Migrate(genVbr),
		)
	}

	// Migrate upgrade module
	if appState[upgradeTypes.ModuleName] == nil {
		genUpgrade := upgrade.GenesisState{}
		appState[upgradeTypes.ModuleName] = v039Codec.MustMarshalJSON(
			genUpgrade,
		)
	}

	// Remove creditrisk pool
	// TODO: what do we do with the tokens already in the pool?
	if appState[creditriskModuleName] != nil {
		delete(appState, creditriskModuleName)
	}

	appState = commercioMintMigrate(appState, govAddr)

	return appState
}

func migrateMemberships(oldState v212memberships.GenesisState, govAddress sdk.AccAddress, expiryAt time.Time) commKyc.GenesisState {

	memberships := commKycTypes.Memberships{}

	for _, oldMembership := range oldState.Memberships {
		membership := migrateMembership(oldMembership, govAddress, expiryAt)
		memberships = append(memberships, membership)
	}

	return commKyc.GenesisState{
		LiquidityPoolAmount:     oldState.LiquidityPoolAmount,
		TrustedServiceProviders: oldState.TrustedServiceProviders,
		Invites:                 oldState.Invites,
		Memberships:             memberships,
	}
}

func migrateMembership(oldMembership v212memberships.Membership, govAddress sdk.AccAddress, expiryAt time.Time) commKycTypes.Membership {
	return commKycTypes.Membership{
		TspAddress:     govAddress,
		ExpiryAt:       expiryAt,
		Owner:          oldMembership.Owner,
		MembershipType: oldMembership.MembershipType,
	}
}
