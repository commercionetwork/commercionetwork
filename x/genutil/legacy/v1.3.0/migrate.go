package v1_3_0

import (
	"github.com/cosmos/cosmos-sdk/codec"

	v120docs "github.com/commercionetwork/commercionetwork/x/documents/legacy/v1.2.0"
	v130docs "github.com/commercionetwork/commercionetwork/x/documents/legacy/v1.3.0"

	"github.com/cosmos/cosmos-sdk/x/genutil"
)

// Migrate migrates exported state from v1.2.0 to a v1.3.0 genesis state.
func Migrate(appState genutil.AppMap) genutil.AppMap {
	v120Codec := codec.New()
	codec.RegisterCrypto(v120Codec)

	v130Codec := codec.New()
	codec.RegisterCrypto(v130Codec)

	// Migrate documents state
	if appState[v120docs.ModuleName] != nil {
		var genDocs v120docs.GenesisState
		v120Codec.MustUnmarshalJSON(appState[v120docs.ModuleName], &genDocs)

		delete(appState, v120docs.ModuleName) // delete old key in case the name changed
		appState[v130docs.ModuleName] = v130Codec.MustMarshalJSON(
			v130docs.Migrate(genDocs),
		)
	}

	return appState
}
