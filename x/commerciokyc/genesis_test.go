package commerciokyc_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	commerciokyc "github.com/commercionetwork/commercionetwork/x/commerciokyc"
	"github.com/commercionetwork/commercionetwork/x/commerciokyc/keeper"
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
	ctx, _, _, k := keeper.SetupTestInput()
	require.Equal(t, &types.GenesisState{LiquidityPoolAmount: sdk.Coins(nil), Invites: []*types.Invite(nil), TrustedServiceProviders: nil, Memberships: []*types.Membership(nil)}, defGen)
	commerciokyc.InitGenesis(ctx, k, *defGen)
	export := commerciokyc.ExportGenesis(ctx, k)
	require.Equal(t, &types.GenesisState{LiquidityPoolAmount: sdk.Coins(nil), Invites: []*types.Invite{}, TrustedServiceProviders: nil, Memberships: []*types.Membership{}}, export)
}

/*func TestExportGenesis(t *testing.T) {

}

func TestValidateGenesis(t *testing.T) {

}*/
