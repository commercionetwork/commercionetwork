package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/upgrade/types"
)

var _ types.QueryServer = Keeper{}
