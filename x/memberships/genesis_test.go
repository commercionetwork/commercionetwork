package memberships_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships"

	"github.com/stretchr/testify/require"
)

func TestDefaultGenesisState(t *testing.T) {
	expted := memberships.GenesisState{}
	require.Equal(t, expted, memberships.DefaultGenesisState("uccc"))
}

func TestInitGenesis(t *testing.T) {
	/*defGen := memberships.DefaultGenesisState("uccc")
	ctx, _, _, k := SetupTestInput()
	require.Equal(t, memberships.GenesisState{LiquidityPoolAmount: sdk.Coins(nil), Invites: types.Invites{}, TrustedServiceProviders: ctypes.Addresses{}, Memberships: types.Memberships{}}, defGen)
	memberships.InitGenesis(ctx, k, k.SupplyKeeper, defGen)*/
	/*export := memberships.ExportGenesis(ctx, k)
	require.Equal(t, commerciomint.GenesisState{Positions: []types.Position{}, LiquidityPoolAmount: sdk.Coins(nil), CreditsDenom: "test", CollateralRate: sdk.NewDec(2)}, export)
	*/
}

func TestExportGenesis(t *testing.T) {

}

func TestValidateGenesis(t *testing.T) {

}
