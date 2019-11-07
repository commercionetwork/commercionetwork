package types_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestRequestStatus_Validate(t *testing.T) {

	err := types.NewRequestStatus("invalid", "message").Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Result().Log, "Invalid status type: invalid")

	assert.NoError(t, types.NewRequestStatus("rejected", "").Validate())
	assert.NoError(t, types.NewRequestStatus("canceled", "").Validate())
}
