package v3_0_0

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"

	sdkLegacy "github.com/cosmos/cosmos-sdk/x/genutil/legacy/v040"

	v220government "github.com/commercionetwork/commercionetwork/x/government/legacy/v2.2.0"
	v300government "github.com/commercionetwork/commercionetwork/x/government/legacy/v3.0.0"

	v220docs "github.com/commercionetwork/commercionetwork/x/documents/legacy/v2.2.0"
	v300docs "github.com/commercionetwork/commercionetwork/x/documents/legacy/v3.0.0"

	v220did "github.com/commercionetwork/commercionetwork/x/did/legacy/v2.2.0"
	v300did "github.com/commercionetwork/commercionetwork/x/did/legacy/v3.0.0"

	v220commerciomint "github.com/commercionetwork/commercionetwork/x/commerciomint/legacy/v2.2.0"
	v300commerciomint "github.com/commercionetwork/commercionetwork/x/commerciomint/legacy/v3.0.0"

	v220commerciokyc "github.com/commercionetwork/commercionetwork/x/commerciokyc/legacy/v2.2.0"
	v300commerciokyc "github.com/commercionetwork/commercionetwork/x/commerciokyc/legacy/v3.0.0"

	"github.com/CosmWasm/wasmd/x/wasm"
)

func Migrate(appState types.AppMap, clientCtx client.Context) types.AppMap {
	v039Codec := codec.NewLegacyAmino()
	v040Codec := clientCtx.JSONMarshaler

	appState = sdkLegacy.Migrate(appState, clientCtx)

	if appState[v220government.ModuleName] != nil {
		var govGenState v220government.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v220government.ModuleName], &govGenState)
		appState[v300government.ModuleName] = v040Codec.MustMarshalJSON(v300government.Migrate(govGenState))
	}

	if appState[v220did.ModuleName] != nil {
		var didGenState v220did.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v220docs.ModuleName], &didGenState)
		appState[v300did.ModuleName] = v040Codec.MustMarshalJSON(v300did.Migrate(didGenState))
		delete(appState, v220did.ModuleName)
	}

	/*if appState[v220docs.ModuleName] != nil {
		var docGenState v220docs.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v220docs.ModuleName], &docGenState)
		appState[v300docs.ModuleName] = v040Codec.MustMarshalJSON(v300docs.Migrate(docGenState))
	}*/
	appState[v300docs.ModuleName] = appState[v220docs.ModuleName]

	if appState[v220commerciomint.ModuleName] != nil {
		var commerciomintGenState v220commerciomint.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v220commerciomint.ModuleName], &commerciomintGenState)
		appState[v300commerciomint.ModuleName] = v040Codec.MustMarshalJSON(v300commerciomint.Migrate(commerciomintGenState))

	}

	if appState[v220commerciokyc.ModuleName] != nil {
		var commerciokycGenState v220commerciokyc.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v220commerciokyc.ModuleName], &commerciokycGenState)
		appState[v300commerciokyc.ModuleName] = v040Codec.MustMarshalJSON(v300commerciokyc.Migrate(commerciokycGenState))

	}

	//appState[wasm.ModuleName] = wasmKeeper.InitGenesis()
	appState[wasm.ModuleName] = nil

	return appState

}
