package keeper

import (
	"reflect"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/government/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Test_msgServer_SetGovAddress(t *testing.T) {

	tests := []struct {
		name    string
		msg     *types.MsgSetGovAddress
		want    *types.MsgSetGovAddressResponse
		wantErr bool
	}{
		{
			name:    "Government address setup",
			msg:     types.NewMsgSetGovAddress(),
			want:    &types.MsgSetGovAddressResponse{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeperWithGovernmentAddress(t, governmentTestAddress)

			msgServer := NewMsgServerImpl(*k)

			got, err := msgServer.SetGovAddress(sdk.WrapSDKContext(ctx), tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("msgServer.SetGovAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("msgServer.SetGovAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
