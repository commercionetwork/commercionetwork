package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func getIdentityStoreKey(owner sdk.AccAddress) []byte {
	return append([]byte(types.IdentitiesStorePrefix), owner...)
}

func getDepositRequestStoreKey(proof string) []byte {
	return []byte(types.DidDepositRequestStorePrefix + proof)
}

func getDidPowerUpRequestStoreKey(proof string) []byte {
	return []byte(types.DidPowerUpRequestStorePrefix + proof)
}

func getHandledPowerUpRequestsReferenceStoreKey(reference string) []byte {
	return []byte(types.HandledPowerUpRequestsReferenceStorePrefix + reference)
}
