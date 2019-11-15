package v1_3_2

import "github.com/cosmos/cosmos-sdk/x/genutil"

// Migrate migrates exported state from v1.3.1 to a v1.3.2 genesis state.
func Migrate(appState genutil.AppMap) genutil.AppMap {
	return appState
}
