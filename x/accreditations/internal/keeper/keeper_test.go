package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_SetAccrediter_NoAccrediter(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.AcKeeper.StoreKey)
	store.Delete(TestUser)

	err := TestUtils.AcKeeper.SetAccrediter(TestUtils.Ctx, TestUser, TestAccrediter)
	assert.Nil(t, err)

	accreditationBz := store.Get(TestUser)
	assert.NotNil(t, accreditationBz)

	var accreditation types.Accreditation
	TestUtils.Cdc.MustUnmarshalBinaryBare(accreditationBz, &accreditation)

}
