package types

import (
	"testing"
	"time"

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

func TestMsgMintCCC_ValidateBasic(t *testing.T) {
	uuid := "1480ab35-8544-405a-9729-595ae78c8fda"
	require.Error(t, NewMsgMintCCC(nil, sdk.NewCoins(sdk.NewInt64Coin("uccc", 100)), uuid).ValidateBasic())
	require.Error(t, NewMsgMintCCC(testOwner, sdk.NewCoins(), uuid).ValidateBasic())
	require.Error(t, NewMsgMintCCC(testOwner, sdk.NewCoins(), "").ValidateBasic())
	require.Error(t, NewMsgMintCCC(testOwner, sdk.NewCoins(sdk.NewInt64Coin("atom", 100)), uuid).ValidateBasic())
	require.NoError(t, NewMsgMintCCC(testOwner, sdk.NewCoins(sdk.NewInt64Coin("uccc", 100)), uuid).ValidateBasic())
}

func TestMsgBurnCCC_ValidateBasic(t *testing.T) {
	uuid := "1480ab35-8544-405a-9729-595ae78c8fda"
	require.Error(t, NewMsgBurnCCC(nil, uuid, sdk.NewCoin("uccc", sdk.NewInt(100))).ValidateBasic())
	require.Error(t, NewMsgBurnCCC(testOwner, uuid, sdk.NewCoin("atom", sdk.NewInt(100))).ValidateBasic())
	require.Error(t, NewMsgBurnCCC(testOwner, "", sdk.NewCoin("uccc", sdk.NewInt(100))).ValidateBasic())
	require.NoError(t, NewMsgBurnCCC(testOwner, uuid, sdk.NewCoin("uccc", sdk.NewInt(100))).ValidateBasic())
}

func TestMsgSetCCCConversionRate_ValidateBasic(t *testing.T) {
	type fields struct {
	}
	tests := []struct {
		name           string
		signer         sdk.AccAddress
		collateralRate sdk.Dec
		wantErr        bool
	}{
		{"empty signer", nil, sdk.NewDec(2), true},
		{"ok", []byte("test"), sdk.NewDec(2), false},
		{"zero collateral rate", []byte("test"), sdk.NewDec(0), true},
		{"negative collateral rate", []byte("test"), sdk.NewDec(-1), true},
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

func TestMsgSetCCCFreezePeriod_ValidateBasic(t *testing.T) {
	type fields struct {
	}
	tests := []struct {
		name         string
		signer       sdk.AccAddress
		freezePeriod time.Duration
		wantErr      bool
	}{
		{"ok", []byte("test"), 60, false},
		{"Negative duration", []byte("test"), -60, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			msg := NewMsgSetCCCFreezePeriod(tt.signer, tt.freezePeriod)
			require.Equal(t, "commerciomint", msg.Route())
			require.Equal(t, "setEtpsFreezePeriod", msg.Type())
			require.Equal(t, tt.wantErr, msg.ValidateBasic() != nil)
		})
	}
}
