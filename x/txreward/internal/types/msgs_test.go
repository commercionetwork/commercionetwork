package types

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/txreward/internal/keeper"
	"github.com/stretchr/testify/assert"
)

var msgIncrementsBRPool = MsgIncrementBlockRewardsPool{
	Funder: keeper.TestFunder,
	Amount: keeper.TestAmount,
}

func TestMsgIncrementBlockRewardsPool_Route(t *testing.T) {
	actual := msgIncrementsBRPool.Route()
	expected := ModuleName

	assert.Equal(t, expected, actual)
}

func TestMsgIncrementBlockRewardsPool_Type(t *testing.T) {
	actual := msgIncrementsBRPool.Type()
	expected := MsgTypeIncrementBlockRewardsPool

	assert.Equal(t, expected, actual)
}

func TestMsgIncrementBlockRewardsPool_ValidateBasic_valid(t *testing.T) {
	actual := msgIncrementsBRPool.ValidateBasic()

	assert.Nil(t, actual)
}

func TestMsgIncrementBlockRewardsPool_ValidateBasic_invalid(t *testing.T) {

}
