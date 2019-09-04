package txreward

import (
	"github.com/commercionetwork/commercionetwork/x/txreward/internal/keeper"
)

var msgIncrementsBRPool = MsgIncrementsBlockRewardsPool{
	Funder: keeper.TestFunder,
	Amount: keeper.TestAmount,
}

var handler = NewHandler()
