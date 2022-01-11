package epochs_test

import (
	"testing"
	"time"

	simapp "github.com/commercionetwork/commercionetwork/testutil/simapp"
	"github.com/commercionetwork/commercionetwork/x/epochs"
	"github.com/commercionetwork/commercionetwork/x/epochs/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestEpochsExportGenesis(t *testing.T) {
	//app := simapp.Setup(false)
	app := simapp.New("")
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	chainStartTime := ctx.BlockTime()
	genState := types.DefaultGenesis()
	epochs.InitGenesis(ctx, app.EpochsKeeper, *genState)
	genesis := epochs.ExportGenesis(ctx, app.EpochsKeeper)
	require.Len(t, genesis.Epochs, 5)

	require.Equal(t, genesis.Epochs[0].Identifier, "day")
	require.Equal(t, genesis.Epochs[0].StartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[0].Duration, time.Hour*24)
	require.Equal(t, genesis.Epochs[0].CurrentEpoch, int64(0))
	require.Equal(t, genesis.Epochs[0].CurrentEpochStartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[0].EpochCountingStarted, false)
	require.Equal(t, genesis.Epochs[0].CurrentEpochEnded, true)
	require.Equal(t, genesis.Epochs[4].Identifier, "week")
	require.Equal(t, genesis.Epochs[4].StartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[4].Duration, time.Hour*24*7)
	require.Equal(t, genesis.Epochs[4].CurrentEpoch, int64(0))
	require.Equal(t, genesis.Epochs[4].CurrentEpochStartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[4].EpochCountingStarted, false)
	require.Equal(t, genesis.Epochs[4].CurrentEpochEnded, true)
}

func TestEpochsInitGenesis(t *testing.T) {
	//app := simapp.Setup(false)

	app := simapp.New("")
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	// On init genesis, default epochs information is set
	// To check init genesis again, should make it fresh status
	epochInfos := app.EpochsKeeper.AllEpochInfos(ctx)
	for _, epochInfo := range epochInfos {
		app.EpochsKeeper.DeleteEpochInfo(ctx, epochInfo.Identifier)
	}

	now := time.Now()
	ctx = ctx.WithBlockHeight(1)
	ctx = ctx.WithBlockTime(now)

	epochs.InitGenesis(ctx, app.EpochsKeeper, types.GenesisState{
		Epochs: []types.EpochInfo{
			{
				Identifier:            "monthly",
				StartTime:             time.Time{},
				Duration:              time.Hour * 24,
				CurrentEpoch:          0,
				CurrentEpochStartTime: time.Time{},
				EpochCountingStarted:  true,
				CurrentEpochEnded:     true,
			},
		},
	})

	epochInfo := app.EpochsKeeper.GetEpochInfo(ctx, "monthly")
	require.Equal(t, epochInfo.Identifier, "monthly")
	require.Equal(t, epochInfo.StartTime.UTC().String(), now.UTC().String())
	require.Equal(t, epochInfo.Duration, time.Hour*24)
	require.Equal(t, epochInfo.CurrentEpoch, int64(0))
	require.Equal(t, epochInfo.CurrentEpochStartTime.UTC().String(), ctx.BlockTime().UTC().String())
	require.Equal(t, epochInfo.EpochCountingStarted, true)
	require.Equal(t, epochInfo.CurrentEpochEnded, true)
}
