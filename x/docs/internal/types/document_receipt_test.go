package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
				sender,
				recipient,
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipt{
				"uuid",
				sender,
				recipient,
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
				sender,
				recipient,
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipt{
				"uuid",
				sender,
				recipient,
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
				sender,
				recipient,
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipt{
				"uuid",
				sender,
				recipient,
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
				sender,
				recipient,
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipt{
				"uuid",
				sender,
				recipient,
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
				sender,
				recipient,
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipt{
				"uuid",
				sender,
				nil,
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
				sender,
				recipient,
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipt{
				"uuid",
				nil,
				recipient,
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
				sender,
				recipient,
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipt{
				"",
				sender,
				recipient,
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
			assert.Equal(t, tt.equal, tt.us.Equals(tt.them))
		})
	}
}
