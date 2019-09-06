package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// --------------------------
// --- MsgSetAccrediter
// --------------------------

// MsgSetAccrediter should be used when wanting to set a specific accrediter
// for a specific user.
// Note that the associated signer should be a trustworthy one in order to avoid
// unauthorized users to perform such assignment.
type MsgSetAccrediter struct {
	User       sdk.AccAddress `json:"user"`
	Accrediter sdk.AccAddress `json:"accrediter"`
	Signer     sdk.AccAddress `json:"signer"`
}

func NewMsgSetAccrediter(user, accrediter, signer sdk.AccAddress) MsgSetAccrediter {
	return MsgSetAccrediter{
		User:       user,
		Accrediter: accrediter,
		Signer:     signer,
	}
}

// RouterKey Implements Msg.
func (msg MsgSetAccrediter) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSetAccrediter) Type() string { return MsgTypeSetAccrediter }

// ValidateBasic Implements Msg.
func (msg MsgSetAccrediter) ValidateBasic() sdk.Error {
	if msg.User.Empty() {
		return sdk.ErrInvalidAddress(msg.User.String())
	}
	if msg.Accrediter.Empty() {
		return sdk.ErrInvalidAddress(msg.Accrediter.String())
	}
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSetAccrediter) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSetAccrediter) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// --------------------------
// --- MsgDistributeReward
// --------------------------

// MsgDistributeReward should be used when wanting to distribute a
// specific reward to an accrediter verifying that he has previously
// accreditated all the users.
// Note that the signer should be a trustworthy one in order to avoid
// unauthorized reward distributions.
type MsgDistributeReward struct {
	Accrediter sdk.AccAddress `json:"accrediter"`
	User       sdk.AccAddress `json:"user"`
	Signer     sdk.AccAddress `json:"signer"`
	Reward     sdk.Coins      `json:"reward"`
}

func NewMsgDistributeReward(accrediter sdk.AccAddress, reward sdk.Coins, user, signer sdk.AccAddress) MsgDistributeReward {
	return MsgDistributeReward{
		Accrediter: accrediter,
		User:       user,
		Signer:     signer,
		Reward:     reward,
	}
}

// RouterKey Implements Msg.
func (msg MsgDistributeReward) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgDistributeReward) Type() string { return MsgTypeDistributeReward }

// ValidateBasic Implements Msg.
func (msg MsgDistributeReward) ValidateBasic() sdk.Error {
	if msg.Accrediter.Empty() {
		return sdk.ErrInvalidAddress(msg.Accrediter.String())
	}
	if msg.User.Empty() {
		return sdk.ErrInvalidAddress(msg.User.String())
	}
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}
	if msg.Reward.Empty() {
		return sdk.ErrUnknownRequest("reward cannot be empty")
	}
	if msg.Reward.IsAnyNegative() {
		return sdk.ErrUnknownRequest("rewards cannot be negative")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgDistributeReward) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgDistributeReward) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// --------------------------------
// --- MsgDepositIntoLiquidityPool
// --------------------------------

// MsgDepositIntoLiquidityPool should be used when wanting to deposit a specific
// amount into the liquidity pool which contains the total amount of rewards to
// be distributed upon an accreditation
type MsgDepositIntoLiquidityPool struct {
	Depositor sdk.AccAddress `json:"depositor"`
	Amount    sdk.Coins      `json:"amount"`
}

func NewMsgMsgDepositIntoLiquidityPool(amount sdk.Coins, depositor sdk.AccAddress) MsgDepositIntoLiquidityPool {
	return MsgDepositIntoLiquidityPool{
		Depositor: depositor,
		Amount:    amount,
	}
}

// RouterKey Implements Msg.
func (msg MsgDepositIntoLiquidityPool) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgDepositIntoLiquidityPool) Type() string { return MsgTypesDepositIntoLiquidityPool }

// ValidateBasic Implements Msg.
func (msg MsgDepositIntoLiquidityPool) ValidateBasic() sdk.Error {
	if msg.Depositor.Empty() {
		return sdk.ErrInvalidAddress(msg.Depositor.String())
	}
	if msg.Amount.Empty() {
		return sdk.ErrUnknownRequest("amount cannot be empty")
	}
	if msg.Amount.IsAnyNegative() {
		return sdk.ErrUnknownRequest("amount cannot be negative")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgDepositIntoLiquidityPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgDepositIntoLiquidityPool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Depositor}
}
