package types_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateHex(t *testing.T) {
	assert.False(t, types.ValidateHex(""))
	assert.False(t, types.ValidateHex("    "))
	assert.False(t, types.ValidateHex("dasfasdfdf897987"))
	assert.True(t, types.ValidateHex("6369616f6369616f63"))
}
