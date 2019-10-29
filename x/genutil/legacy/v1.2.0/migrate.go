package v1_2_0

import (
	"github.com/cosmos/cosmos-sdk/codec"

	v110docs "github.com/commercionetwork/commercionetwork/x/docs/legacy/v1.1.0"
	v120docs "github.com/commercionetwork/commercionetwork/x/docs/legacy/v1.2.0"

	"github.com/cosmos/cosmos-sdk/x/genutil"
)

// Migrate migrates exported state from v1.1.0 to a v1.2.0 genesis state.
func Migrate(appState genutil.AppMap) genutil.AppMap {
	v110Codec := codec.New()
	codec.RegisterCrypto(v110Codec)

	v120Codec := codec.New()
	codec.RegisterCrypto(v120Codec)

	// Migrate docs state
	if appState[v110docs.ModuleName] != nil {
		var genDocs v110docs.GenesisState
		v110Codec.MustUnmarshalJSON(appState[v110docs.ModuleName], &genDocs)

		delete(appState, v110docs.ModuleName) // delete old key in case the name changed
		appState[v120docs.ModuleName] = v120Codec.MustMarshalJSON(
			v120docs.Migrate(genDocs),
		)
	}

	return appState
}
