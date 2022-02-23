package types_test

import (
	"testing"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

// Test vars
var user, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var sender, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")

// ----------------------
// --- MsgInviteUser
// ----------------------

var msgInviteUser = types.NewMsgInviteUser(user.String(), sender.String())

func TestMsgInviteUser_Route(t *testing.T) {
	require.Equal(t, types.RouterKey, msgInviteUser.Route())
}

func TestMsgInviteUser_Type(t *testing.T) {
	require.Equal(t, types.MsgTypeInviteUser, msgInviteUser.Type())
}

func TestMsgInviteUser_ValidateBasic(t *testing.T) {
	missRecipient := types.MsgInviteUser{Recipient: "", Sender: sender.String()}
	missSender := types.MsgInviteUser{Recipient: user.String(), Sender: ""}
	tests := []struct {
		name  string
		msg   types.MsgInviteUser
		error error
	}{
		{
			name:  "Valid message returns no error",
			msg:   *msgInviteUser,
			error: nil,
		},
		{
			name:  "Missing recipient returns error",
			msg:   missRecipient,
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid recipient address:  (empty address string is not allowed)"),
		},
		{
			name:  "Missing sender returns error",
			msg:   missSender,
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid sender address:  (empty address string is not allowed)"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if test.error != nil {
				require.Equal(t, test.error.Error(), test.msg.ValidateBasic().Error())
			} else {
				require.NoError(t, test.msg.ValidateBasic())
			}
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
	require.Equal(t, msgInviteUser.Recipient, actual[0].String())
	expRecipient, _ := sdk.AccAddressFromBech32(msgInviteUser.Recipient)
	require.Equal(t, expRecipient, actual[0])
}

// TODO check this test
func TestMsgInviteUser_UnmarshalJson(t *testing.T) {
	//json := `{"type":"commercio/MsgInviteUser","value":{"recipient":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","sender":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"}}`
	json := `{"recipient":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","sender":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"}`

	var msg types.MsgInviteUser
	types.ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	require.Equal(t, user.String(), msg.Recipient)
	require.Equal(t, sender.String(), msg.Sender)
}

// --------------------------------
// --- MsgDepositIntoLiquidityPool
// --------------------------------

var amount = sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))
var msgDepositIntoLiquidityPool = types.NewMsgDepositIntoLiquidityPool(amount, user.String())

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
			msg:   *msgDepositIntoLiquidityPool,
			error: nil,
		},
		{
			name:  "Missing deposit returns error",
			msg:   types.MsgDepositIntoLiquidityPool{Depositor: "", Amount: amount},
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid depositor address: "),
		},
		{
			name:  "Empty deposit amount returns error",
			msg:   types.MsgDepositIntoLiquidityPool{Depositor: user.String(), Amount: nil},
			error: sdkErr.Wrap(sdkErr.ErrInvalidCoins, "Invalid deposit amount: "),
		},
		{
			name: "Negative deposit amount returns error",
			msg: types.MsgDepositIntoLiquidityPool{
				Depositor: user.String(),
				Amount:    sdk.Coins{sdk.Coin{Denom: "uatom", Amount: sdk.NewInt(-100)}},
			},
			error: sdkErr.Wrap(sdkErr.ErrInvalidCoins, "Invalid deposit amount: -100uatom"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if test.error != nil {
				require.Equal(t, test.error.Error(), test.msg.ValidateBasic().Error())
			} else {
				require.NoError(t, test.msg.ValidateBasic())
			}
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
	require.Equal(t, msgDepositIntoLiquidityPool.Depositor, actual[0].String())
}

func TestMsgDepositIntoLiquidityPool_UnmarshalJson(t *testing.T) {
	//json := `{"type":"commercio/MsgDepositIntoLiquidityPool","value":{"depositor":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","amount":[{"denom":"uatom","amount":"100"}]}}`
	json := `{"depositor":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","amount":[{"denom":"uatom","amount":"100"}]}`

	var msg types.MsgDepositIntoLiquidityPool
	types.ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	require.Equal(t, user.String(), msg.Depositor)
	require.Equal(t, amount, msg.Amount)
}

// --------------------------------
// --- MsgAddTsp
// --------------------------------

var government, _ = sdk.AccAddressFromBech32("cosmos1ct4ym78j7ksv9weyua4mzlksgwc9qq7q3wvhqg")
var tsp, _ = sdk.AccAddressFromBech32("cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee")

var msgAddTsp = types.NewMsgAddTsp(tsp.String(), government.String())

func TestMsgAddTsp_Route(t *testing.T) {
	require.Equal(t, types.RouterKey, msgAddTsp.Route())
}

func TestMsgAddTsp_Type(t *testing.T) {
	require.Equal(t, types.MsgTypeAddTsp, msgAddTsp.Type())
}

func TestMsgAddTsp_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   types.MsgAddTsp
		error error
	}{
		{
			name:  "Valid message does not return any error",
			msg:   *msgAddTsp,
			error: nil,
		},
		{
			name:  "Missing government returns error",
			msg:   types.MsgAddTsp{Government: "", Tsp: tsp.String()},
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid government address: "),
		},
		{
			name:  "Missing tsp returns error",
			msg:   types.MsgAddTsp{Government: government.String(), Tsp: ""},
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid TSP address: "),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if test.error != nil {
				require.Equal(t, test.error.Error(), test.msg.ValidateBasic().Error())
			} else {
				require.NoError(t, test.msg.ValidateBasic())
			}
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
	require.Equal(t, msgAddTsp.Government, actual[0].String())
}

func TestMsgAddTsp_UnmarshalJson(t *testing.T) {
	//json := `{"type":"commercio/MsgAddTsp","value":{"government":"cosmos1ct4ym78j7ksv9weyua4mzlksgwc9qq7q3wvhqg","tsp":"cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee"}}`
	json := `{"government":"cosmos1ct4ym78j7ksv9weyua4mzlksgwc9qq7q3wvhqg","tsp":"cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee"}`

	var msg types.MsgAddTsp
	types.ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	require.Equal(t, tsp.String(), msg.Tsp)
	require.Equal(t, government.String(), msg.Government)
}

// --------------------------------
// --- MsgRemoveTsp
// --------------------------------

var msgRemoveTsp = types.NewMsgRemoveTsp(tsp.String(), government.String())

func TestMsgRemoveTsp_Route(t *testing.T) {
	require.Equal(t, types.RouterKey, msgRemoveTsp.Route())
}

func TestMsgRemoveTsp_Type(t *testing.T) {
	require.Equal(t, types.MsgTypeRemoveTsp, msgRemoveTsp.Type())
}

func TestMsgRemoveTsp_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   types.MsgRemoveTsp
		error error
	}{
		{
			name:  "Valid message does not return any error",
			msg:   *msgRemoveTsp,
			error: nil,
		},
		{
			name:  "Missing government returns error",
			msg:   types.MsgRemoveTsp{Government: "", Tsp: tsp.String()},
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid government address: "),
		},
		{
			name:  "Missing tsp returns error",
			msg:   types.MsgRemoveTsp{Government: government.String(), Tsp: ""},
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid TSP address: "),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if test.error != nil {
				require.Equal(t, test.error.Error(), test.msg.ValidateBasic().Error())
			} else {
				require.NoError(t, test.msg.ValidateBasic())
			}
		})
	}
}

func TestMsgRemoveTsp_GetSignBytes(t *testing.T) {
	actual := msgRemoveTsp.GetSignBytes()
	expected := sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msgRemoveTsp))
	require.Equal(t, expected, actual)
}

func TestMsgRemoveTsp_GetSigners(t *testing.T) {
	actual := msgRemoveTsp.GetSigners()
	require.Equal(t, 1, len(actual))
	require.Equal(t, msgRemoveTsp.Government, actual[0].String())
}

func TestMsgRemoveTsp_UnmarshalJson(t *testing.T) {
	//json := `{"type":"commercio/MsgRemoveTsp","value":{"government":"cosmos1ct4ym78j7ksv9weyua4mzlksgwc9qq7q3wvhqg","tsp":"cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee"}}`
	json := `{"government":"cosmos1ct4ym78j7ksv9weyua4mzlksgwc9qq7q3wvhqg","tsp":"cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee"}`

	var msg types.MsgRemoveTsp
	types.ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	require.Equal(t, tsp.String(), msg.Tsp)
	require.Equal(t, government.String(), msg.Government)
}

// ---------------------------
// --- MsgBuyMembership
// ---------------------------

var testBuyer, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var testTsp, _ = sdk.AccAddressFromBech32("cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee")
var TestMembershipType = "bronze"

var membership = types.Membership{
	Owner:          testBuyer.String(),
	TspAddress:     testTsp.String(),
	MembershipType: TestMembershipType,
}

var msgBuyMembership = types.NewMsgBuyMembership(TestMembershipType, testBuyer, testTsp)

//var msgBuyMembership = types.NewMsgBuyMembership(membership)

func TestMsgBuyMembership_Route(t *testing.T) {
	require.Equal(t, types.RouterKey, msgBuyMembership.Route())
}

func TestMsgBuyMembership_Type(t *testing.T) {
	require.Equal(t, types.MsgTypeBuyMembership, msgBuyMembership.Type())
}

func TestMsgBuyMembership_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   types.MsgBuyMembership
		error error
	}{
		{
			name:  "Valid message does not return any error",
			msg:   *msgBuyMembership,
			error: nil,
		},
		{
			name:  "Missing buyer returns error",
			msg:   *types.NewMsgBuyMembership(TestMembershipType, nil, testTsp),
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid buyer address: "),
		},
		{
			name:  "Missing membership returns error",
			msg:   *types.NewMsgBuyMembership("", testBuyer, testTsp),
			error: sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid membership type: "),
		},
		{
			name:  "Invalid membership returns error",
			msg:   *types.NewMsgBuyMembership("grn", testBuyer, testTsp),
			error: sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid membership type: grn"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if test.error != nil {
				require.Equal(t, test.error.Error(), test.msg.ValidateBasic().Error())
			} else {
				require.NoError(t, test.msg.ValidateBasic())
			}
		})
	}
}

func TestMsgBuyMembership_GetSignBytes(t *testing.T) {
	//expected := `{"type":"commercio/MsgBuyMembership","value":{"buyer":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","membership_type":"bronze","tsp":"cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee"}}`
	expected := `{"buyer":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","membership_type":"bronze","tsp":"cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee"}`
	require.Equal(t, expected, string(msgBuyMembership.GetSignBytes()))
}

func TestMsgBuyMembership_GetSigners(t *testing.T) {
	//expected := []sdk.AccAddress{msgBuyMembership.Tsp}
	exp, _ := sdk.AccAddressFromBech32(msgBuyMembership.Tsp)
	expected := []sdk.AccAddress{exp}
	require.Equal(t, expected, msgBuyMembership.GetSigners())
}

// ---------------------------
// --- MsgSetMembership
// ---------------------------

var msgSetBlackmembership = types.NewMsgSetMembership(testBuyer.String(), government.String(), "black")

func TestNewMsgSetMembership_Route(t *testing.T) {
	require.Equal(t, types.RouterKey, msgSetBlackmembership.Route())
}

func TestNewMsgSetMembership_Type(t *testing.T) {
	require.Equal(t, types.MsgTypeSetMembership, msgSetBlackmembership.Type())
}

func TestNewMsgSetMembership_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   types.MsgSetMembership
		error string
	}{
		{
			name:  "Valid message does not return any error",
			msg:   *msgSetBlackmembership,
			error: "",
		},
		{
			name:  "Missing gov address returns error",
			msg:   *types.NewMsgSetMembership(testBuyer.String(), "", "black"),
			error: "Invalid government address: ",
		},
		{
			name:  "Missing subscriber returns error",
			msg:   *types.NewMsgSetMembership("", government.String(), "black"),
			error: "Invalid subscriber address: ",
		},
		{
			name:  "Missing membership returns error",
			msg:   *types.NewMsgSetMembership(testBuyer.String(), government.String(), ""),
			error: "new membership must not be empty",
		},
		/*{
			name:  "Error membership returns error",
			msg:   types.NewMsgSetMembership(testBuyer, government, "nvm"),
			error: "Invalid membership type",
		},*/
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if test.error != "" {
				require.Contains(t, test.msg.ValidateBasic().Error(), test.error)
			} else {
				require.NoError(t, test.msg.ValidateBasic())
			}
		})
	}
}

func TestNewMsgSetMembership_GetSignBytes(t *testing.T) {
	//expected := `{"type":"commercio/MsgSetMembership","value":{"government":"cosmos1ct4ym78j7ksv9weyua4mzlksgwc9qq7q3wvhqg","new_membership":"black","subscriber":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"}}`
	expected := `{"government":"cosmos1ct4ym78j7ksv9weyua4mzlksgwc9qq7q3wvhqg","new_membership":"black","subscriber":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"}`
	require.Equal(t, expected, string(msgSetBlackmembership.GetSignBytes()))
}

func TestNewMsgSetMembership_GetSigners(t *testing.T) {
	//expected := []sdk.AccAddress{msgSetBlackmembership.Government}
	exp, _ := sdk.AccAddressFromBech32(msgSetBlackmembership.Government)
	expected := []sdk.AccAddress{exp}
	require.Equal(t, expected, msgSetBlackmembership.GetSigners())
}

// ---------------------------
// --- MsgRemoveMembership
// ---------------------------

var testRemoveMembershipGov, _ = sdk.AccAddressFromBech32("cosmos1ct4ym78j7ksv9weyua4mzlksgwc9qq7q3wvhqg")
var testRemoveMembershipSubscriber, _ = sdk.AccAddressFromBech32("cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee")
var msgRemoveMembership = types.NewMsgRemoveMembership(testRemoveMembershipGov.String(), testRemoveMembershipSubscriber.String())

func TestNewMsgRemoveMembership_Route(t *testing.T) {
	require.Equal(t, types.RouterKey, msgRemoveMembership.Route())
}

func TestNewMsgRemoveMembership_Type(t *testing.T) {
	require.Equal(t, types.MsgTypeRemoveMembership, msgRemoveMembership.Type())
}

func TestNewMsgRemoveMembership_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   types.MsgRemoveMembership
		error string
	}{
		{
			name:  "Valid message does not return any error",
			msg:   *msgRemoveMembership,
			error: "",
		},
		{
			name:  "Missing gov address returns error",
			msg:   *types.NewMsgRemoveMembership("", testRemoveMembershipSubscriber.String()),
			error: "Invalid government address: ",
		},
		{
			name:  "Missing subscriber returns error",
			msg:   *types.NewMsgRemoveMembership(testRemoveMembershipGov.String(), ""),
			error: "Invalid subscriber address: ",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if test.error != "" {
				require.Contains(t, test.msg.ValidateBasic().Error(), test.error)
			} else {
				require.NoError(t, test.msg.ValidateBasic())
			}
		})
	}
}

func TestNewMsgRemoveMembership_GetSignBytes(t *testing.T) {
	//expected := `{"type":"commercio/MsgRemoveMembership","value":{"government":"cosmos1ct4ym78j7ksv9weyua4mzlksgwc9qq7q3wvhqg","subscriber":"cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee"}}`
	expected := `{"government":"cosmos1ct4ym78j7ksv9weyua4mzlksgwc9qq7q3wvhqg","subscriber":"cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee"}`
	require.Equal(t, expected, string(msgRemoveMembership.GetSignBytes()))
}

func TestNewMsgRemoveMembership_GetSigners(t *testing.T) {
	//expected := []sdk.AccAddress{msgRemoveMembership.Government}
	exp, _ := sdk.AccAddressFromBech32(msgSetBlackmembership.Government)
	expected := []sdk.AccAddress{exp}
	require.Equal(t, expected, msgRemoveMembership.GetSigners())
}
