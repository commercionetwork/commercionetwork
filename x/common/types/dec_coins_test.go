package types

import (
	"testing"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestIsAllGTE(t *testing.T) {
	type args struct {
		coins      sdk.DecCoins
		otherCoins sdk.DecCoins
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "both empty",
			want: true,
		},
		{
			name: "other empty",
			args: args{
				coins: []sdk.DecCoin{sdk.NewDecCoin(sdk.DefaultBondDenom, math.OneInt())},
			},
			want: true,
		},
		{
			name: "coins empty",
			args: args{
				otherCoins: []sdk.DecCoin{sdk.NewDecCoin(sdk.DefaultBondDenom, math.OneInt())},
			},
		},
		{
			name: "coins equal to other",
			args: args{
				coins:      []sdk.DecCoin{sdk.NewDecCoin(sdk.DefaultBondDenom, math.OneInt())},
				otherCoins: []sdk.DecCoin{sdk.NewDecCoin(sdk.DefaultBondDenom, math.OneInt())},
			},
			want: true,
		},
		{
			name: "coins less than other",
			args: args{
				coins:      []sdk.DecCoin{sdk.NewDecCoin(sdk.DefaultBondDenom, math.ZeroInt())},
				otherCoins: []sdk.DecCoin{sdk.NewDecCoin(sdk.DefaultBondDenom, math.OneInt())},
			},
		},
		// {
		// 	name: "",
		// 	args: args{
		// 		coins:      []sdk.DecCoin{},
		// 		otherCoins: []sdk.DecCoin{},
		// 	},
		// 	want: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAllGTE(tt.args.coins, tt.args.otherCoins); got != tt.want {
				t.Errorf("IsAllGTE() = %v, want %v", got, tt.want)
			}
		})
	}
}
