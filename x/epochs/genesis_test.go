package epochs_test

// import (
// 	"testing"
// 	"time"

// 	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
// 	simapp "github.com/commercionetwork/commercionetwork/testutil/simapp"
// 	"github.com/commercionetwork/commercionetwork/x/epochs"
// 	"github.com/commercionetwork/commercionetwork/x/epochs/types"
// 	"github.com/stretchr/testify/require"
// )

// func TestEpochsExportGenesis(t *testing.T) {
// 	app := simapp.New("")
// 	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

// 	chainStartTime := ctx.BlockTime()
// 	genState := types.DefaultGenesis()
// 	epochs.InitGenesis(ctx, app.EpochsKeeper, *genState)
// 	genesis := epochs.ExportGenesis(ctx, app.EpochsKeeper)
// 	require.Len(t, genesis.Epochs, 5)

// 	require.Equal(t, genesis.Epochs[0].Identifier, "day")
// 	require.Equal(t, genesis.Epochs[0].StartTime, chainStartTime)
// 	require.Equal(t, genesis.Epochs[0].Duration, time.Hour*24)
// 	require.Equal(t, genesis.Epochs[0].CurrentEpoch, int64(0))
// 	require.Equal(t, genesis.Epochs[0].CurrentEpochStartTime, chainStartTime)
// 	require.Equal(t, genesis.Epochs[0].EpochCountingStarted, false)
// 	require.Equal(t, genesis.Epochs[0].CurrentEpochEnded, true)

// 	require.Equal(t, genesis.Epochs[1].Identifier, "hour")
// 	require.Equal(t, genesis.Epochs[1].StartTime, chainStartTime)
// 	require.Equal(t, genesis.Epochs[1].Duration, time.Hour)
// 	require.Equal(t, genesis.Epochs[1].CurrentEpoch, int64(0))
// 	require.Equal(t, genesis.Epochs[1].CurrentEpochStartTime, chainStartTime)
// 	require.Equal(t, genesis.Epochs[1].EpochCountingStarted, false)
// 	require.Equal(t, genesis.Epochs[1].CurrentEpochEnded, true)

// 	require.Equal(t, genesis.Epochs[2].Identifier, "minute")
// 	require.Equal(t, genesis.Epochs[2].StartTime, chainStartTime)
// 	require.Equal(t, genesis.Epochs[2].Duration, time.Minute)
// 	require.Equal(t, genesis.Epochs[2].CurrentEpoch, int64(0))
// 	require.Equal(t, genesis.Epochs[2].CurrentEpochStartTime, chainStartTime)
// 	require.Equal(t, genesis.Epochs[2].EpochCountingStarted, false)
// 	require.Equal(t, genesis.Epochs[2].CurrentEpochEnded, true)

// 	require.Equal(t, genesis.Epochs[3].Identifier, "month")
// 	require.Equal(t, genesis.Epochs[3].StartTime, chainStartTime)
// 	require.Equal(t, genesis.Epochs[3].Duration, time.Hour*24*30)
// 	require.Equal(t, genesis.Epochs[3].CurrentEpoch, int64(0))
// 	require.Equal(t, genesis.Epochs[3].CurrentEpochStartTime, chainStartTime)
// 	require.Equal(t, genesis.Epochs[3].EpochCountingStarted, false)
// 	require.Equal(t, genesis.Epochs[3].CurrentEpochEnded, true)

// 	require.Equal(t, genesis.Epochs[4].Identifier, "week")
// 	require.Equal(t, genesis.Epochs[4].StartTime, chainStartTime)
// 	require.Equal(t, genesis.Epochs[4].Duration, time.Hour*24*7)
// 	require.Equal(t, genesis.Epochs[4].CurrentEpoch, int64(0))
// 	require.Equal(t, genesis.Epochs[4].CurrentEpochStartTime, chainStartTime)
// 	require.Equal(t, genesis.Epochs[4].EpochCountingStarted, false)
// 	require.Equal(t, genesis.Epochs[4].CurrentEpochEnded, true)
// }

// func TestEpochsInitGenesis(t *testing.T) {
// 	app := simapp.New("")
// 	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

// 	// On init genesis, default epochs information is set
// 	// To check init genesis again, should make it fresh status
// 	epochInfos := app.EpochsKeeper.AllEpochInfos(ctx)
// 	for _, epochInfo := range epochInfos {
// 		app.EpochsKeeper.DeleteEpochInfo(ctx, epochInfo.Identifier)
// 	}

// 	now := time.Now()
// 	ctx = ctx.WithBlockHeight(1)
// 	ctx = ctx.WithBlockTime(now)

// 	epochs.InitGenesis(ctx, app.EpochsKeeper, types.GenesisState{
// 		Epochs: []types.EpochInfo{
// 			{
// 				Identifier:            "month",
// 				StartTime:             time.Time{},
// 				Duration:              time.Hour * 24 * 30,
// 				CurrentEpoch:          0,
// 				CurrentEpochStartTime: time.Time{},
// 				EpochCountingStarted:  true,
// 				CurrentEpochEnded:     true,
// 			},
// 			{
// 				Identifier:            "week",
// 				StartTime:             time.Time{},
// 				Duration:              time.Hour * 24 * 7,
// 				CurrentEpoch:          0,
// 				CurrentEpochStartTime: time.Time{},
// 				EpochCountingStarted:  true,
// 				CurrentEpochEnded:     true,
// 			},
// 			{
// 				Identifier:            "day",
// 				StartTime:             time.Time{},
// 				Duration:              time.Hour * 24,
// 				CurrentEpoch:          0,
// 				CurrentEpochStartTime: time.Time{},
// 				EpochCountingStarted:  true,
// 				CurrentEpochEnded:     true,
// 			},
// 			{
// 				Identifier:            "hour",
// 				StartTime:             time.Time{},
// 				Duration:              time.Hour,
// 				CurrentEpoch:          0,
// 				CurrentEpochStartTime: time.Time{},
// 				EpochCountingStarted:  true,
// 				CurrentEpochEnded:     true,
// 			},
// 			{
// 				Identifier:            "minute",
// 				StartTime:             time.Time{},
// 				Duration:              time.Minute,
// 				CurrentEpoch:          0,
// 				CurrentEpochStartTime: time.Time{},
// 				EpochCountingStarted:  true,
// 				CurrentEpochEnded:     true,
// 			},
// 		},
// 	})

// 	epochInfomonth := app.EpochsKeeper.GetEpochInfo(ctx, "month")
// 	require.Equal(t, epochInfomonth.Identifier, "month")
// 	require.Equal(t, epochInfomonth.StartTime.UTC().String(), now.UTC().String())
// 	require.Equal(t, epochInfomonth.Duration, time.Hour*24*30)
// 	require.Equal(t, epochInfomonth.CurrentEpoch, int64(0))
// 	require.Equal(t, epochInfomonth.CurrentEpochStartTime.UTC().String(), ctx.BlockTime().UTC().String())
// 	require.Equal(t, epochInfomonth.EpochCountingStarted, true)
// 	require.Equal(t, epochInfomonth.CurrentEpochEnded, true)

// 	epochInfoWeek := app.EpochsKeeper.GetEpochInfo(ctx, "week")
// 	require.Equal(t, epochInfoWeek.Identifier, "week")
// 	require.Equal(t, epochInfoWeek.StartTime.UTC().String(), now.UTC().String())
// 	require.Equal(t, epochInfoWeek.Duration, time.Hour*24*7)
// 	require.Equal(t, epochInfoWeek.CurrentEpoch, int64(0))
// 	require.Equal(t, epochInfoWeek.CurrentEpochStartTime.UTC().String(), ctx.BlockTime().UTC().String())
// 	require.Equal(t, epochInfoWeek.EpochCountingStarted, true)
// 	require.Equal(t, epochInfoWeek.CurrentEpochEnded, true)

// 	epochInfoDay := app.EpochsKeeper.GetEpochInfo(ctx, "day")
// 	require.Equal(t, epochInfoDay.Identifier, "day")
// 	require.Equal(t, epochInfoDay.StartTime.UTC().String(), now.UTC().String())
// 	require.Equal(t, epochInfoDay.Duration, time.Hour*24)
// 	require.Equal(t, epochInfoDay.CurrentEpoch, int64(0))
// 	require.Equal(t, epochInfoDay.CurrentEpochStartTime.UTC().String(), ctx.BlockTime().UTC().String())
// 	require.Equal(t, epochInfoDay.EpochCountingStarted, true)
// 	require.Equal(t, epochInfoDay.CurrentEpochEnded, true)

// 	epochInfoHour := app.EpochsKeeper.GetEpochInfo(ctx, "hour")
// 	require.Equal(t, epochInfoHour.Identifier, "hour")
// 	require.Equal(t, epochInfoHour.StartTime.UTC().String(), now.UTC().String())
// 	require.Equal(t, epochInfoHour.Duration, time.Hour)
// 	require.Equal(t, epochInfoHour.CurrentEpoch, int64(0))
// 	require.Equal(t, epochInfoHour.CurrentEpochStartTime.UTC().String(), ctx.BlockTime().UTC().String())
// 	require.Equal(t, epochInfoHour.EpochCountingStarted, true)
// 	require.Equal(t, epochInfoHour.CurrentEpochEnded, true)

// 	epochInfoMinute := app.EpochsKeeper.GetEpochInfo(ctx, "minute")
// 	require.Equal(t, epochInfoMinute.Identifier, "minute")
// 	require.Equal(t, epochInfoMinute.StartTime.UTC().String(), now.UTC().String())
// 	require.Equal(t, epochInfoMinute.Duration, time.Minute)
// 	require.Equal(t, epochInfoMinute.CurrentEpoch, int64(0))
// 	require.Equal(t, epochInfoMinute.CurrentEpochStartTime.UTC().String(), ctx.BlockTime().UTC().String())
// 	require.Equal(t, epochInfoMinute.EpochCountingStarted, true)
// 	require.Equal(t, epochInfoMinute.CurrentEpochEnded, true)
// }
