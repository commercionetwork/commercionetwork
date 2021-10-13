package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/did/types"
)

var _ types.QueryServer = Keeper{}
