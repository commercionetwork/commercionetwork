package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// Test vars
var user, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var accrediter, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var signer, _ = sdk.AccAddressFromBech32("cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee")

// ----------------------
// --- MsgSetAccrediter
// ----------------------

var msgSetAccrediter = MsgSetAccrediter{
	User:       user,
	Accrediter: accrediter,
	Signer:     signer,
}

func TestMsgSetAccrediter_Route(t *testing.T) {
	assert.Equal(t, QuerierRoute, msgSetAccrediter.Route())
}

func TestMsgSetAccrediter_Type(t *testing.T) {
	assert.Equal(t, MsgTypeSetAccrediter, msgSetAccrediter.Type())
}

func TestMsgSetAccrediter_ValidateBasic_ValidMsg(t *testing.T) {
	assert.Nil(t, msgSetAccrediter.ValidateBasic())
}

func TestMsgSetAccrediter_ValidateBasic_MissingUser(t *testing.T) {
	msg := MsgSetAccrediter{User: nil, Accrediter: accrediter, Signer: signer}
	assert.NotNil(t, msg.ValidateBasic())
}

func TestMsgSetAccrediter_ValidateBasic_MissingAccrediter(t *testing.T) {
	msg := MsgSetAccrediter{User: user, Accrediter: nil, Signer: signer}
	assert.NotNil(t, msg.ValidateBasic())
}

func TestMsgSetAccrediter_GetSignBytes(t *testing.T) {
	actual := msgSetAccrediter.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSetAccrediter))
	assert.Equal(t, expected, actual)
}

func TestMsgSetAccrediter_GetSigners(t *testing.T) {
	actual := msgSetAccrediter.GetSigners()
	assert.Equal(t, 1, len(actual))
	assert.Equal(t, msgSetAccrediter.Signer, actual[0])
}

func TestMsgSetAccrediter_UnmarshalJson(t *testing.T) {
	json := `{"type":"commercio/MsgSetAccrediter","value":{"user":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","accrediter":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","signer":"cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee"}}`

	var msg MsgSetAccrediter
	ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	assert.Equal(t, user, msg.User)
	assert.Equal(t, accrediter, msg.Accrediter)
	assert.Equal(t, signer, msg.Signer)
}

// --------------------------
// --- MsgDistributeReward
// --------------------------

var reward = sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))
var msgDistributeReward = MsgDistributeReward{
	User:       user,
	Accrediter: accrediter,
	Signer:     signer,
	Reward:     reward,
}

func TestMsgDistributeReward_Route(t *testing.T) {
	assert.Equal(t, QuerierRoute, msgDistributeReward.Route())
}

func TestMsgDistributeReward_Type(t *testing.T) {
	assert.Equal(t, MsgTypeDistributeReward, msgDistributeReward.Type())
}

func TestMsgDistributeReward_ValidateBasic_ValidMsg(t *testing.T) {
	assert.Nil(t, msgDistributeReward.ValidateBasic())
}

func TestMsgDistributeReward_ValidateBasic_MissingAccrediter(t *testing.T) {
	msg := MsgDistributeReward{Accrediter: nil, User: user, Signer: signer, Reward: reward}
	assert.NotNil(t, msg.ValidateBasic())
}

func TestMsgDistributeReward_ValidateBasic_MissingUser(t *testing.T) {
	msg := MsgDistributeReward{Accrediter: accrediter, User: nil, Signer: signer, Reward: reward}
	assert.NotNil(t, msg.ValidateBasic())
}

func TestMsgDistributeReward_ValidateBasic_MissingSinger(t *testing.T) {
	msg := MsgDistributeReward{Accrediter: accrediter, User: user, Signer: nil, Reward: reward}
	assert.NotNil(t, msg.ValidateBasic())
}

func TestMsgDistributeReward_ValidateBasic_MissingReward(t *testing.T) {
	msg := MsgDistributeReward{Accrediter: accrediter, User: user, Signer: signer, Reward: nil}
	assert.NotNil(t, msg.ValidateBasic())
}

func TestMsgDistributeReward_ValidateBasic_NegativeReward(t *testing.T) {
	reward := sdk.Coins{sdk.Coin{Denom: "uatom", Amount: sdk.NewInt(-100)}}
	msg := MsgDistributeReward{Accrediter: accrediter, User: user, Signer: signer, Reward: reward}
	assert.NotNil(t, msg.ValidateBasic())
}

func TestMsgDistributeReward_GetSignBytes(t *testing.T) {
	actual := msgDistributeReward.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgDistributeReward))
	assert.Equal(t, expected, actual)
}

func TestMsgDistributeReward_GetSigners(t *testing.T) {
	actual := msgDistributeReward.GetSigners()
	assert.Equal(t, 1, len(actual))
	assert.Equal(t, msgDistributeReward.Signer, actual[0])
}

func TestMsgDistributeReward_UnmarshalJson(t *testing.T) {
	json := `{"type":"commercio/MsgDistributeReward","value":{"user":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","accrediter":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","signer":"cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee","reward":[{"denom":"uatom","amount":"100"}]}}`

	var msg MsgDistributeReward
	ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	assert.Equal(t, user, msg.User)
	assert.Equal(t, accrediter, msg.Accrediter)
	assert.Equal(t, signer, msg.Signer)
	assert.Equal(t, reward, msg.Reward)
}

// --------------------------------
// --- MsgDepositIntoLiquidityPool
// --------------------------------

var amount = sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))
var msgDepositIntoLiquidityPool = MsgDepositIntoLiquidityPool{
	Depositor: user,
	Amount:    amount,
}

func TestMsgDepositIntoLiquidityPool_Route(t *testing.T) {
	assert.Equal(t, QuerierRoute, msgDepositIntoLiquidityPool.Route())
}

func TestMsgDepositIntoLiquidityPool_Type(t *testing.T) {
	assert.Equal(t, MsgTypesDepositIntoLiquidityPool, msgDepositIntoLiquidityPool.Type())
}

func TestMsgDepositIntoLiquidityPool_ValidateBasic_ValidMsg(t *testing.T) {
	assert.Nil(t, msgDepositIntoLiquidityPool.ValidateBasic())
}

func TestMsgDepositIntoLiquidityPool_ValidateBasic_MissingDepositor(t *testing.T) {
	msg := MsgDepositIntoLiquidityPool{Depositor: nil, Amount: amount}
	assert.NotNil(t, msg.ValidateBasic())
}

func TestMsgDepositIntoLiquidityPool_ValidateBasic_MissingAmount(t *testing.T) {
	msg := MsgDepositIntoLiquidityPool{Depositor: user, Amount: nil}
	assert.NotNil(t, msg.ValidateBasic())
}

func TestMsgDepositIntoLiquidityPool_ValidateBasic_NegativeAmount(t *testing.T) {
	amount := sdk.Coins{sdk.Coin{Denom: "uatom", Amount: sdk.NewInt(-100)}}
	msg := MsgDepositIntoLiquidityPool{Depositor: user, Amount: amount}
	assert.NotNil(t, msg.ValidateBasic())
}

func TestMsgDepositIntoLiquidityPool_GetSignBytes(t *testing.T) {
	actual := msgDepositIntoLiquidityPool.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgDepositIntoLiquidityPool))
	assert.Equal(t, expected, actual)
}

func TestMsgDepositIntoLiquidityPool_GetSigners(t *testing.T) {
	actual := msgDepositIntoLiquidityPool.GetSigners()
	assert.Equal(t, 1, len(actual))
	assert.Equal(t, msgDepositIntoLiquidityPool.Depositor, actual[0])
}

func TestMsgDepositIntoLiquidityPool_UnmarshalJson(t *testing.T) {
	json := `{"type":"commercio/MsgDepositIntoLiquidityPool","value":{"depositor":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","amount":[{"denom":"uatom","amount":"100"}]}}`

	var msg MsgDepositIntoLiquidityPool
	ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	assert.Equal(t, user, msg.Depositor)
	assert.Equal(t, amount, msg.Amount)
}

// --------------------------------
// --- MsgAddTrustedSigner
// --------------------------------

var government, _ = sdk.AccAddressFromBech32("cosmos1ct4ym78j7ksv9weyua4mzlksgwc9qq7q3wvhqg")
var msgAddTrustedSigner = MsgAddTrustedSigner{
	Government:    government,
	TrustedSigner: signer,
}

func TestMsgAddTrustedSigner_Route(t *testing.T) {
	assert.Equal(t, QuerierRoute, msgAddTrustedSigner.Route())
}

func TestMsgAddTrustedSigner_Type(t *testing.T) {
	assert.Equal(t, MsgTypeAddTrustedSigner, msgAddTrustedSigner.Type())
}

func TestMsgAddTrustedSigner_ValidateBasic_ValidMsg(t *testing.T) {
	assert.Nil(t, msgAddTrustedSigner.ValidateBasic())
}

func TestMsgAddTrustedSigner_ValidateBasic_MissingGovernment(t *testing.T) {
	msg := MsgAddTrustedSigner{Government: nil, TrustedSigner: signer}
	assert.NotNil(t, msg.ValidateBasic())
}

func TestMsgAddTrustedSigner_ValidateBasic_MissingSigner(t *testing.T) {
	msg := MsgAddTrustedSigner{Government: government, TrustedSigner: nil}
	assert.NotNil(t, msg.ValidateBasic())
}

func TestMsgAddTrustedSigner_GetSignBytes(t *testing.T) {
	actual := msgAddTrustedSigner.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgAddTrustedSigner))
	assert.Equal(t, expected, actual)
}

func TestMsgAddTrustedSigner_GetSigners(t *testing.T) {
	actual := msgAddTrustedSigner.GetSigners()
	assert.Equal(t, 1, len(actual))
	assert.Equal(t, msgAddTrustedSigner.Government, actual[0])
}

func TestMsgAddTrustedSigner_UnmarshalJson(t *testing.T) {
	json := `{"type":"commercio/MsgAddTrustedSigner","value":{"government":"cosmos1ct4ym78j7ksv9weyua4mzlksgwc9qq7q3wvhqg","signer":"cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee"}}`

	var msg MsgAddTrustedSigner
	ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	assert.Equal(t, signer, msg.TrustedSigner)
	assert.Equal(t, government, msg.Government)
}
