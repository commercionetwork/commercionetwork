package keeper

import (
	"time"
	"testing"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"

	"github.com/commercionetwork/commercionetwork/x/upgrade/types"
	governmentKeeper "github.com/commercionetwork/commercionetwork/x/government/keeper"
	governmentTypes "github.com/commercionetwork/commercionetwork/x/government/types"
	upgradeKeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// Testing variables
var governmentTestAddress, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var notGovernmentAddress, _ = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")

func setupKeeper(t testing.TB) (*Keeper, sdk.Context) {
	storeKeys := sdk.NewKVStoreKeys(
			types.StoreKey,
			governmentTypes.GovernmentStoreKey,
			upgradeTypes.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	memStoreKeyGov := storetypes.NewMemoryStoreKey(governmentTypes.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	for _, key := range storeKeys {
		stateStore.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	}
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(memStoreKeyGov, sdk.StoreTypeMemory, nil)
	
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, tmproto.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	ctx = ctx.WithBlockTime(time.Now())

	registry := codectypes.NewInterfaceRegistry()
	govk := governmentKeeper.NewKeeper(codec.NewProtoCodec(registry), storeKeys[governmentTypes.GovernmentStoreKey], memStoreKeyGov)
	upgradek := upgradeKeeper.NewKeeper(map[int64]bool{}, storeKeys[upgradeTypes.StoreKey], &codec.ProtoCodec{}, "")
	
	keeper := NewKeeper(
		codec.NewProtoCodec(registry),
		storeKeys[types.StoreKey],
		memStoreKey,
		*govk,
		upgradek,
	)

	_ = govk.SetGovernmentAddress(ctx, governmentTestAddress)
	store := ctx.KVStore(storeKeys[governmentTypes.GovernmentStoreKey])

	store.Set([]byte(governmentTypes.GovernmentStoreKey), governmentTestAddress)

	return keeper, ctx
}


func TestKeeper_GetUpgradePlan(t *testing.T) {

	tests := []struct {
		name       string
		planExists bool
	}{
		{"a correctly scheduled plan exists",
			true,
		},
		{"no plan scheduled",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			var plan upgradeTypes.Plan
			if tt.planExists {
				plan = upgradeTypes.Plan{
					Name:   "name",
					Info:   "info info info",
					Height: 100,
				}
				err := k.UpgradeKeeper.ScheduleUpgrade(ctx, plan)
				require.NoError(t, err)
			} else {
				plan = upgradeTypes.Plan{Name: "", Time: time.Time{}, Height: 0, Info: ""}
			}

			upgradePlan, havePlan := k.GetUpgradePlan(ctx)

			require.Equal(t, tt.planExists, havePlan)
			require.Equal(t, plan, upgradePlan)
		})
	}
}

func TestKeeper_ScheduleUpgrade(t *testing.T) {

	timePast := time.Now()
	timeFuture := time.Now().Add(time.Duration(123456789))

	tests := []struct {
		name                 string
		proposedByGovernment bool
		time                 *time.Time
		height               int64
		wantErr              bool
	}{
		{"government schedules with future block height",
			true,
			nil,
			100,
			false,
		},
		{"government schedules with future time",
			true,
			&timeFuture,
			-1,
			false,
		},
		{"government schedules with time that has already passed yields error",
			true,
			&timePast,
			-1,
			true,
		},
		{"fake government trying to schedule an upgrade yields error",
			false,
			nil,
			100,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			plan := upgradeTypes.Plan{
				Name: "name",
				Info: "info info info",
			}

			if tt.time != nil {
				plan.Time = *tt.time
			} else {
				plan.Height = tt.height
			}

			var proposer sdk.AccAddress
			if tt.proposedByGovernment {
				proposer = governmentTestAddress
			} else {
				proposer = notGovernmentAddress
			}

			err := k.ScheduleUpgradeGov(ctx, proposer, plan)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				upgradePlan, havePlan := k.GetUpgradePlan(ctx)
				require.True(t, havePlan)
				require.NotNil(t, upgradePlan)
			}
		})
	}
}

func TestKeeper_DeleteUpgrade(t *testing.T) {

	tests := []struct {
		name                 string
		proposedByGovernment bool
		planExists           bool
		wantErr              bool
	}{
		{"government deletes the scheduled plan",
			true,
			true,
			false,
		},
		{"government deletes nothing",
			true,
			false,
			false,
		},
		{"fake government trying to delete an upgrade raises error",
			false,
			true,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			var plan upgradeTypes.Plan
			if tt.planExists {
				plan = upgradeTypes.Plan{
					Name:   "name",
					Info:   "info info info",
					Height: 100,
				}
				err := k.ScheduleUpgradeGov(ctx, governmentTestAddress, plan)
				require.NoError(t, err)
			}
			nilPlan := upgradeTypes.Plan{Name: "", Time: time.Time{}, Height: 0, Info: ""}

			var proposer sdk.AccAddress
			if tt.proposedByGovernment {
				proposer = governmentTestAddress
			} else {
				proposer = notGovernmentAddress
			}

			err := k.DeleteUpgradeGov(ctx, proposer)

			if tt.wantErr {
				require.Error(t, err)

				upgradePlan, havePlan := k.GetUpgradePlan(ctx)
				require.True(t, havePlan)
				require.Equal(t, plan, upgradePlan)

			} else {
				require.NoError(t, err)

				upgradePlan, havePlan := k.GetUpgradePlan(ctx)
				require.False(t, havePlan)
				require.Equal(t, nilPlan, upgradePlan)
			}
		})
	}
}
