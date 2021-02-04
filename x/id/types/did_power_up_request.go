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
	Proof    string         `json:"proof" example:"S5hyg4slMxm9fK8PTNDs8tHmQcBfWXG0vqrNHLXY5K1qUz3QwZYjR9nzJoNDJh18aPsXper7rNBbyZPOm5K//x8Bqm2EJkdnHd7woa5eFqpziGaHxqvgPaLGspH47tnVilARTeF23L2NVHWcEWuo9U5cWg52l1lOixOG+DehT3vC9KjLqg0YqBoL2u0LTLqQMON4UUjC8JwzT/RMs30OYGsWuLc9s48RtJCQJZ+yAg3U6jZn3OokGwWWjYxF9tAsMR48KilHsPigsa9WPnaAyCMSJ05hOqjBxWiSHYiH1nAefFqHtNFXhJF3LRUCJ2xnSHxJC5Ndj4HFzUjyK4aiV1mtRlRcsqmXU80HEk7IzI74HYpW74F8LzXNsh8Pbl7HXoIzEiOHB5XStFnrxkIL3sYAJGH/pGbX3SxeyfoZhY4ikEyqX3OB7Pat2yHh/63XSPThRVpD7g0gy5N2aKBz3vrHCPhe3QQTzWmKlJOcg1FE5ZtSUEHdVQbm1GD9zP6KZDfbekh9+xU0EFczW9JF/we61LTvMF1KoxaBpL46O/J6ROEOQsb03hLEMadBKxZ+XaqAHiQWKu6G5YH2opNTGKcvSyNfDInOvAygUOfzLgTCWp7JOU09hWBKW1ya2yJNJMZ6q9giEAlqS/qqYy4gAqZKjt7nF0siOb3Vz6zEaXdhCcqrfnNN6n/kFXWz24yAucW+/EHt+hsygEVUZQ=="`
	ID       string         `json:"id" example:"d423c645-fd50-4841-8138-192ee8e23dde"`
	ProofKey string         `json:"proof_key" example:"L0QIWxtHeWeUQhmfWqB2n+MZXFqEYctltilM0j69tBd1drUoUSz/vUkaPadQAdKqtQOD43Py7/JZt5IFyx7iDdphzJEX7bqq+B6nC2DQUeISEiXwtDmJYMp20/N23DY2T7L/Z/dzbxRZDWoUhtr9fRPeJL8NHtPqU9YZw2f1tgMk2t/ZMKtBhYzO5BnF8Crmshjw6b6KA3fK+j7YrmF8fVpVFCdz5jd7cprf5RIqwVjt4w1cYZWeKvGLWeGVX3oiCB67EzXZVUCsD03evr90GDY9qGLfUaWJdBkNjByDotLY0OhrKpcZ+O0IZyZv1+YKx7ZDoPAsEJqpqw4M9bGQRg=="`
}
