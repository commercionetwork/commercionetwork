package types_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestService_Equals(t *testing.T) {
	service := types.NewService("id", "type", "endpoint")

	tests := []struct {
		name  string
		us    types.Service
		them  types.Service
		equal bool
	}{
		{
			"different ID",
			service,
			types.NewService(service.ID+"2", service.Type, service.ServiceEndpoint),
			false,
		},
		{
			"different type",
			service,
			types.NewService(service.ID, service.Type+"other", service.ServiceEndpoint),
			false,
		},
		{
			"different service endpoint",
			service,
			types.NewService(service.ID, service.Type, service.ServiceEndpoint+"/v2"),
			false,
		},
		{
			"equal services",
			service,
			service,
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.equal, tt.us.Equals(tt.them))
		})
	}
}

func TestService_Validate(t *testing.T) {
	tests := []struct {
		name string
		ts   types.Service
		want sdk.Error
	}{
		{
			"missing id",
			types.NewService("  ", "type", "endpoint"),
			sdk.ErrUnknownRequest("Service id cannot be empty"),
		},
		{
			"missing type",
			types.NewService("id", "  ", "endpoint"),
			sdk.ErrUnknownRequest("Service type cannot be empty"),
		},
		{
			"missing endpoint",
			types.NewService("id", "type", "  "),
			sdk.ErrUnknownRequest("Service endpoint cannot be empty"),
		},
		{
			"well-formed service",
			types.NewService("id", "type", "endpoint"),
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != nil {
				require.EqualError(t, tt.ts.Validate(), tt.want.Error())
			} else {
				require.NoError(t, tt.ts.Validate())
			}
		})
	}
}

func TestServices_Eqwuals(t *testing.T) {
	service1 := types.NewService("id-1", "type-1", "endpoint-1")
	service2 := types.NewService("id-2", "type-2", "endpoint-2")

	require.False(t, types.Services{}.Equals(types.Services{service1}))
	require.False(t, types.Services{service1}.Equals(types.Services{service1, service2}))
	require.False(t, types.Services{service1, service2}.Equals(types.Services{service2, service1}))
	require.True(t, types.Services{service1, service2}.Equals(types.Services{service1, service2}))
}

func TestServices_Equals(t *testing.T) {
	service1 := types.NewService("id-1", "type-1", "endpoint-1")
	service2 := types.NewService("id-2", "type-2", "endpoint-2")

	tests := []struct {
		name  string
		us    types.Services
		them  types.Services
		equal bool
	}{
		{
			"empty service",
			types.Services{},
			types.Services{service1},
			false,
		},
		{
			"a slice with one service vs a slice with two services",
			types.Services{service1},
			types.Services{service1, service2},
			false,
		},
		{
			"two slices identical in content but with different ordering",
			types.Services{service1, service2},
			types.Services{service2, service1},
			false,
		},
		{
			"two equal services",
			types.Services{service1, service2},
			types.Services{service1, service2},
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.equal, tt.us.Equals(tt.them))
		})
	}
}
