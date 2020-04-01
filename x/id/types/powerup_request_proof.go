package types

import "github.com/cosmos/cosmos-sdk/types"

// PowerUpRequestProof represents the structure of a PowerUp request proof a client should assemble when sending a
// MsgRequestDidPowerUp.
// It isn't strictly used in the blockchain-side of the application, and is kept here
// for order/ease of development's sake.
type PowerUpRequestProof struct {
	SenderDid   types.AccAddress `json:"sender_did"`
	PairwiseDid types.AccAddress `json:"pairwise_did"`
	Timestamp   int64            `json:"timestamp"`
	Signature   string           `json:"signature"`
}
