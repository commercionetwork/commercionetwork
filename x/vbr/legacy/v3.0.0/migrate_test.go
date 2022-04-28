package v3_0_0

import (
	"reflect"
	"testing"

	v220vbr "github.com/commercionetwork/commercionetwork/x/vbr/legacy/v2.2.0"
	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var params = types.NewParams(types.EpochDay, sdk.NewDecWithPrec(5, 1))

func TestMigrate(t *testing.T) {
	type args struct {
		v220GenState v220vbr.GenesisState
	}
	tests := []struct {
		name string
		args args
		want *types.GenesisState
	}{
		{
			name: "empty",
			args: args{v220GenState: v220vbr.GenesisState{}},
			want: &types.GenesisState{
				Params: params,
			},
		},
		{
			name: "genesis state correctly migrated",
			args: args{
				v220GenState: v220vbr.GenesisState{
					PoolAmount: sdk.DecCoins{
						{
							Denom:  types.BondDenom,
							Amount: sdk.NewDec(1000000),
						},
					},
					RewardRate:        sdk.NewDecWithPrec(112, 5),
					AutomaticWithdraw: true,
				},
			},
			want: &types.GenesisState{
				PoolAmount: sdk.DecCoins{
					{
						Denom:  types.BondDenom,
						Amount: sdk.NewDec(1000000),
					},
				},
				Params: params,
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
