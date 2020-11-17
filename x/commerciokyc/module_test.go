package commerciokyc_test

import (
	"encoding/json"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/auth"

	commerciokyc "github.com/commercionetwork/commercionetwork/x/commerciokyc"
	"github.com/commercionetwork/commercionetwork/x/commerciokyc/keeper"
	government "github.com/commercionetwork/commercionetwork/x/government/keeper"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
)

func TestAppModuleBasic(t *testing.T) {
	amb := commerciokyc.AppModuleBasic{}
	require.Equal(t, "commerciokyc", amb.Name())
	require.Equal(t, `{"liquidity_pool_amount":[],"invites":null,"trusted_service_providers":null,"memberships":null}`, string(amb.DefaultGenesis()))

	require.NoError(t, amb.ValidateGenesis(json.RawMessage(`{}`)))
	require.Error(t, amb.ValidateGenesis(json.RawMessage(``)))

	require.NotNil(t, amb.GetTxCmd(nil))
	require.NotNil(t, amb.GetQueryCmd(nil))
	require.NotPanics(t, func() { amb.RegisterCodec(codec.New()) })
}

func TestAppModule(t *testing.T) {
	moduleName := "commerciokyc"
	am := commerciokyc.NewAppModule(keeper.Keeper{}, supply.Keeper{}, government.Keeper{}, auth.AccountKeeper{})
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
	am := commerciokyc.NewAppModule(k, k.SupplyKeeper, gk, auth.AccountKeeper{})
	require.Equal(t, 0, len(am.InitGenesis(ctx, json.RawMessage(`{}`))))
	require.Panics(t, func() { am.InitGenesis(sdk.Context{}, json.RawMessage(`{}`)) })
}

/*
func TestAppModule_ExportGenesis(t *testing.T) {
	ctx, bk, gk, k := SetupTestInput()
	require.Equal(t, `{"oracles":null,"assets":null,"raw_prices":[],"denom_blacklist":null}`,
		string(commerciokyc.NewAppModule(k, bk.supply.keeper, gk, auth.AccountKeeper{}).ExportGenesis(ctx)))
}
*/
