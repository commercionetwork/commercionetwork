package accreditations

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// --------------------------
// --- handleSetAccrediter
// --------------------------

func Test_handleSetAccrediter_NotTrustedAccrediter(t *testing.T) {
	ctx, _, _, _, governmentKeeper, accreditationKeeper := GetTestInput()

	handler := NewHandler(accreditationKeeper, governmentKeeper)
	msg := MsgSetAccrediter{User: TestUser, Accrediter: TestAccrediter, Signer: TestSigner}
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
}

func Test_handleSetAccrediter_ExistingAccrediter(t *testing.T) {
	ctx, cdc, _, _, governmentKeeper, accreditationKeeper := GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Set([]byte(TrustedSignersStoreKey), cdc.MustMarshalBinaryBare(&[]sdk.AccAddress{TestSigner}))

	accreditation := Accreditation{Accrediter: TestAccrediter, User: TestUser, Rewarded: false}
	store.Set(TestUser, cdc.MustMarshalBinaryBare(accreditation))

	handler := NewHandler(accreditationKeeper, governmentKeeper)
	msg := MsgSetAccrediter{User: TestUser, Accrediter: TestAccrediter, Signer: TestSigner}
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
}

func Test_handleSetAccrediter_ValidData(t *testing.T) {
	ctx, cdc, _, _, governmentKeeper, accreditationKeeper := GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Set([]byte(TrustedSignersStoreKey), cdc.MustMarshalBinaryBare(&[]sdk.AccAddress{TestSigner}))

	handler := NewHandler(accreditationKeeper, governmentKeeper)
	msg := MsgSetAccrediter{User: TestUser, Accrediter: TestAccrediter, Signer: TestSigner}
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())
	assert.Equal(t, sdk.CodeOK, res.Code)
}

// -----------------------------
// --- handleDistributeReward
// -----------------------------

func Test_handleDistributeReward_NilAccrediter(t *testing.T) {
	ctx, _, _, _, governmentKeeper, accreditationKeeper := GetTestInput()

	handler := NewHandler(accreditationKeeper, governmentKeeper)
	reward := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))
	msg := MsgDistributeReward{User: TestUser, Accrediter: TestAccrediter, Reward: reward, Signer: TestSigner}
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
}

func Test_handleDistributeReward_InvalidAccrediter(t *testing.T) {
	ctx, cdc, _, _, governmentKeeper, accreditationKeeper := GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)

	reward := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))

	accreditation := Accreditation{Accrediter: TestUser, User: TestUser, Rewarded: false}
	store.Set(TestUser, cdc.MustMarshalBinaryBare(accreditation))

	handler := NewHandler(accreditationKeeper, governmentKeeper)
	msg := MsgDistributeReward{User: TestUser, Accrediter: TestAccrediter, Reward: reward, Signer: TestSigner}
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
}

func Test_handleDistributeReward_RewardedAccrediter(t *testing.T) {
	ctx, cdc, _, _, governmentKeeper, accreditationKeeper := GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)

	reward := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))

	accreditation := Accreditation{Accrediter: TestAccrediter, User: TestUser, Rewarded: true}
	store.Set(TestUser, cdc.MustMarshalBinaryBare(accreditation))

	handler := NewHandler(accreditationKeeper, governmentKeeper)
	msg := MsgDistributeReward{User: TestUser, Accrediter: TestAccrediter, Reward: reward, Signer: TestSigner}
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
}

func Test_handleDistributeReward_ValidData(t *testing.T) {
	ctx, cdc, _, _, governmentKeeper, accreditationKeeper := GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)

	pool := []sdk.Coin{sdk.NewCoin("uatom", sdk.NewInt(100))}
	store.Set([]byte(LiquidityPoolKey), cdc.MustMarshalBinaryBare(&pool))

	reward := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))

	accreditation := Accreditation{Accrediter: TestAccrediter, User: TestUser, Rewarded: false}
	store.Set(TestUser, cdc.MustMarshalBinaryBare(accreditation))

	handler := NewHandler(accreditationKeeper, governmentKeeper)
	msg := MsgDistributeReward{User: TestUser, Accrediter: TestAccrediter, Reward: reward, Signer: TestSigner}
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())
	assert.Equal(t, sdk.CodeOK, res.Code)
}

// -----------------------------
// --- handleAddTrustedSigner
// -----------------------------

func Test_handleAddTrustedSigner_InvalidGovernment(t *testing.T) {
	ctx, _, _, _, governmentKeeper, accreditationKeeper := GetTestInput()

	government, _ := sdk.AccAddressFromBech32("cosmos15ne6fy8uukkyyf072qklkeleh2zf39k52mcg2f")

	err := governmentKeeper.SetGovernmentAddress(ctx, TestUser)
	assert.Nil(t, err)

	handler := NewHandler(accreditationKeeper, governmentKeeper)
	msg := MsgAddTrustedSigner{Government: government, TrustedSigner: TestSigner}
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
	msg := MsgAddTrustedSigner{Government: government, TrustedSigner: TestSigner}
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())
	assert.Equal(t, sdk.CodeOK, res.Code)
}
