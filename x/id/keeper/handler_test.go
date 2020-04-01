package keeper

import (
	"errors"
	"fmt"
	"testing"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/id/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestValidMsg_StoreDoc(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()

	handler := NewHandler(k, govK)
	msgSetIdentity := types.MsgSetIdentity(TestDidDocument)
	_, err := handler(ctx, msgSetIdentity)

	require.NoError(t, err)
}

func TestInvalidMsg(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()

	tm := sdk.NewTestMsg()
	handler := NewHandler(k, govK)
	_, err := handler(ctx, tm)

	require.Error(t, err)
	require.Equal(t, fmt.Sprintf("unknown request: Unrecognized %s message type: %s", types.ModuleName, tm.Type()), err.Error())
}

// ----------------------------
// --- Did power up requests
// --------------------------

var msgRequestDidPowerUp = types.MsgRequestDidPowerUp{
	Claimant: TestDidPowerUpRequest.Claimant,
	Amount:   TestDidPowerUpRequest.Amount,
	Proof:    TestDidPowerUpRequest.Proof,
	ID:       TestDidPowerUpRequest.ID,
	ProofKey: TestDidPowerUpRequest.ProofKey,
}

func Test_handleMsgRequestDidPowerUp_NewRequest(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()

	handler := NewHandler(k, govK)
	_, err := handler(ctx, msgRequestDidPowerUp)
	require.NoError(t, err)

	stored, err := k.GetPowerUpRequestByID(ctx, TestDidPowerUpRequest.ID)
	require.NoError(t, err)
	require.Equal(t, TestDidPowerUpRequest.Proof, stored.Proof)
	require.Equal(t, TestDidPowerUpRequest.Amount.String(), stored.Amount.String())
	require.Equal(t, TestDidPowerUpRequest.ID, stored.ID)
	require.Equal(t, TestDidPowerUpRequest.Claimant.String(), stored.Claimant.String())
	require.Equal(t, TestDidPowerUpRequest.ProofKey, stored.ProofKey)
}

func Test_handleMsgRequestDidPowerUp_ExistingRequest(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	handler := NewHandler(k, govK)
	_, err := handler(ctx, msgRequestDidPowerUp)

	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrUnknownRequest))
}

func Test_handleMsgPowerUpDid_InvalidTumbler(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()

	msg := types.MsgChangePowerUpStatus{
		Recipient: TestDidPowerUpRequest.Claimant,
		Signer:    TestDidPowerUpRequest.Claimant,
	}
	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)

	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrInvalidAddress))
	require.Contains(t, err.Error(), "tumbler")
}
