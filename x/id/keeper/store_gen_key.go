package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/id/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func getIdentityStoreKey(owner sdk.AccAddress) []byte {
	return append([]byte(types.IdentitiesStorePrefix), owner...)
}

func getDepositRequestStoreKey(proof string) []byte {
	return []byte(types.DidDepositRequestStorePrefix + proof)
}

func getDidPowerUpRequestStoreKey(id string) []byte {
	return []byte(types.DidPowerUpRequestStorePrefix + id)
}

func getHandledPowerUpRequestsReferenceStoreKey(reference string) []byte {
	return []byte(types.HandledPowerUpRequestsReferenceStorePrefix + reference)
}
