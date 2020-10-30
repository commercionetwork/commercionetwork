package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgBasics(t *testing.T) {
	require.Equal(t, "commerciomint", MsgMintCCC{}.Route())
	require.Equal(t, "mintCCC", MsgMintCCC{}.Type())
	require.Equal(t, 1, len(MsgMintCCC{}.GetSigners()))
	require.NotNil(t, MsgMintCCC{}.GetSignBytes())

	msg := NewMsgBurnCCC(nil, "id", sdk.NewCoin("denom", sdk.NewInt(1)))
	require.Equal(t, "commerciomint", msg.Route())
	require.Equal(t, "burnCCC", msg.Type())
	require.Equal(t, 1, len(msg.GetSigners()))
	require.NotNil(t, msg.GetSignBytes())
}

func TestMsgOpenCdp_ValidateBasic(t *testing.T) {
	require.Error(t, NewMsgMintCCC(nil, sdk.NewCoins(sdk.NewInt64Coin("atom", 100))).ValidateBasic())
	require.Error(t, NewMsgMintCCC(testOwner, sdk.NewCoins()).ValidateBasic())
	require.NoError(t, NewMsgMintCCC(testOwner, sdk.NewCoins(sdk.NewInt64Coin("uccc", 100))).ValidateBasic())
}

func TestMsgSetCdpCollateralRate_ValidateBasic(t *testing.T) {
	type fields struct {
	}
	tests := []struct {
		name           string
		signer         sdk.AccAddress
		collateralRate sdk.Int
		wantErr        bool
	}{
		{"empty signer", nil, sdk.NewInt(2), true},
		{"ok", []byte("test"), sdk.NewInt(2), false},
		{"zero collateral rate", []byte("test"), sdk.NewInt(0), true},
		{"negative collateral rate", []byte("test"), sdk.NewInt(-1), true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			msg := NewMsgSetCCCConversionRate(tt.signer, tt.collateralRate)
			require.Equal(t, "commerciomint", msg.Route())
			require.Equal(t, "setEtpsConversionRate", msg.Type())
			require.Equal(t, tt.wantErr, msg.ValidateBasic() != nil)
		})
	}
}
