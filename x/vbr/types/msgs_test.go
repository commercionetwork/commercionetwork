package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var funderAddr, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var validAmount = sdk.NewCoins(sdk.Coin{
	Denom:  BondDenom,
	Amount: sdk.NewInt(100),
})

var validMsgIncrementBlockRewardsPool = *NewMsgIncrementBlockRewardsPool(
	funderAddr.String(),
	validAmount,
)

func TestMsgIncrementBlockRewardsPool_Route(t *testing.T) {
	expected := ModuleName
	actual := validMsgIncrementBlockRewardsPool.Route()

	require.Equal(t, expected, actual)
}

func TestMsgIncrementBlockRewardsPool_Type(t *testing.T) {
	expected := MsgTypeIncrementBlockRewardsPool
	actual := validMsgIncrementBlockRewardsPool.Type()

	require.Equal(t, expected, actual)
}

func TestMsgIncrementBlockRewardsPool_GetSignBytes(t *testing.T) {
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&validMsgIncrementBlockRewardsPool))
	actual := validMsgIncrementBlockRewardsPool.GetSignBytes()

	require.Equal(t, expected, actual)
}

func TestMsgIncrementBlockRewardsPool_GetSigners(t *testing.T) {
	funderAddr, _ := sdk.AccAddressFromBech32(validMsgIncrementBlockRewardsPool.Funder)
	expected := []sdk.AccAddress{funderAddr}
	actual := validMsgIncrementBlockRewardsPool.GetSigners()

	require.Equal(t, expected, actual)
}
func TestMsgIncrementBlockRewardsPool_ValidateBasic(t *testing.T) {
	type fields struct {
		Funder string
		Amount sdk.Coins
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "ok",
			fields:  fields(validMsgIncrementBlockRewardsPool),
			wantErr: false,
		},
		{
			name: "invalid funder",
			fields: fields{
				Funder: "",
				Amount: validAmount,
			},
			wantErr: true,
		},
		{
			name: "invalid amount: zero",
			fields: fields{
				Funder: funderAddr.String(),
				Amount: sdk.NewCoins(sdk.NewCoin(BondDenom, sdk.ZeroInt())),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &MsgIncrementBlockRewardsPool{
				Funder: tt.fields.Funder,
				Amount: tt.fields.Amount,
			}
			if err := msg.ValidateBasic(); (err != nil) != tt.wantErr {
				t.Errorf("MsgIncrementBlockRewardsPool.ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// -------------------------
// --- MsgSetParams
// -------------------------
var governmentAddress, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")

var validMsgSetParams = *NewMsgSetParams(
	governmentAddress.String(),
	validDistrEpochIdentifier,
	validEarnRate,
)

func TestMsgSetParams_Route(t *testing.T) {
	expected := RouterKey
	actual := validMsgSetParams.Route()

	require.Equal(t, expected, actual)
}

func TestMsgSetParams_Type(t *testing.T) {
	expected := MsgTypeSetParams
	actual := validMsgSetParams.Type()

	require.Equal(t, expected, actual)
}

func TestMsgSetParams_GetSignBytes(t *testing.T) {
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&validMsgSetParams))
	actual := validMsgSetParams.GetSignBytes()

	require.Equal(t, expected, actual)
}

func TestMsgSetParams_GetSigners(t *testing.T) {
	expectedAddr, _ := sdk.AccAddressFromBech32(validMsgSetParams.Government)
	expected := []sdk.AccAddress{expectedAddr}
	actual := validMsgSetParams.GetSigners()

	require.Equal(t, expected, actual)
}

func TestMsgSetParams_ValidateBasic(t *testing.T) {
	type fields struct {
		Government           string
		DistrEpochIdentifier string
		EarnRate             sdk.Dec
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "ok",
			fields:  fields(validMsgSetParams),
			wantErr: false,
		},
		{
			name: "invalid Government",
			fields: fields{
				Government:           "",
				DistrEpochIdentifier: validMsgSetParams.DistrEpochIdentifier,
				EarnRate:             validMsgSetParams.EarnRate,
			},
			wantErr: true,
		},
		{
			name: "invalid DistrEpochIdentifier",
			fields: fields{
				Government:           validMsgSetParams.Government,
				DistrEpochIdentifier: "",
				EarnRate:             validMsgSetParams.EarnRate,
			},
			wantErr: true,
		},
		{
			name: "invalid EarnRate",
			fields: fields{
				Government:           validMsgSetParams.Government,
				DistrEpochIdentifier: validMsgSetParams.DistrEpochIdentifier,
				EarnRate:             invalidEarnRate,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &MsgSetParams{
				Government:           tt.fields.Government,
				DistrEpochIdentifier: tt.fields.DistrEpochIdentifier,
				EarnRate:             tt.fields.EarnRate,
			}
			if err := msg.ValidateBasic(); (err != nil) != tt.wantErr {
				t.Errorf("MsgSetParams.ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
