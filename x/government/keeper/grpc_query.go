package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/government/types"
)

var _ types.QueryServer = Keeper{}
