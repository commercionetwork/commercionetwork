package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"

	"github.com/stretchr/testify/require"
)

// TODO: fail test
func TestDefaultGenesisState(t *testing.T) {
	expted := types.GenesisState{}
	require.Equal(t, expted, *types.DefaultGenesis())
}

func TestInitGenesis(t *testing.T) {
	defGen := types.DefaultGenesis()
	ctx, _, _, k := SetupTestInput()
	require.Equal(t, &types.GenesisState{LiquidityPoolAmount: sdk.Coins(nil), Invites: []*types.Invite(nil), TrustedServiceProviders: nil, Memberships: []*types.Membership(nil)}, defGen)
	k.InitGenesis(ctx, *defGen)
	export := k.ExportGenesis(ctx)
	require.Equal(t, &types.GenesisState{LiquidityPoolAmount: sdk.Coins(nil), Invites: []*types.Invite{}, TrustedServiceProviders: nil, Memberships: []*types.Membership{}}, export)

	var tsps []string
	tsps = append(tsps, "cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")
	tsps = append(tsps, "cosmos1h7tw92a66gr58pxgmf6cc336lgxadpjz5d5psf")
	tsps = append(tsps, "cosmos14lultfckehtszvzw4ehu0apvsr77afvyhgqhwh")

	var invites []*types.Invite

	invites = append(invites, &types.Invite{Sender: "--", SenderMembership: "bronze", User: "---", Status: 1})
	var memberships []*types.Membership
	now := time.Now()
	now = now.Add(time.Hour * 24 * 7)
	memberships = append(memberships, &types.Membership{Owner: "cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0",
		TspAddress:     "cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae",
		MembershipType: "black", ExpiryAt: &now})

	genStateWithData := types.GenesisState{
		LiquidityPoolAmount:     sdk.Coins(nil),
		Invites:                 invites,
		Memberships:             memberships,
		TrustedServiceProviders: tsps,
	}
	k.InitGenesis(ctx, genStateWithData)

	export = k.ExportGenesis(ctx)

	require.Equal(t, genStateWithData.Invites, export.Invites)
	//require.Equal(t, genStateWithData.Memberships, export.Memberships) // TODO fix expiryAt
	require.Equal(t, genStateWithData.TrustedServiceProviders, export.TrustedServiceProviders)

	//require.Equal(t, *export, genStateWithData)

}

/*func TestExportGenesis(t *testing.T) {

}*/

func TestValidateGenesis(t *testing.T) {

}
