package memberships_test

import (
	"encoding/json"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/auth"

	government "github.com/commercionetwork/commercionetwork/x/government/keeper"
	"github.com/commercionetwork/commercionetwork/x/memberships"
	"github.com/commercionetwork/commercionetwork/x/memberships/keeper"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
)

func TestAppModuleBasic(t *testing.T) {
	amb := memberships.AppModuleBasic{}
	require.Equal(t, "accreditations", amb.Name())
	require.Equal(t, `{"liquidity_pool_amount":[],"invites":null,"trusted_service_providers":null,"memberships":null}`, string(amb.DefaultGenesis()))

	require.NoError(t, amb.ValidateGenesis(json.RawMessage(`{}`)))
	require.Error(t, amb.ValidateGenesis(json.RawMessage(``)))

	require.NotNil(t, amb.GetTxCmd(nil))
	require.NotNil(t, amb.GetQueryCmd(nil))
	require.NotPanics(t, func() { amb.RegisterCodec(codec.New()) })
}

func TestAppModule(t *testing.T) {
	moduleName := "accreditations"
	am := memberships.NewAppModule(keeper.Keeper{}, supply.Keeper{}, government.Keeper{}, auth.AccountKeeper{})
	require.Equal(t, moduleName, am.Name())
	require.Equal(t, moduleName, am.QuerierRoute())
	require.Equal(t, moduleName, am.Route())

	handler := am.NewHandler()
	require.NotNil(t, handler)
	_, err := handler(sdk.Context{}, &sdk.TestMsg{})
	require.Error(t, err)

	require.NotNil(t, am.NewQuerierHandler())
}

func TestAppModule_InitGenesis(t *testing.T) {
	ctx, _, gk, k := keeper.SetupTestInput()
	am := memberships.NewAppModule(k, k.SupplyKeeper, gk, auth.AccountKeeper{})
	require.Equal(t, 0, len(am.InitGenesis(ctx, json.RawMessage(`{}`))))
	require.Panics(t, func() { am.InitGenesis(sdk.Context{}, json.RawMessage(`{}`)) })
}

/*
func TestAppModule_ExportGenesis(t *testing.T) {
	ctx, bk, gk, k := SetupTestInput()
	require.Equal(t, `{"oracles":null,"assets":null,"raw_prices":[],"denom_blacklist":null}`,
		string(memberships.NewAppModule(k, bk.supply.keeper, gk, auth.AccountKeeper{}).ExportGenesis(ctx)))
}
*/
