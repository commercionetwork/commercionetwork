package simapp

import (
	"time"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	sims "github.com/cosmos/cosmos-sdk/testutil/sims"

	"cosmossdk.io/log"
	tmdb "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"

	"github.com/commercionetwork/commercionetwork/app"
)

// New creates application instance with in-memory database and disabled logging.
func New(dir string) *app.App {
	db := tmdb.NewMemDB()
	logger := log.NewNopLogger()

	encoding := app.MakeEncodingConfig()

	appOpts := sims.EmptyAppOptions{}
	var wasmOpts []wasmkeeper.Option

	a := app.New(
		logger,
		db,
		nil,
		true,
		map[int64]bool{},
		dir,
		0,
		encoding,
		appOpts,
		app.GetEnabledProposals(),
		wasmOpts,
	)
	// InitChain updates deliverState which is required when app.NewContext is called
	a.InitChain(abci.RequestInitChain{
		ConsensusParams: defaultConsensusParams,
		AppStateBytes:   []byte("{}"),
		ChainId:         "commercionetwork",
	})
	return a
}

var defaultConsensusParams = &tmproto.ConsensusParams{
	Block: &tmproto.BlockParams{
		MaxBytes: 200000,
		MaxGas:   2000000,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}
