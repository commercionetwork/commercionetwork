package types_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Test vars
var user, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var sender, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")

// ----------------------
// --- MsgInviteUser
// ----------------------

var msgInviteUser = types.NewMsgInviteUser(user, sender)

func TestMsgInviteUser_Route(t *testing.T) {
	require.Equal(t, types.RouterKey, msgInviteUser.Route())
}

func TestMsgInviteUser_Type(t *testing.T) {
	require.Equal(t, types.MsgTypeInviteUser, msgInviteUser.Type())
}

func TestMsgInviteUser_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   types.MsgInviteUser
		error error
	}{
		{
			name:  "Valid message returns no error",
			msg:   msgInviteUser,
			error: nil,
		},
		{
			name:  "Missing recipient returns error",
			msg:   types.MsgInviteUser{Recipient: nil, Sender: sender},
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid recipient address: "),
		},
		{
			name:  "Missing sender returns error",
			msg:   types.MsgInviteUser{Recipient: user, Sender: nil},
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid sender address: "),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.error, test.msg.ValidateBasic())
		})
	}
}

func TestMsgInviteUser_GetSignBytes(t *testing.T) {
	actual := msgInviteUser.GetSignBytes()
	expected := sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msgInviteUser))
	require.Equal(t, expected, actual)
}

func TestMsgInviteUser_GetSigners(t *testing.T) {
	actual := msgInviteUser.GetSigners()
	require.Equal(t, 1, len(actual))
	require.Equal(t, msgInviteUser.Recipient, actual[0])
}

func TestMsgInviteUser_UnmarshalJson(t *testing.T) {
	json := `{"type":"commercio/MsgInviteUser","value":{"recipient":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","sender":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"}}`

	var msg types.MsgInviteUser
	types.ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	require.Equal(t, user, msg.Recipient)
	require.Equal(t, sender, msg.Sender)
}

// ---------------------------
// --- MsgSetUserVerified
// ---------------------------

var msgSetUserVerified = types.NewMsgSetUserVerified(user, tsp)

func TestMsgSetUserVerified_Route(t *testing.T) {
	require.Equal(t, types.RouterKey, msgSetUserVerified.Route())
}

func TestMsgSetUserVerified_Type(t *testing.T) {
	require.Equal(t, types.MsgTypeSetUserVerified, msgSetUserVerified.Type())
}

func TestMsgSetUserVerified_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   types.MsgSetUserVerified
		error error
	}{
		{
			name:  "Valid msg does not return any error",
			msg:   msgSetUserVerified,
			error: nil,
		},
		{
			name:  "Missing user returns error",
			msg:   types.NewMsgSetUserVerified(nil, tsp),
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid user address: "),
		},
		{
			name:  "Missing signer returns error",
			msg:   types.NewMsgSetUserVerified(user, nil),
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid verifier address: "),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.error, test.msg.ValidateBasic())
		})
	}
}

func TestMsgSetUserVerified_GetSignBytes(t *testing.T) {
	json := `{"type":"commercio/MsgSetUserVerified","value":{"user":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","verifier":"cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee"}}`
	require.Equal(t, json, string(msgSetUserVerified.GetSignBytes()))
}

func TestMsgSetUserVerified_GetSigners(t *testing.T) {
	actual := msgSetUserVerified.GetSigners()
	require.Equal(t, 1, len(actual))
	require.Equal(t, msgSetUserVerified.Verifier, actual[0])
}

func TestMsgSetUserVerified_UnmarshalJson(t *testing.T) {
	json := `{"type":"commercio/MsgSetUserVerified","value":{"user":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","verifier":"cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee"}}`

	var msg types.MsgSetUserVerified
	types.ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	require.Equal(t, user, msg.User)
	require.Equal(t, tsp, msg.Verifier)
}

// --------------------------------
// --- MsgDepositIntoLiquidityPool
// --------------------------------

var amount = sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))
var msgDepositIntoLiquidityPool = types.NewMsgDepositIntoLiquidityPool(amount, user)

func TestMsgDepositIntoLiquidityPool_Route(t *testing.T) {
	require.Equal(t, types.RouterKey, msgDepositIntoLiquidityPool.Route())
}

func TestMsgDepositIntoLiquidityPool_Type(t *testing.T) {
	require.Equal(t, types.MsgTypesDepositIntoLiquidityPool, msgDepositIntoLiquidityPool.Type())
}

func TestMsgDepositIntoLiquidityPool_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   types.MsgDepositIntoLiquidityPool
		error error
	}{
		{
			name:  "Valid message returns no error",
			msg:   msgDepositIntoLiquidityPool,
			error: nil,
		},
		{
			name:  "Missing deposit returns error",
			msg:   types.MsgDepositIntoLiquidityPool{Depositor: nil, Amount: amount},
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid depositor address: "),
		},
		{
			name:  "Empty deposit amount returns error",
			msg:   types.MsgDepositIntoLiquidityPool{Depositor: user, Amount: nil},
			error: sdk.ErrInvalidCoins("Invalid deposit amount: "),
		},
		{
			name: "Negative deposit amount returns error",
			msg: types.MsgDepositIntoLiquidityPool{
				Depositor: user,
				Amount:    sdk.Coins{sdk.Coin{Denom: "uatom", Amount: sdk.NewInt(-100)}},
			},
			error: sdk.ErrInvalidCoins("Invalid deposit amount: -100uatom"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.error, test.msg.ValidateBasic())
		})
	}
}

func TestMsgDepositIntoLiquidityPool_GetSignBytes(t *testing.T) {
	actual := msgDepositIntoLiquidityPool.GetSignBytes()
	expected := sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msgDepositIntoLiquidityPool))
	require.Equal(t, expected, actual)
}

func TestMsgDepositIntoLiquidityPool_GetSigners(t *testing.T) {
	actual := msgDepositIntoLiquidityPool.GetSigners()
	require.Equal(t, 1, len(actual))
	require.Equal(t, msgDepositIntoLiquidityPool.Depositor, actual[0])
}

func TestMsgDepositIntoLiquidityPool_UnmarshalJson(t *testing.T) {
	json := `{"type":"commercio/MsgDepositIntoLiquidityPool","value":{"depositor":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","amount":[{"denom":"uatom","amount":"100"}]}}`

	var msg types.MsgDepositIntoLiquidityPool
	types.ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	require.Equal(t, user, msg.Depositor)
	require.Equal(t, amount, msg.Amount)
}

// --------------------------------
// --- MsgAddTsp
// --------------------------------

var government, _ = sdk.AccAddressFromBech32("cosmos1ct4ym78j7ksv9weyua4mzlksgwc9qq7q3wvhqg")
var tsp, _ = sdk.AccAddressFromBech32("cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee")
var msgAddTsp = types.MsgAddTsp{
	Government: government,
	Tsp:        tsp,
}

func TestMsgAddTsp_Route(t *testing.T) {
	require.Equal(t, types.RouterKey, msgAddTsp.Route())
}

func TestMsgAddTsp_Type(t *testing.T) {
	require.Equal(t, types.MsgTypeAddTsp, msgAddTsp.Type())
}

func TestMsgAddTsp_ValidateBasic_ValidMsg(t *testing.T) {
	tests := []struct {
		name  string
		msg   types.MsgAddTsp
		error error
	}{
		{
			name:  "Valid message does not return any error",
			msg:   msgAddTsp,
			error: nil,
		},
		{
			name:  "Missing government returns error",
			msg:   types.MsgAddTsp{Government: nil, Tsp: tsp},
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid government address: "),
		},
		{
			name:  "Missing tsp returns error",
			msg:   types.MsgAddTsp{Government: government, Tsp: nil},
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid TSP address: "),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.error, test.msg.ValidateBasic())
		})
	}
}

func TestMsgAddTsp_GetSignBytes(t *testing.T) {
	actual := msgAddTsp.GetSignBytes()
	expected := sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msgAddTsp))
	require.Equal(t, expected, actual)
}

func TestMsgAddTsp_GetSigners(t *testing.T) {
	actual := msgAddTsp.GetSigners()
	require.Equal(t, 1, len(actual))
	require.Equal(t, msgAddTsp.Government, actual[0])
}

func TestMsgAddTsp_UnmarshalJson(t *testing.T) {
	json := `{"type":"commercio/MsgAddTsp","value":{"government":"cosmos1ct4ym78j7ksv9weyua4mzlksgwc9qq7q3wvhqg","tsp":"cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee"}}`

	var msg types.MsgAddTsp
	types.ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	require.Equal(t, tsp, msg.Tsp)
	require.Equal(t, government, msg.Government)
}

// ---------------------------
// --- MsgBuyMemberships
// ---------------------------

var testBuyer, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var TestMembershipType = "bronze"
var msgBuyMembership = types.NewMsgBuyMembership(TestMembershipType, testBuyer)

func TestMsgBuyMembership_Route(t *testing.T) {
	require.Equal(t, types.RouterKey, msgBuyMembership.Route())
}

func TestMsgBuyMembership_Type(t *testing.T) {
	require.Equal(t, types.MsgTypeBuyMembership, msgBuyMembership.Type())
}

func TestMsgBuyMembership_ValidateBasic_AllFieldsCorrect(t *testing.T) {
	tests := []struct {
		name  string
		msg   types.MsgBuyMembership
		error error
	}{
		{
			name:  "Valid message does not return any error",
			msg:   msgBuyMembership,
			error: nil,
		},
		{
			name:  "Missing buyer returns error",
			msg:   types.NewMsgBuyMembership(TestMembershipType, nil),
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid buyer address: "),
		},
		{
			name:  "Missing membership returns error",
			msg:   types.NewMsgBuyMembership("", testBuyer),
			error: sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid membership type: "),
		},
		{
			name:  "Invalid membership returns error",
			msg:   types.NewMsgBuyMembership("grn", testBuyer),
			error: sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid membership type: grn"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.error, test.msg.ValidateBasic())
		})
	}
}

func TestMsgBuyMembership_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgBuyMembership","value":{"buyer":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","membership_type":"bronze"}}`
	require.Equal(t, expected, string(msgBuyMembership.GetSignBytes()))
}

func TestMsgBuyMembership_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgBuyMembership.Buyer}
	require.Equal(t, expected, msgBuyMembership.GetSigners())
}
