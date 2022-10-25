package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgIncrementBlockRewardsPool_Route(t *testing.T) {
	expected := ModuleName
	actual := ValidMsgIncrementBlockRewardsPool.Route()

	require.Equal(t, expected, actual)
}

func TestMsgIncrementBlockRewardsPool_Type(t *testing.T) {
	expected := MsgTypeIncrementBlockRewardsPool
	actual := ValidMsgIncrementBlockRewardsPool.Type()

	require.Equal(t, expected, actual)
}

func TestMsgIncrementBlockRewardsPool_GetSignBytes(t *testing.T) {
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&ValidMsgIncrementBlockRewardsPool))
	actual := ValidMsgIncrementBlockRewardsPool.GetSignBytes()

	require.Equal(t, expected, actual)
}

func TestMsgIncrementBlockRewardsPool_GetSigners(t *testing.T) {
	funderAddr, _ := sdk.AccAddressFromBech32(ValidMsgIncrementBlockRewardsPool.Funder)
	expected := []sdk.AccAddress{funderAddr}
	actual := ValidMsgIncrementBlockRewardsPool.GetSigners()

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
			fields:  fields(ValidMsgIncrementBlockRewardsPool),
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

func TestMsgSetParams_Route(t *testing.T) {
	expected := RouterKey
	actual := ValidMsgSetParams.Route()

	require.Equal(t, expected, actual)
}

func TestMsgSetParams_Type(t *testing.T) {
	expected := MsgTypeSetParams
	actual := ValidMsgSetParams.Type()

	require.Equal(t, expected, actual)
}

func TestMsgSetParams_GetSignBytes(t *testing.T) {
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&ValidMsgSetParams))
	actual := ValidMsgSetParams.GetSignBytes()

	require.Equal(t, expected, actual)
}

func TestMsgSetParams_GetSigners(t *testing.T) {
	expectedAddr, _ := sdk.AccAddressFromBech32(ValidMsgSetParams.Government)
	expected := []sdk.AccAddress{expectedAddr}
	actual := ValidMsgSetParams.GetSigners()

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
			fields:  fields(ValidMsgSetParams),
			wantErr: false,
		},
		{
			name: "invalid Government address format",
			fields: fields{
				Government:           "",
				DistrEpochIdentifier: ValidMsgSetParams.DistrEpochIdentifier,
				EarnRate:             ValidMsgSetParams.EarnRate,
			},
			wantErr: true,
		},
		{
			name: "invalid DistrEpochIdentifier",
			fields: fields{
				Government:           ValidMsgSetParams.Government,
				DistrEpochIdentifier: "",
				EarnRate:             ValidMsgSetParams.EarnRate,
			},
			wantErr: true,
		},
		{
			name: "invalid EarnRate",
			fields: fields{
				Government:           ValidMsgSetParams.Government,
				DistrEpochIdentifier: ValidMsgSetParams.DistrEpochIdentifier,
				EarnRate:             InvalidEarnRate,
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
