package types_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/id/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var msgSetIdentity types.MsgSetIdentity
var msgRequestDidPowerUp types.MsgRequestDidPowerUp
var msgChangePowerUpStatus types.MsgChangePowerUpStatus
var requestRecipient sdk.AccAddress
var requestSender sdk.AccAddress

func init() {
	types.ConfigTestPrefixes()

	var testZone, _ = time.LoadLocation("UTC")
	var testTime = time.Date(2016, 2, 8, 16, 2, 20, 0, testZone)
	var testOwnerAddress, _ = sdk.AccAddressFromBech32("did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf")
	msgSetIdentity = types.NewMsgSetIdentity(types.DidDocument{
		Context: types.ContextDidV1,
		ID:      testOwnerAddress,
		Proof: types.Proof{
			Type:               types.KeyTypeSecp256k12019,
			Created:            testTime,
			ProofPurpose:       types.ProofPurposeAuthentication,
			Controller:         testOwnerAddress.String(),
			SignatureValue:     "4T2jhs4C0k7p649tdzQAOLqJ0GJsiFDP/NnsSkFpoXAxcgn6h/EgvOpHxW7FMNQ9RDgQbcE6FWP6I2UsNv1qXQ==",
			VerificationMethod: "did:com:pub1addwnpepqwzc44ggn40xpwkfhcje9y7wdz6sunuv2uydxmqjrvcwff6npp2exy5dn6c",
		},
		PubKeys: types.PubKeys{
			types.PubKey{
				ID:         fmt.Sprintf("%s#keys-1", testOwnerAddress),
				Type:       "RsaVerificationKey2018",
				Controller: testOwnerAddress,
				PublicKey: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqOoLR843vgkFGudQsjch
2K85QJ4Hh7l2jjrMesQFDWVcW1xr//eieGzxDogWx7tMOtQ0hw77NAURhldek1Bh
Co06790YHAE97JqgRQ+IR9Dl3GaGVQ2WcnknO4B1cvTRJmdsqrN1Bs4Qfd+jjKIM
V1tz8zU9NmdR+DvGkAYYxoIx74YaTAxH+GCArfWMG1tRJPI9MELZbOWd9xkKlPic
bLp8coZh9NgLajMDWKXpuHQ8cdJSxQ/ekZaTuEy7qbjbGBMVzbjhPjcxffQmGV1W
gNY1BGplZz9mbBmH7siKnKIVZ5Bp55uLfEw+u2yOVx/0yKUdsmZoe4jhevCSq3aw
GwIDAQAB
-----END PUBLIC KEY-----`,
			},
			types.PubKey{
				ID:         fmt.Sprintf("%s#keys-2", testOwnerAddress),
				Type:       "RsaSignatureKey2018",
				Controller: testOwnerAddress,
				PublicKey: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+Juw6xqYchTNFYUznmoB
CzKfQG75v2Pv1Db1Z5EJgP6i0yRsBG1VqIOY4icRnyhDDVFi1omQjjUuCRxWGjsc
B1UkSnybm0WC+g82HL3mUzbZja27NFJPuNaMaUlNbe0daOG88FS67jq5J2LsZH/V
cGZBX5bbtCe0Niq39mQdJxdHq3D5ROMA73qeYvLkmXS6Dvs0w0fHsy+DwJtdOnOj
xt4F5hIEXGP53qz2tBjCRL6HiMP/cLSwAd7oc67abgQxfnf9qldyd3X0IABpti1L
irJNugfN6HuxHDm6dlXVReOhHRbkEcWedv82Ji5d/sDZ+WT+yWILOq03EJo/LXJ1
SQIDAQAB
-----END PUBLIC KEY-----`,
			},
		},
	})

	requestSender, _ = sdk.AccAddressFromBech32("did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf")
	requestRecipient, _ = sdk.AccAddressFromBech32("did:com:1sqnp7cmasyv2yathd8ye8xlhhaqaw953sc5lp6")

	// --------------------------
	// --- MsgRequestDidPowerUp
	// --------------------------

	msgRequestDidPowerUp = types.MsgRequestDidPowerUp{
		Claimant: requestSender,
		Amount:   sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		ID:       "47F8F05F-AA3C-4E2B-9944-34EB9FB24BAE",
		Proof:    "68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
		ProofKey: "eW91IGxvc3QgdGhlIGdhbWUK",
	}

	msgChangePowerUpStatus = types.MsgChangePowerUpStatus{
		Recipient: requestRecipient,
		Status: types.RequestStatus{
			Type:    types.StatusApproved,
			Message: "",
		},
		PowerUpID: "47F8F05F-AA3C-4E2B-9944-34EB9FB24BAE",
		Signer:    requestSender,
	}
}

// ----------------------------------
// --- SetIdentity
// ----------------------------------

func TestMsgSetIdentity_Route(t *testing.T) {
	require.Equal(t, types.ModuleName, msgSetIdentity.Route())
}

func TestMsgSetIdentity_Type(t *testing.T) {
	require.Equal(t, types.MsgTypeSetIdentity, msgSetIdentity.Type())
}

func TestMsgSetIdentity_ValidateBasic(t *testing.T) {
	testData := []struct {
		name          string
		message       types.MsgSetIdentity
		shouldBeValid bool
	}{
		{
			name:          "Valid message",
			message:       msgSetIdentity,
			shouldBeValid: true,
		},
		{
			name:          "Invalid address",
			message:       types.MsgSetIdentity{},
			shouldBeValid: false,
		},
	}

	for _, test := range testData {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if test.shouldBeValid {
				require.NoError(t, test.message.ValidateBasic())
			} else {
				require.Error(t, test.message.ValidateBasic())
			}
		})
	}
}

func TestMsgSetIdentity_GetSignBytes(t *testing.T) {
	var m struct {
		Type  string               `json:"type"`
		Value types.MsgSetIdentity `json:"value"`
	}
	cdc := codec.New()
	err := cdc.UnmarshalJSON(msgSetIdentity.GetSignBytes(), &m)
	require.NoError(t, err)
	require.Equal(t, msgSetIdentity, m.Value)
}

func TestMsgSetIdentity_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgSetIdentity.ID}
	require.Equal(t, expected, msgSetIdentity.GetSigners())
}

func TestMsgRequestDidPowerUp_Route(t *testing.T) {
	require.Equal(t, types.ModuleName, msgRequestDidPowerUp.Route())
}

func TestMsgRequestDidPowerUp_Type(t *testing.T) {
	require.Equal(t, types.MsgTypeRequestDidPowerUp, msgRequestDidPowerUp.Type())
}

func TestMsgRequestDidPowerUp_ValidateBasic(t *testing.T) {
	claimant, _ := sdk.AccAddressFromBech32("did:com:1zla8arsc5rju9wekz00yz54zguj20a96jn9cy6")
	amount := sdk.NewCoins(sdk.NewInt64Coin("uatom", 100))

	tests := []struct {
		name  string
		msg   types.MsgRequestDidPowerUp
		error error
	}{
		{
			name:  "Invalid claimant returns error",
			msg:   types.MsgRequestDidPowerUp{Claimant: sdk.AccAddress{}},
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid claimant: "),
		},
		{
			name:  "Invalid amount returns error",
			msg:   types.MsgRequestDidPowerUp{Claimant: claimant, Amount: sdk.NewCoins()},
			error: sdkErr.Wrap(sdkErr.ErrInvalidCoins, "Power up amount not valid: "),
		},
		{
			name: "Valid message returns no error",
			msg:  msgRequestDidPowerUp,
		},
		{
			name:  "Invalid proof returns error",
			msg:   types.MsgRequestDidPowerUp{Claimant: claimant, Amount: amount, Proof: "230sd"},
			error: sdkErr.Wrap(sdkErr.ErrUnknownRequest, "proof must be base64-encoded"),
		},
		{
			name: "Invalid encryption key returns error",
			msg: types.MsgRequestDidPowerUp{
				Claimant: claimant,
				Amount:   amount,
				Proof:    msgRequestDidPowerUp.Proof,
				ID:       msgRequestDidPowerUp.ID,
				ProofKey: "1230xcv",
			},
			error: sdkErr.Wrap(sdkErr.ErrUnknownRequest, "proof key must be base64-encoded"),
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

func TestMsgRequestDidPowerUp_GetSignBytes(t *testing.T) {
	var m struct {
		Type  string                     `json:"type"`
		Value types.MsgRequestDidPowerUp `json:"value"`
	}
	cdc := codec.New()
	err := cdc.UnmarshalJSON(msgRequestDidPowerUp.GetSignBytes(), &m)
	require.NoError(t, err)
	require.Equal(t, msgRequestDidPowerUp, m.Value)
}

func TestMsgRequestDidPowerUp_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgRequestDidPowerUp.Claimant}
	require.Equal(t, expected, msgRequestDidPowerUp.GetSigners())
}

// ------------------------
// --- MsgChangePowerUpStatus
// ------------------------

func TestMsgChangePowerUpStatus_Route(t *testing.T) {
	require.Equal(t, types.ModuleName, msgChangePowerUpStatus.Route())
}

func TestMsgChangePowerUpStatus_Type(t *testing.T) {
	require.Equal(t, types.MsgTypeChangePowerUpStatus, msgChangePowerUpStatus.Type())
}

func TestMsgChangePowerUpStatus_ValidateBasic(t *testing.T) {
	claimant, _ := sdk.AccAddressFromBech32("did:com:1zla8arsc5rju9wekz00yz54zguj20a96jn9cy6")

	tests := []struct {
		name  string
		msg   types.MsgChangePowerUpStatus
		error error
	}{
		{
			name:  "Invalid recipient returns error",
			msg:   types.MsgChangePowerUpStatus{Recipient: sdk.AccAddress{}},
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid recipient address: "),
		},
		{
			name: "Invalid Status returns error",
			msg: types.MsgChangePowerUpStatus{
				Recipient: msgChangePowerUpStatus.Recipient,
				PowerUpID: msgChangePowerUpStatus.PowerUpID,
				Signer:    msgChangePowerUpStatus.Signer,
				Status:    types.RequestStatus{},
			},
			error: sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid status type: "),
		},
		{
			name: "Valid message returns no error",
			msg:  msgChangePowerUpStatus,
		},
		{
			name: "Invalid PowerUp ID returns error",
			msg: types.MsgChangePowerUpStatus{
				Recipient: msgChangePowerUpStatus.Recipient,
				PowerUpID: "oh no",
				Signer:    msgChangePowerUpStatus.Signer,
				Status:    msgChangePowerUpStatus.Status,
			},
			error: sdkErr.Wrap(sdkErr.ErrUnauthorized, "invalid PowerUpID, must be a valid UUID"),
		},
		{
			name:  "Invalid signer returns error",
			msg:   types.MsgChangePowerUpStatus{Recipient: claimant, Status: msgChangePowerUpStatus.Status, PowerUpID: msgChangePowerUpStatus.PowerUpID, Signer: sdk.AccAddress{}},
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid signer address: "),
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

func TestMsgChangePowerUpStatus_GetSignBytes(t *testing.T) {
	var m struct {
		Type  string                       `json:"type"`
		Value types.MsgChangePowerUpStatus `json:"value"`
	}
	cdc := codec.New()
	err := cdc.UnmarshalJSON(msgChangePowerUpStatus.GetSignBytes(), &m)
	require.NoError(t, err)
	require.Equal(t, msgChangePowerUpStatus, m.Value)
}

func TestMsgChangePowerUpStatus_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgChangePowerUpStatus.Signer}
	require.Equal(t, expected, msgChangePowerUpStatus.GetSigners())
}
