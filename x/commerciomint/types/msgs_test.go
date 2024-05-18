package types

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var governmentAddress, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")

func TestMsgBasics(t *testing.T) {
	coinPos := sdk.NewCoin(CreditsDenom, math.NewInt(1))
	exchangeRate := math.LegacyNewDec(1)
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

	msgBurn := NewMsgBurnCCC(nil, "id", sdk.NewCoin("denom", math.NewInt(1)))
	require.Equal(t, "commerciomint", msgBurn.Route())
	require.Equal(t, "burnCCC", msgBurn.Type())
	require.Equal(t, 0, len(msgBurn.GetSigners()))
	require.NotNil(t, msgBurn.GetSignBytes())
}

func TestMsgMintCCC_ValidateBasic(t *testing.T) {
	uuid := "1480ab35-8544-405a-9729-595ae78c8fda"
	coinPos := sdk.NewCoin(CreditsDenom, math.NewInt(1))
	exchangeRate := math.LegacyNewDec(1)
	position := Position{"", 10, &coinPos, &time.Time{}, uuid, exchangeRate}
	require.Error(t, NewMsgMintCCC(position).ValidateBasic())
	//require.Error(t, NewMsgMintCCC(nil, sdk.NewCoins(sdk.NewInt64Coin(CreditsDenom, 100)), uuid).ValidateBasic())
	//coinPos = sdk.NewCoin("denom", math.NewInt(0))
	position = Position{testOwner.String(), 0, &coinPos, &time.Time{}, uuid, exchangeRate}
	require.Error(t, NewMsgMintCCC(position).ValidateBasic())
	//require.Error(t, NewMsgMintCCC(testOwner, sdk.NewCoins(), uuid).ValidateBasic())
	position = Position{testOwner.String(), 10, &coinPos, &time.Time{}, "", exchangeRate}
	require.Error(t, NewMsgMintCCC(position).ValidateBasic())
	//require.Error(t, NewMsgMintCCC(testOwner, sdk.NewCoins(), "").ValidateBasic())

	// ---------------------------------------
	// TODO control mint message
	/*
		coinPos = sdk.NewCoin("uatom", math.NewInt(1))
		position = Position{testOwner.String(), 10, &coinPos, "1", uuid, &exchangeRate}
		require.Error(t, NewMsgMintCCC(position).ValidateBasic())
	*/
	//require.Error(t, NewMsgMintCCC(testOwner, sdk.NewCoins(sdk.NewInt64Coin("atom", 100)), uuid).ValidateBasic())
	// ---------------------------------------

	position = Position{testOwner.String(), 10, &coinPos, &time.Time{}, uuid, exchangeRate}
	require.NoError(t, NewMsgMintCCC(position).ValidateBasic())
	//require.NoError(t, NewMsgMintCCC(testOwner, sdk.NewCoins(sdk.NewInt64Coin(CreditsDenom, 100)), uuid).ValidateBasic())
}

func TestMsgBurnCCC_ValidateBasic(t *testing.T) {
	uuid := "1480ab35-8544-405a-9729-595ae78c8fda"
	require.Error(t, NewMsgBurnCCC(nil, uuid, sdk.NewCoin(CreditsDenom, math.NewInt(100))).ValidateBasic())
	require.Error(t, NewMsgBurnCCC(testOwner, uuid, sdk.NewCoin("atom", math.NewInt(100))).ValidateBasic())
	require.Error(t, NewMsgBurnCCC(testOwner, "", sdk.NewCoin(CreditsDenom, math.NewInt(100))).ValidateBasic())
	require.NoError(t, NewMsgBurnCCC(testOwner, uuid, sdk.NewCoin(CreditsDenom, math.NewInt(100))).ValidateBasic())
}

func TestMsgSetParams_ValidateBasic(t *testing.T) {
	type fields struct {
		Signer string
		Params *Params
	}
	tests := []struct {
		name         string
		msgSetParams MsgSetParams
		wantErr      bool
	}{
		{
			name: "ok",
			msgSetParams: *NewMsgSetParams(
				governmentAddress.String(),
				validConversionRate,
				validFreezePeriod,
			),
			wantErr: false,
		},
		{
			name: "invalid government",
			msgSetParams: *NewMsgSetParams(
				"",
				validConversionRate,
				validFreezePeriod,
			),
			wantErr: true,
		},
		{
			name: "invalid conversion rate",
			msgSetParams: *NewMsgSetParams(
				governmentAddress.String(),
				invalidConversionRate,
				validFreezePeriod,
			),
			wantErr: true,
		},
		{
			name: "invalid freeze period",
			msgSetParams: *NewMsgSetParams(
				governmentAddress.String(),
				validConversionRate,
				invalidFreezePeriod,
			),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tt.msgSetParams
			if err := msg.ValidateBasic(); (err != nil) != tt.wantErr {
				t.Errorf("MsgSetParams.ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
