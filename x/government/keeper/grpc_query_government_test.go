package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/government/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_GovernmentAddr(t *testing.T) {

	tests := []struct {
		name       string
		request    *types.QueryGovernmentAddrRequest
		government sdk.AccAddress
		wantErr    bool
	}{
		{
			"valid request",
			&types.QueryGovernmentAddrRequest{},
			governmentTestAddress,
			false,
		},
		{
			"empty request",
			nil,
			governmentTestAddress,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeperWithGovernmentAddress(t, tt.government)

			c := sdk.WrapSDKContext(ctx)

			got, err := k.GovernmentAddr(c, tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Keeper.GovernmentAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else {
				require.Equal(t, &types.QueryGovernmentAddrResponse{GovernmentAddress: tt.government.String()}, got)
			}

		})
	}
}
