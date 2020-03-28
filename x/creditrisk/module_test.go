package creditrisk_test

import (
	"encoding/json"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/commercionetwork/commercionetwork/x/creditrisk"
	"github.com/commercionetwork/commercionetwork/x/creditrisk/types"
)

func TestAppModuleBasic(t *testing.T) {
	amb := creditrisk.AppModuleBasic{}
	require.Equal(t, "creditrisk", amb.Name())
	require.Equal(t, "{\"pool\":[]}", string(amb.DefaultGenesis()))

	require.NoError(t, amb.ValidateGenesis(json.RawMessage(`{}`)))
	require.NoError(t, amb.ValidateGenesis(json.RawMessage(`{"pool":[{"denom":"coin","amount":"10"}]}`)))
	require.Error(t, amb.ValidateGenesis(json.RawMessage(``)))

	require.Nil(t, amb.GetTxCmd(nil))
	require.NotNil(t, amb.GetQueryCmd(codec.New()))
	require.NotPanics(t, func() { amb.RegisterCodec(codec.New()) })
}

func TestAppModule(t *testing.T) {
	am := creditrisk.NewAppModule(creditrisk.Keeper{})

	require.Equal(t, "creditrisk", am.QuerierRoute())
	require.Equal(t, "creditrisk", am.Route())
	require.NotNil(t, am.NewQuerierHandler())

	handler := am.NewHandler()
	require.NotNil(t, handler)
	_, err := handler(sdk.Context{}, &sdk.TestMsg{})
	require.Error(t, err, "creditrisk doesn't handle messages")

	require.NotPanics(t, func() { am.BeginBlock(sdk.Context{}, abci.RequestBeginBlock{}) })
	require.Equal(t, 0, len(am.EndBlock(sdk.Context{}, abci.RequestEndBlock{})))
	require.NotPanics(t, func() { am.RegisterRESTRoutes(context.CLIContext{}, mux.NewRouter()) })
}

func TestAppModule_InitGenesis(t *testing.T) {
	ctx, _, k := SetupTestInput()
	am := creditrisk.NewAppModule(k)
	require.Equal(t, 0, len(am.InitGenesis(ctx, json.RawMessage(`{}`))))
	require.Panics(t, func() { am.InitGenesis(sdk.Context{}, json.RawMessage(`{}`)) })
}

func TestAppModule_ExportGenesis(t *testing.T) {
	ctx, sk, k := SetupTestInput()
	am := creditrisk.NewAppModule(k)

	require.Equal(t, `{"pool":[]}`, string(am.ExportGenesis(ctx)))

	modAcc := sk.GetModuleAccount(ctx, types.ModuleName)
	newcoins := modAcc.GetCoins().Add(sdk.NewInt64Coin("coin", 10))
	modAcc.SetCoins(newcoins)
	sk.SetModuleAccount(ctx, modAcc)

	require.Equal(t, `{"pool":[{"denom":"coin","amount":"10"}]}`, string(am.ExportGenesis(ctx)))
}
