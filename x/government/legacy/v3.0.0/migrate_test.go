package v3_0_0

import (
	"reflect"
	"testing"

	v220government "github.com/commercionetwork/commercionetwork/x/government/legacy/v2.2.0"
	"github.com/commercionetwork/commercionetwork/x/government/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	testGov, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
	testTumblerAddress, _    = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
)

func TestMigrate(t *testing.T) {
	type args struct {
		v220GenState v220government.GenesisState
	}
	tests := []struct {
		name string
		args args
		want *types.GenesisState
	}{
		{
			name: "empty",
			args: args{v220GenState: v220government.GenesisState{}},
			want: &types.GenesisState{},
		},
		{
			name: "genesis state correctly migrated",
			args: args{
				v220GenState: v220government.GenesisState{
					GovernmentAddress: testGov,
					TumblerAddress: testTumblerAddress,
				},
			},
			want: &types.GenesisState{
				GovernmentAddress: testGov.String(),
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