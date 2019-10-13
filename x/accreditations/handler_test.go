package accreditations

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// --------------------------
// --- handleMsgInviteUser
// --------------------------

func Test_handleMsgInviteUser_ExistingInvite(t *testing.T) {
	ctx, _, _, _, govK, k := GetTestInput()

	invite := Invite{Sender: TestInviteSender, User: TestUser, Rewarded: false}
	k.SaveInvite(ctx, invite)

	handler := NewHandler(k, govK)
	msg := NewMsgInviteUser(TestUser, TestInviteSender)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
	assert.Contains(t, res.Log, "has already been invited")
}

func Test_handleMsgInviteUser_NewInvite(t *testing.T) {
	ctx, _, _, _, govK, k := GetTestInput()

	handler := NewHandler(k, govK)
	msg := NewMsgInviteUser(TestUser, TestInviteSender)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())
}

// -----------------------------
// --- handleMsgSetUserVerified
// -----------------------------

func Test_handleMsgSetUserVerified_InvalidSigner(t *testing.T) {
	ctx, _, _, _, govK, k := GetTestInput()

	handler := NewHandler(k, govK)
	msg := NewMsgSetUserVerified(TestUser, TestTimestamp, nil)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnauthorized, res.Code)
}

func Test_handleMsgSetUserVerified_ExistingCredential(t *testing.T) {
	ctx, _, _, _, govK, k := GetTestInput()

	credential := types.Credential{User: TestUser, Verifier: TestTsp, Timestamp: TestTimestamp}
	k.SaveCredential(ctx, credential)

	handler := NewHandler(k, govK)
	msg := NewMsgSetUserVerified(TestUser, TestTimestamp, TestTsp)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
}

func Test_handleMsgSetUserVerified_NewCredential(t *testing.T) {
	ctx, _, _, _, govK, k := GetTestInput()

	handler := NewHandler(k, govK)
	msg := NewMsgSetUserVerified(TestUser, TestTimestamp, TestTsp)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
}

// -----------------------------
// --- handleMsgAddTrustedSigner
// -----------------------------

func Test_handleAddTrustedSigner_InvalidGovernment(t *testing.T) {
	ctx, _, _, _, governmentKeeper, accreditationKeeper := GetTestInput()

	government, _ := sdk.AccAddressFromBech32("cosmos15ne6fy8uukkyyf072qklkeleh2zf39k52mcg2f")

	err := governmentKeeper.SetGovernmentAddress(ctx, TestUser)
	assert.Nil(t, err)

	handler := NewHandler(accreditationKeeper, governmentKeeper)
	msg := MsgAddTrustedSigner{Government: government, Tsp: TestTsp}
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
}

func Test_handleAddTrustedSigner_ValidGovernment(t *testing.T) {
	ctx, _, _, _, governmentKeeper, accreditationKeeper := GetTestInput()

	government, _ := sdk.AccAddressFromBech32("cosmos15ne6fy8uukkyyf072qklkeleh2zf39k52mcg2f")

	err := governmentKeeper.SetGovernmentAddress(ctx, government)
	assert.Nil(t, err)

	handler := NewHandler(accreditationKeeper, governmentKeeper)
	msg := MsgAddTrustedSigner{Government: government, Tsp: TestTsp}
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())
	assert.Equal(t, sdk.CodeOK, res.Code)
}
