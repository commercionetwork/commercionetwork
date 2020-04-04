package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgBasics(t *testing.T) {
	require.Equal(t, "commerciomint", MsgOpenCdp{}.Route())
	require.Equal(t, "openCdp", MsgOpenCdp{}.Type())
	require.Equal(t, 1, len(MsgOpenCdp{}.GetSigners()))
	require.NotNil(t, MsgOpenCdp{}.GetSignBytes())

	msg := NewMsgCloseCdp(nil, 0)
	require.Equal(t, "commerciomint", msg.Route())
	require.Equal(t, "closeCdp", msg.Type())
	require.Equal(t, 1, len(msg.GetSigners()))
	require.NotNil(t, msg.GetSignBytes())
}

func TestMsgOpenCdp_ValidateBasic(t *testing.T) {
	require.Error(t, NewMsgOpenCdp(nil, sdk.NewInt64Coin("atom", 100)).ValidateBasic())
	require.Error(t, NewMsgOpenCdp(testOwner, sdk.Coin{}).ValidateBasic())
	require.NoError(t, NewMsgOpenCdp(testOwner, sdk.NewInt64Coin("atom", 100)).ValidateBasic())
}

func TestMsgSetCdpCollateralRate_ValidateBasic(t *testing.T) {
	type fields struct {
		Signer            sdk.AccAddress
		CdpCollateralRate sdk.Dec
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"empty signer", fields{nil, sdk.NewDec(2)}, true},
		{"ok", fields{[]byte("test"), sdk.NewDec(2)}, false},
		{"zero collateral rate", fields{[]byte("test"), sdk.NewDec(0)}, true},
		{"negative collateral rate", fields{[]byte("test"), sdk.NewDec(-1)}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			msg := NewMsgSetCdpCollateralRate(tt.fields.Signer, tt.fields.CdpCollateralRate)
			require.Equal(t, "commerciomint", msg.Route())
			require.Equal(t, "setCdpCollateralRate", msg.Type())
			require.Equal(t, tt.wantErr, msg.ValidateBasic() != nil)
		})
	}
}
