package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/id/types"
)

func getIdentityStoreKey(owner sdk.AccAddress) []byte {
	return append([]byte(types.IdentitiesStorePrefix), owner...)
}

func getDidPowerUpRequestStoreKey(id string) []byte {
	return []byte(types.DidPowerUpRequestStorePrefix + id)
}

func getHandledPowerUpRequestsReferenceStoreKey(reference string) []byte {
	return []byte(types.HandledPowerUpRequestsReferenceStorePrefix + reference)
}
