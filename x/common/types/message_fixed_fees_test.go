package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStandardBinding(t *testing.T) {
	tests := []struct {
		name        string
		messageName string
		want        MessageFeeBinding
	}{
		{
			"MsgTypeShareDoc standard fee binding",
			"MsgTypeShareDoc",
			MessageFeeBinding{
				Name: "MsgTypeShareDoc",
				Fee:  StandardFIATFee,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mb := NewStandardBinding(tt.messageName)
			require.Equal(t, tt.want, mb)
		})
	}
}
