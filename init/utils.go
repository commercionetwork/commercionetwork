package init

import (
	app "commercio-network"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/tendermint/tendermint/libs/common"
)

func initializeEmptyGenesis(cdc *codec.Codec, genFile string, overwrite bool) (appState json.RawMessage, err error) {
	if !overwrite && common.FileExists(genFile) {
		return nil, fmt.Errorf("genesis.json file already exists: %v", genFile)
	}

	return codec.MarshalJSONIndent(cdc, app.NewDefaultGenesisState())
}
