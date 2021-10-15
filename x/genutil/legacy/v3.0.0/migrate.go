package v3_0_0

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	sdkLegacy "github.com/cosmos/cosmos-sdk/x/genutil/legacy/v040"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibctransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	ibc "github.com/cosmos/cosmos-sdk/x/ibc/core/types"

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

	v220vbr "github.com/commercionetwork/commercionetwork/x/vbr/legacy/v2.2.0"
	v300vbr "github.com/commercionetwork/commercionetwork/x/vbr/legacy/v3.0.0"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
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
		v039Codec.MustUnmarshalJSON(appState[v220did.ModuleName], &didGenState)
		appState[v300did.ModuleName] = v040Codec.MustMarshalJSON(v300did.Migrate(didGenState))
		delete(appState, v220did.ModuleName)
	}

	if appState[v220docs.ModuleName] != nil {
		var docGenState v220docs.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v220docs.ModuleName], &docGenState)
		appState[v300docs.ModuleName] = v040Codec.MustMarshalJSON(v300docs.Migrate(docGenState))
	}
	//appState[v300docs.ModuleName] = appState[v220docs.ModuleName]

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

	if appState[v220vbr.ModuleName] != nil {
		var vbrGenState v220vbr.GenesisState
		v039Codec.MustUnmarshalJSON(appState[v220vbr.ModuleName], &vbrGenState)
		appState[v300vbr.ModuleName] = v040Codec.MustMarshalJSON(v300vbr.Migrate(vbrGenState))
	}

	//appState[wasm.ModuleName] = wasmKeeper.InitGenesis()
	wasmModule := &wasmTypes.GenesisState{}

	wasmModule.Params.InstantiateDefaultPermission = 3
	wasmModule.Params.CodeUploadAccess.Permission = 3
	wasmModule.Params.MaxWasmCodeSize = 1228800
	appState[wasm.ModuleName] = v040Codec.MustMarshalJSON(wasmModule)
	appState[ibctransfertypes.ModuleName] = v040Codec.MustMarshalJSON(ibctransfertypes.DefaultGenesisState())

	appState["ibc"] = v040Codec.MustMarshalJSON(ibc.DefaultGenesisState())
	appState[capabilitytypes.ModuleName] = v040Codec.MustMarshalJSON(capabilitytypes.DefaultGenesis())
	appState[evidencetypes.ModuleName] = v040Codec.MustMarshalJSON(evidencetypes.DefaultGenesisState())
	appState[evidencetypes.ModuleName] = v040Codec.MustMarshalJSON(evidencetypes.DefaultGenesisState())
	appState[govtypes.ModuleName] = v040Codec.MustMarshalJSON(govtypes.DefaultGenesisState())

	return appState

}
