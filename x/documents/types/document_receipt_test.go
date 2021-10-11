package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDocumentReceipt_Equals(t *testing.T) {

	tests := []struct {
		name  string
		us    DocumentReceipt
		them  DocumentReceipt
		equal bool
	}{
		{
			"two equal DocumentReceipt",
			DocumentReceipt{
				"uuid",
				sender.String(),
				recipient.String(),
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipt{
				"uuid",
				sender.String(),
				recipient.String(),
				"txhash",
				"documentuuid",
				"proof",
			},
			true,
		},
		{
			"different in proof",
			DocumentReceipt{
				"uuid",
				sender.String(),
				recipient.String(),
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipt{
				"uuid",
				sender.String(),
				recipient.String(),
				"txhash",
				"documentuuid",
				"",
			},
			false,
		},
		{
			"different in documentuuid",
			DocumentReceipt{
				"uuid",
				sender.String(),
				recipient.String(),
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipt{
				"uuid",
				sender.String(),
				recipient.String(),
				"txhash",
				"",
				"proof",
			},
			false,
		},
		{
			"different in txhash",
			DocumentReceipt{
				"uuid",
				sender.String(),
				recipient.String(),
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipt{
				"uuid",
				sender.String(),
				recipient.String(),
				"",
				"documentuuid",
				"proof",
			},
			false,
		},
		{
			"different in recipient",
			DocumentReceipt{
				"uuid",
				sender.String(),
				recipient.String(),
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipt{
				"uuid",
				sender.String(),
				"",
				"txhash",
				"documentuuid",
				"proof",
			},
			false,
		},
		{
			"different in sender",
			DocumentReceipt{
				"uuid",
				sender.String(),
				recipient.String(),
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipt{
				"uuid",
				"",
				recipient.String(),
				"txhash",
				"documentuuid",
				"proof",
			},
			false,
		},
		{
			"different in uuid",
			DocumentReceipt{
				"uuid",
				sender.String(),
				recipient.String(),
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipt{
				"",
				sender.String(),
				recipient.String(),
				"txhash",
				"documentuuid",
				"proof",
			},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.equal, tt.us.Equals(tt.them))
		})
	}
}
