package commerciomint_test

import (
	"encoding/json"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/commerciomint"
	"github.com/commercionetwork/commercionetwork/x/commerciomint/keeper"
)

func TestAppModuleBasic(t *testing.T) {
	amb := commerciomint.AppModuleBasic{}
	require.Equal(t, "commerciomint", amb.Name())
	require.Equal(t, `{"positions":[],"pool_amount":[],"collateral_rate":"2"}`, string(amb.DefaultGenesis()))

	require.Panics(t, func() { amb.ValidateGenesis(json.RawMessage(`{}`)) })
	require.NoError(t, amb.ValidateGenesis(json.RawMessage(`{"positions":[],"pool_amount":[],"collateral_rate":"2"}`)))
	require.Error(t, amb.ValidateGenesis(json.RawMessage(``)))

	require.NotNil(t, amb.GetTxCmd(codec.New()))
	require.NotNil(t, amb.GetQueryCmd(codec.New()))
	require.NotPanics(t, func() { amb.RegisterCodec(codec.New()) })
}

func TestAppModule(t *testing.T) {
	am := commerciomint.NewAppModule(keeper.Keeper{}, supply.Keeper{})
	require.Equal(t, "commerciomint", am.Name())
	require.Equal(t, "commerciomint", am.QuerierRoute())
	require.Equal(t, "commerciomint", am.Route())

	handler := am.NewHandler()
	require.NotNil(t, handler)
	_, err := handler(sdk.Context{}, &sdk.TestMsg{})
	require.Error(t, err)

	require.NotNil(t, am.NewQuerierHandler())
}
