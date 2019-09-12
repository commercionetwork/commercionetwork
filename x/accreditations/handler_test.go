package accreditations

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// --------------------------
// --- handleSetAccrediter
// --------------------------

func Test_handleSetAccrediter_NotTrustedAccrediter(t *testing.T) {
	ctx, _, _, _, accreditationKeeper := keeper.GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Delete([]byte(types.TrustedSignersStoreKey))

	handler := NewHandler(accreditationKeeper)
	msg := MsgSetAccrediter{User: keeper.TestUser, Accrediter: keeper.TestAccrediter, Signer: keeper.TestSigner}
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
}

func Test_handleSetAccrediter_ExistingAccrediter(t *testing.T) {
	ctx, cdc, _, _, accreditationKeeper := keeper.GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Set([]byte(types.TrustedSignersStoreKey), cdc.MustMarshalBinaryBare(&[]sdk.AccAddress{keeper.TestSigner}))

	accreditation := types.Accreditation{Accrediter: keeper.TestAccrediter, User: keeper.TestUser, Rewarded: false}
	store.Set(keeper.TestUser, cdc.MustMarshalBinaryBare(accreditation))

	handler := NewHandler(accreditationKeeper)
	msg := MsgSetAccrediter{User: keeper.TestUser, Accrediter: keeper.TestAccrediter, Signer: keeper.TestSigner}
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
}

func Test_handleSetAccrediter_ValidData(t *testing.T) {
	ctx, cdc, _, _, accreditationKeeper := keeper.GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Set([]byte(types.TrustedSignersStoreKey), cdc.MustMarshalBinaryBare(&[]sdk.AccAddress{keeper.TestSigner}))

	handler := NewHandler(accreditationKeeper)
	msg := MsgSetAccrediter{User: keeper.TestUser, Accrediter: keeper.TestAccrediter, Signer: keeper.TestSigner}
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())
	assert.Equal(t, sdk.CodeOK, res.Code)
}

// -----------------------------
// --- handleDistributeReward
// -----------------------------

func Test_handleDistributeReward_NilAccrediter(t *testing.T) {
	ctx, _, _, _, accreditationKeeper := keeper.GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Delete(keeper.TestUser)

	handler := NewHandler(accreditationKeeper)
	reward := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))
	msg := MsgDistributeReward{User: keeper.TestUser, Accrediter: keeper.TestAccrediter, Reward: reward, Signer: keeper.TestSigner}
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
}

func Test_handleDistributeReward_InvalidAccrediter(t *testing.T) {
	ctx, cdc, _, _, accreditationKeeper := keeper.GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)

	reward := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))

	accreditation := types.Accreditation{Accrediter: keeper.TestUser, User: keeper.TestUser, Rewarded: false}
	store.Set(keeper.TestUser, cdc.MustMarshalBinaryBare(accreditation))

	handler := NewHandler(accreditationKeeper)
	msg := MsgDistributeReward{User: keeper.TestUser, Accrediter: keeper.TestAccrediter, Reward: reward, Signer: keeper.TestSigner}
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
}

func Test_handleDistributeReward_RewardedAccrediter(t *testing.T) {
	ctx, cdc, _, _, accreditationKeeper := keeper.GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)

	reward := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))

	accreditation := types.Accreditation{Accrediter: keeper.TestAccrediter, User: keeper.TestUser, Rewarded: true}
	store.Set(keeper.TestUser, cdc.MustMarshalBinaryBare(accreditation))

	handler := NewHandler(accreditationKeeper)
	msg := MsgDistributeReward{User: keeper.TestUser, Accrediter: keeper.TestAccrediter, Reward: reward, Signer: keeper.TestSigner}
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
}

func Test_handleDistributeReward_ValidData(t *testing.T) {
	ctx, cdc, _, _, accreditationKeeper := keeper.GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)

	pool := []sdk.Coin{sdk.NewCoin("uatom", sdk.NewInt(100))}
	store.Set([]byte(types.LiquidityPoolKey), cdc.MustMarshalBinaryBare(&pool))

	reward := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))

	accreditation := types.Accreditation{Accrediter: keeper.TestAccrediter, User: keeper.TestUser, Rewarded: false}
	store.Set(keeper.TestUser, cdc.MustMarshalBinaryBare(accreditation))

	handler := NewHandler(accreditationKeeper)
	msg := MsgDistributeReward{User: keeper.TestUser, Accrediter: keeper.TestAccrediter, Reward: reward, Signer: keeper.TestSigner}
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())
	assert.Equal(t, sdk.CodeOK, res.Code)
}
