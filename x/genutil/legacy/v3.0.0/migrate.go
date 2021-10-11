package v3_0_0

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"

	v220government "github.com/commercionetwork/commercionetwork/x/government/legacy/v2.2.0"
	v300government "github.com/commercionetwork/commercionetwork/x/government/legacy/v3.0.0"

	v220docs "github.com/commercionetwork/commercionetwork/x/documents/legacy/v2.2.0"
	v300docs "github.com/commercionetwork/commercionetwork/x/documents/legacy/v3.0.0"

	v220commerciomint "github.com/commercionetwork/commercionetwork/x/commerciomint/legacy/v2.2.0"
	v300commerciomint "github.com/commercionetwork/commercionetwork/x/commerciomint/legacy/v3.0.0"
)

func Migrate(appState types.AppMap, clientCtx client.Context) types.AppMap {
	v039Codec := codec.NewLegacyAmino()
	v040Codec := clientCtx.JSONMarshaler

	if appState[v220government.ModuleName] != nil {
		var govGenState v220government.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v220government.ModuleName], &govGenState)
		appState[v300government.ModuleName] = v040Codec.MustMarshalJSON(v300government.Migrate(govGenState))

	}

	if appState[v220docs.ModuleName] != nil {
		var docGenState v220docs.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v220docs.ModuleName], &docGenState)
		appState[v300docs.ModuleName] = v040Codec.MustMarshalJSON(v300docs.Migrate(docGenState))

	}

	if appState[v220commerciomint.ModuleName] != nil {
		var commerciomintGenState v220commerciomint.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v220docs.ModuleName], &commerciomintGenState)
		appState[v300docs.ModuleName] = v040Codec.MustMarshalJSON(v300commerciomint.Migrate(commerciomintGenState))

	}

	return appState

}
