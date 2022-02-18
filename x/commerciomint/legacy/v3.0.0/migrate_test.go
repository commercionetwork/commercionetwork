package v3_0_0

import (
	"reflect"
	"testing"
	"time"

	v220commerciomint "github.com/commercionetwork/commercionetwork/x/commerciomint/legacy/v2.2.0"
	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	testUser01, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
	testUser02, _ = sdk.AccAddressFromBech32("cosmos14lultfckehtszvzw4ehu0apvsr77afvyhgqhwh")
	timeDuration  = time.Duration(18000000000)
	cAt01         = time.Date(2021, time.April,
		11, 21, 34, 01, 0, time.UTC)
	cAt02 = time.Date(2021, time.February,
		11, 21, 34, 01, 0, time.UTC)
)

func TestMigrate(t *testing.T) {
	type args struct {
		v220GenState v220commerciomint.GenesisState
	}
	tests := []struct {
		name string
		args args
		want *types.GenesisState
	}{
		{
			name: "empty",
			args: args{v220GenState: v220commerciomint.GenesisState{}},
			want: &types.GenesisState{},
		},
		{
			name: "genesis state corectly migrated",
			args: args{
				v220GenState: v220commerciomint.GenesisState{
					Positions: []v220commerciomint.Position{
						v220commerciomint.Position{
							Owner:      testUser01,
							Collateral: sdk.NewInt(10),
							Credits: sdk.Coin{
								Denom:  "ucommerico",
								Amount: sdk.NewInt(5),
							},
							CreatedAt:    cAt01,
							ID:           "some-id-01",
							ExchangeRate: sdk.NewDec(2),
						},
						v220commerciomint.Position{
							Owner:      testUser02,
							Collateral: sdk.NewInt(1000),
							Credits: sdk.Coin{
								Denom:  "ucommerico",
								Amount: sdk.NewInt(500),
							},
							CreatedAt:    cAt02,
							ID:           "some-id-02",
							ExchangeRate: sdk.NewDec(2),
						},
					},
					LiquidityPoolAmount: sdk.Coins{
						sdk.Coin{
							Denom:  "ucommercio",
							Amount: sdk.NewInt(1000000),
						},
					},
					CollateralRate: sdk.NewDec(2),
					FreezePeriod:   timeDuration,
				},
			},
			want: &types.GenesisState{
				Positions: []*types.Position{
					&types.Position{
						Owner:      testUser01.String(),
						Collateral: 10,
						Credits: &sdk.Coin{
							Denom:  "ucommerico",
							Amount: sdk.NewInt(5),
						},
						CreatedAt:    &cAt01,
						ID:           "some-id-01",
						ExchangeRate: sdk.NewDec(2),
					},
					&types.Position{
						Owner:      testUser02.String(),
						Collateral: 1000,
						Credits: &sdk.Coin{
							Denom:  "ucommerico",
							Amount: sdk.NewInt(500),
						},
						CreatedAt:    &cAt02,
						ID:           "some-id-02",
						ExchangeRate: sdk.NewDec(2),
					},
				},
				PoolAmount: sdk.Coins{
					sdk.Coin{
						Denom:  "ucommercio",
						Amount: sdk.NewInt(1000000),
					},
				},
				Params: types.Params{
					ConversionRate: sdk.NewDec(2),
					FreezePeriod:   timeDuration,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Migrate(tt.args.v220GenState); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Migrate() = %v, want %v", got, tt.want)
			}
		})
	}
}
