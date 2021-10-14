package types

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgBasics(t *testing.T) {
	coinPos := sdk.NewCoin("uccc", sdk.NewInt(1))
	exchangeRate := sdk.NewDec(1)
	position := Position{
		"1",
		10,
		&coinPos,
		&time.Time{},
		"1",
		exchangeRate,
	}
	msgMint := NewMsgMintCCC(position)

	require.Equal(t, "commerciomint", msgMint.Route())
	require.Equal(t, "mintCCC", msgMint.Type())
	require.Equal(t, 0, len(msgMint.GetSigners()))
	require.NotNil(t, msgMint.GetSignBytes())

	msgBurn := NewMsgBurnCCC(nil, "id", sdk.NewCoin("denom", sdk.NewInt(1)))
	require.Equal(t, "commerciomint", msgBurn.Route())
	require.Equal(t, "burnCCC", msgBurn.Type())
	require.Equal(t, 0, len(msgBurn.GetSigners()))
	require.NotNil(t, msgBurn.GetSignBytes())
}

func TestMsgMintCCC_ValidateBasic(t *testing.T) {
	uuid := "1480ab35-8544-405a-9729-595ae78c8fda"
	coinPos := sdk.NewCoin("uccc", sdk.NewInt(1))
	exchangeRate := sdk.NewDec(1)
	position := Position{"", 10, &coinPos, &time.Time{}, uuid, exchangeRate}
	require.Error(t, NewMsgMintCCC(position).ValidateBasic())
	//require.Error(t, NewMsgMintCCC(nil, sdk.NewCoins(sdk.NewInt64Coin("uccc", 100)), uuid).ValidateBasic())
	//coinPos = sdk.NewCoin("denom", sdk.NewInt(0))
	position = Position{testOwner.String(), 0, &coinPos, &time.Time{}, uuid, exchangeRate}
	require.Error(t, NewMsgMintCCC(position).ValidateBasic())
	//require.Error(t, NewMsgMintCCC(testOwner, sdk.NewCoins(), uuid).ValidateBasic())
	position = Position{testOwner.String(), 10, &coinPos, &time.Time{}, "", exchangeRate}
	require.Error(t, NewMsgMintCCC(position).ValidateBasic())
	//require.Error(t, NewMsgMintCCC(testOwner, sdk.NewCoins(), "").ValidateBasic())

	// ---------------------------------------
	// TODO control mint message
	/*
		coinPos = sdk.NewCoin("uatom", sdk.NewInt(1))
		position = Position{testOwner.String(), 10, &coinPos, "1", uuid, &exchangeRate}
		require.Error(t, NewMsgMintCCC(position).ValidateBasic())
	*/
	//require.Error(t, NewMsgMintCCC(testOwner, sdk.NewCoins(sdk.NewInt64Coin("atom", 100)), uuid).ValidateBasic())
	// ---------------------------------------

	position = Position{testOwner.String(), 10, &coinPos, &time.Time{}, uuid, exchangeRate}
	require.NoError(t, NewMsgMintCCC(position).ValidateBasic())
	//require.NoError(t, NewMsgMintCCC(testOwner, sdk.NewCoins(sdk.NewInt64Coin("uccc", 100)), uuid).ValidateBasic())
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
			msg := NewMsgSetCCCConversionRate(tt.signer, sdk.DecProto{Dec: tt.collateralRate})
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
			msg := NewMsgSetCCCFreezePeriod(tt.signer, tt.freezePeriod.String()) // TODO control cast
			require.Equal(t, "commerciomint", msg.Route())
			require.Equal(t, "setEtpsFreezePeriod", msg.Type())
			require.Equal(t, tt.wantErr, msg.ValidateBasic() != nil)
		})
	}
}
