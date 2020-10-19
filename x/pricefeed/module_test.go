package pricefeed_test

import (
	"encoding/json"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	government "github.com/commercionetwork/commercionetwork/x/government/keeper"
	"github.com/commercionetwork/commercionetwork/x/pricefeed"
	pricefeedKeeper "github.com/commercionetwork/commercionetwork/x/pricefeed/keeper"
)

func TestAppModuleBasic(t *testing.T) {
	amb := pricefeed.AppModuleBasic{}
	require.Equal(t, "pricefeed", amb.Name())
	require.Equal(t, `{"oracles":null,"assets":null,"raw_prices":null,"denom_blacklist":["uccc"]}`, string(amb.DefaultGenesis()))

	require.NoError(t, amb.ValidateGenesis(json.RawMessage(`{}`)))
	require.NoError(t, amb.ValidateGenesis(json.RawMessage(`{"pool":[{"denom":"coin","amount":"10"}]}`)))
	require.Error(t, amb.ValidateGenesis(json.RawMessage(``)))

	require.NotNil(t, amb.GetTxCmd(nil))
	require.NotNil(t, amb.GetQueryCmd(nil))
	require.NotPanics(t, func() { amb.RegisterCodec(codec.New()) })
}

func TestNewAppModule(t *testing.T) {
	am := pricefeed.NewAppModule(pricefeedKeeper.Keeper{}, government.Keeper{})

	require.Equal(t, "pricefeed", am.QuerierRoute())
	require.Equal(t, "pricefeed", am.Route())
	require.NotNil(t, am.NewQuerierHandler())

	handler := am.NewHandler()
	require.NotNil(t, handler)
	_, err := handler(sdk.Context{}, &sdk.TestMsg{})
	require.Error(t, err, "pricefeed doesn't handle messages")

	require.NotPanics(t, func() { am.BeginBlock(sdk.Context{}, abci.RequestBeginBlock{}) })
	require.Panics(t, func() { am.EndBlock(sdk.Context{}, abci.RequestEndBlock{}) })
	require.NotPanics(t, func() { am.RegisterRESTRoutes(context.CLIContext{}, mux.NewRouter()) })
}

func TestAppModule_InitGenesis(t *testing.T) {
	_, ctx, govk, k := pricefeedKeeper.SetupTestInput()
	am := pricefeed.NewAppModule(k, govk)
	require.Equal(t, 0, len(am.InitGenesis(ctx, json.RawMessage(`{}`))))
	require.Panics(t, func() { am.InitGenesis(sdk.Context{}, json.RawMessage(`{}`)) })
}

func TestAppModule_ExportGenesis(t *testing.T) {
	_, ctx, govk, k := pricefeedKeeper.SetupTestInput()
	require.Equal(t, `{"oracles":null,"assets":null,"raw_prices":[],"denom_blacklist":null}`,
		string(pricefeed.NewAppModule(k, govk).ExportGenesis(ctx)))
}
