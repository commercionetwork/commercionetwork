package commerciokyc_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	commerciokyc "github.com/commercionetwork/commercionetwork/x/commerciokyc"
	"github.com/commercionetwork/commercionetwork/x/commerciokyc/keeper"
	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"

	"github.com/stretchr/testify/require"
)

func TestDefaultGenesisState(t *testing.T) {
	expted := commerciokyc.GenesisState{}
	require.Equal(t, expted, commerciokyc.DefaultGenesisState("uccc"))
}

func TestInitGenesis(t *testing.T) {
	defGen := commerciokyc.DefaultGenesisState("uccc")
	ctx, _, _, k := keeper.SetupTestInput()
	require.Equal(t, commerciokyc.GenesisState{LiquidityPoolAmount: sdk.Coins(nil), Invites: nil, TrustedServiceProviders: nil, Memberships: nil}, defGen)
	commerciokyc.InitGenesis(ctx, k, k.SupplyKeeper, defGen)
	export := commerciokyc.ExportGenesis(ctx, k)
	require.Equal(t, commerciokyc.GenesisState{LiquidityPoolAmount: sdk.Coins(nil), Invites: types.Invites{}, TrustedServiceProviders: nil, Memberships: types.Memberships{}}, export)
}

/*func TestExportGenesis(t *testing.T) {

}

func TestValidateGenesis(t *testing.T) {

}*/
