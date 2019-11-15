package v1_2_1

import "github.com/cosmos/cosmos-sdk/x/genutil"

// Migrate migrates exported state from v1.2.0 to a v1.2.1 genesis state.
func Migrate(appState genutil.AppMap) genutil.AppMap {
	return appState
}
