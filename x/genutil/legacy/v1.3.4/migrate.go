package v1_3_4

import (
	v133tbr "github.com/commercionetwork/commercionetwork/x/vbr/legacy/v1.3.3"
	v134vbr "github.com/commercionetwork/commercionetwork/x/vbr/legacy/v1.3.4"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/genutil"
)

// Migrate migrates exported state from v1.3.3 to a v1.3.4 genesis state
func Migrate(appState genutil.AppMap) genutil.AppMap {
	v133Codec := codec.New()
	codec.RegisterCrypto(v133Codec)

	v134Codec := codec.New()
	codec.RegisterCrypto(v134Codec)

	// Migrate vbr state
	if appState[v133tbr.ModuleName] != nil {
		var genTbr v133tbr.GenesisState
		v133Codec.MustUnmarshalJSON(appState[v133tbr.ModuleName], &genTbr)

		delete(appState, v133tbr.ModuleName) // delete old key in case the name changed
		appState[v134vbr.ModuleName] = v134Codec.MustMarshalJSON(
			v134vbr.Migrate(genTbr),
		)
	}

	return appState
}
