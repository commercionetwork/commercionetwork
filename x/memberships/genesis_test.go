package memberships_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/memberships"
	"github.com/commercionetwork/commercionetwork/x/memberships/keeper"
	"github.com/commercionetwork/commercionetwork/x/memberships/types"

	"github.com/stretchr/testify/require"
)

func TestDefaultGenesisState(t *testing.T) {
	expted := memberships.GenesisState{}
	require.Equal(t, expted, memberships.DefaultGenesisState("uccc"))
}

func TestInitGenesis(t *testing.T) {
	defGen := memberships.DefaultGenesisState("uccc")
	ctx, _, _, k := keeper.SetupTestInput()
	require.Equal(t, memberships.GenesisState{LiquidityPoolAmount: sdk.Coins(nil), Invites: nil, TrustedServiceProviders: nil, Memberships: nil}, defGen)
	memberships.InitGenesis(ctx, k, k.SupplyKeeper, defGen)
	export := memberships.ExportGenesis(ctx, k)
	require.Equal(t, memberships.GenesisState{LiquidityPoolAmount: sdk.Coins(nil), Invites: types.Invites{}, TrustedServiceProviders: nil, Memberships: types.Memberships{}}, export)
}

/*func TestExportGenesis(t *testing.T) {

}

func TestValidateGenesis(t *testing.T) {

}*/
