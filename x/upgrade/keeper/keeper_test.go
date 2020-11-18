package keeper

import (
	governmentKeeper "github.com/commercionetwork/commercionetwork/x/government/keeper"
	governmentTypes "github.com/commercionetwork/commercionetwork/x/government/types"
	"github.com/commercionetwork/commercionetwork/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
	"testing"
	"time"
)

func TestKeeper_GetUpgradePlan(t *testing.T) {

	tests := []struct {
		name       string
		planExists bool
	}{
		{"government deletes the scheduled plan",
			true,
		},
		{"government deletes nothing",
			true,
		},
		{"non-government address trying to delete raises error",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k := SetupTestInput(true)

			var plan upgrade.Plan
			if tt.planExists {
				plan = upgrade.Plan{
					Name:   "name",
					Info:   "info info info",
					Height: 100,
				}
				err := k.UpgradeKeeper.ScheduleUpgrade(ctx, plan)
				require.NoError(t, err)
			} else {
				plan = upgrade.Plan{Name: "", Time: time.Time{}, Height: 0, Info: ""}
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
		{"government schedules with time that has already passed raises error",
			true,
			&timePast,
			-1,
			true,
		},
		{"non-government address trying to schedule raises error",
			false,
			nil,
			100,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k := SetupTestInput(true)

			plan := upgrade.Plan{
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

			err := k.ScheduleUpgrade(ctx, proposer, plan)

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
		{"non-government address trying to delete raises error",
			false,
			true,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k := SetupTestInput(true)

			var plan upgrade.Plan
			if tt.planExists {
				plan = upgrade.Plan{
					Name:   "name",
					Info:   "info info info",
					Height: 100,
				}
				err := k.ScheduleUpgrade(ctx, governmentTestAddress, plan)
				require.NoError(t, err)
			}
			nilPlan := upgrade.Plan{Name: "", Time: time.Time{}, Height: 0, Info: ""}

			var proposer sdk.AccAddress
			if tt.proposedByGovernment {
				proposer = governmentTestAddress
			} else {
				proposer = notGovernmentAddress
			}

			err := k.DeleteUpgrade(ctx, proposer)

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

// This function creates an environment to test modules
func SetupTestInput(setGovAddress bool) (cdc *codec.Codec, ctx sdk.Context, keeper Keeper) {

	memDB := db.NewMemDB()
	cdc = testCodec()

	// Store keys
	keyGovernment := sdk.NewKVStoreKey(governmentTypes.GovernmentStoreKey)
	keyUpgrade := sdk.NewKVStoreKey(upgrade.StoreKey)

	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(keyGovernment, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyUpgrade, sdk.StoreTypeIAVL, memDB)
	_ = ms.LoadLatestVersion()

	ctx = sdk.NewContext(ms, abci.Header{Time: time.Now(), Height: 1, ChainID: "test-chain-id"}, false, log.NewNopLogger())

	govk := governmentKeeper.NewKeeper(cdc, keyGovernment)
	upgk := upgrade.NewKeeper(map[int64]bool{}, keyUpgrade, cdc)
	customUpgradeKeeper := NewKeeper(cdc, govk, upgk)

	if setGovAddress {
		_ = govk.SetGovernmentAddress(ctx, governmentTestAddress)
	}
	return cdc, ctx, customUpgradeKeeper
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	types.RegisterCodec(cdc)
	cdc.Seal()

	return cdc
}

// Testing variables
var governmentTestAddress, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var notGovernmentAddress, _ = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")
