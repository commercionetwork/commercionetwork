package keeper

import (
	"reflect"
	"testing"

	//"cosmossdk.io/simapp"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/commercionetwork/commercionetwork/x/government/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"
)

func TestNewQuerier(t *testing.T) {
	t.Run("default request", func(t *testing.T) {
		k, ctx := setupKeeperWithGovernmentAddress(t, governmentTestAddress)

		app := simapp.Setup(false)
		legacyAmino := app.LegacyAmino()
		querier := NewQuerier(*k, legacyAmino)
		path := []string{"abcd"}
		_, err := querier(ctx, path, abci.RequestQuery{})
		require.Error(t, err)
	})
}

func Test_queryGetGovernmentAddress(t *testing.T) {
	bz, _ := codec.MarshalJSONIndent(codec.NewLegacyAmino(), types.QueryGovernmentAddrResponse{GovernmentAddress: governmentTestAddress.String()})
	tests := []struct {
		name             string
		legacyQuerierCdc *codec.LegacyAmino
		want             []byte
		wantErr          bool
	}{
		{
			name:             "Query Government",
			legacyQuerierCdc: codec.NewLegacyAmino(),
			want:             bz,
			wantErr:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeperWithGovernmentAddress(t, governmentTestAddress)

			app := simapp.Setup(false)
			legacyAmino := app.LegacyAmino()
			querier := NewQuerier(*k, legacyAmino)
			path := []string{types.QueryGovernmentAddress}
			got, err := querier(ctx, path, abci.RequestQuery{})

			if (err != nil) != tt.wantErr {
				t.Errorf("queryGetGovernmentAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryGetGovernmentAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
