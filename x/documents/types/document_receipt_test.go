package types
/*
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
				recipient1.String(),
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipt{
				"uuid",
				sender.String(),
				recipient1.String(),
				"txhash",
				"documentuuid",
				"proof",
			},ValidDocument
				sender.String(),
				recipient1.String(),
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipt{
				"uuid",
				sender.String(),
				recipient1.String(),
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
				recipient1.String(),
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipt{
				"uuid",
				sender.String(),
				recipient1.String(),
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
				recipient1.String(),
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipt{
				"uuid",
				sender.String(),
				recipient1.String(),
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
				recipient1.String(),
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
				recipient1.String(),
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipt{
				"uuid",
				"",
				recipient1.String(),
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
				recipient1.String(),
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipt{
				"",
				sender.String(),
				recipient1.String(),
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

func TestDocumentReceipt_Validate(t *testing.T) {
	tests := []struct {
		name    string
		receipt func() DocumentReceipt
		wantErr bool
	}{
		{
			name: "valid",
			receipt: func() DocumentReceipt {
				return ValidDocumentReceiptRecipient1
			},
			wantErr: false,
		},
		{
			name: "empty UUID",
			receipt: func() DocumentReceipt {
				receipt := ValidDocumentReceiptRecipient1
				receipt.UUID = ""
				return receipt
			},
			wantErr: true,
		},
		{
			name: "invalid sender",
			receipt: func() DocumentReceipt {
				receipt := ValidDocumentReceiptRecipient1
				receipt.Sender = ""
				return receipt
			},
			wantErr: true,
		},
		{
			name: "invalid recipient",
			receipt: func() DocumentReceipt {
				receipt := ValidDocumentReceiptRecipient1
				receipt.Recipient = ""
				return receipt
			},
			wantErr: true,
		},
		{
			name: "empty tx hash",
			receipt: func() DocumentReceipt {
				receipt := ValidDocumentReceiptRecipient1
				receipt.TxHash = ""
				return receipt
			},
			wantErr: true,
		},
		{
			name: "empty document UUID",
			receipt: func() DocumentReceipt {
				receipt := ValidDocumentReceiptRecipient1
				receipt.DocumentUUID = ""
				return receipt
			},
			wantErr: true,
		},
		{
			name: "empty proof",
			receipt: func() DocumentReceipt {
				receipt := ValidDocumentReceiptRecipient2
				receipt.Proof = ""
				return receipt
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.receipt().Validate(); (err != nil) != tt.wantErr {
				t.Errorf("DocumentReceipt.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
*/