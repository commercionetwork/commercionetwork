package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DidDepositRequest represents the request that is sent from a user when he wants to send
// something to his pairwise Did. This request will be read and unencrypted from a central
// identity that will later update the status and send the funds to the pairwise Did
type DidPowerUpRequest struct {
	Status   *RequestStatus `json:"status"`
	Claimant sdk.AccAddress `json:"claimant"`
	Amount   sdk.Coins      `json:"amount"`
	Proof    string         `json:"proof"`
	ID       string         `json:"id"`
}
