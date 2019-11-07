package types_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestService_Equals(t *testing.T) {
	service := types.NewService("id", "type", "endpoint")

	assert.False(t, service.Equals(types.NewService(service.ID+"2", service.Type, service.ServiceEndpoint)))
	assert.False(t, service.Equals(types.NewService(service.ID, service.Type+"other", service.ServiceEndpoint)))
	assert.False(t, service.Equals(types.NewService(service.ID, service.Type, service.ServiceEndpoint+"/v2")))
	assert.True(t, service.Equals(service))
}

func TestService_Validate(t *testing.T) {

	err := types.NewService("  ", "type", "endpoint").Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Result().Log, "Service id cannot be empty")

	err = types.NewService("id", "  ", "endpoint").Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Result().Log, "Service type cannot be empty")

	err = types.NewService("id", "type", "    ").Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Result().Log, "Service endpoint cannot be empty")

	assert.NoError(t, types.NewService("id", "type", "endpoint").Validate())
}

func TestServices_Equals(t *testing.T) {
	service1 := types.NewService("id-1", "type-1", "endpoint-1")
	service2 := types.NewService("id-2", "type-2", "endpoint-2")

	assert.False(t, types.Services{}.Equals(types.Services{service1}))
	assert.False(t, types.Services{service1}.Equals(types.Services{service1, service2}))
	assert.False(t, types.Services{service1, service2}.Equals(types.Services{service2, service1}))
	assert.True(t, types.Services{service1, service2}.Equals(types.Services{service1, service2}))
}
