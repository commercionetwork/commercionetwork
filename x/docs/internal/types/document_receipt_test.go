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

func TestDocumentReceipts_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		receipts DocumentReceipts
		empty    bool
	}{
		{
			"really empty",
			DocumentReceipts{},
			true,
		},
		{
			"not empty",
			DocumentReceipts{
				DocumentReceipt{
					"uuid",
					sender,
					recipient,
					"txhash",
					"documentuuid",
					"proof",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.empty, tt.receipts.IsEmpty())
		})
	}
}

func TestDocumentReceipts_AppendIfMissing(t *testing.T) {
	tests := []struct {
		name           string
		us             DocumentReceipts
		newData        DocumentReceipt
		want           DocumentReceipts
		alreadyPresent bool
	}{
		{
			"adding a new element",
			DocumentReceipts{},
			DocumentReceipt{
				"uuid",
				sender,
				recipient,
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipts{
				DocumentReceipt{
					"uuid",
					sender,
					recipient,
					"txhash",
					"documentuuid",
					"proof",
				},
			},
			true,
		},
		{
			"adding an already present element",
			DocumentReceipts{
				DocumentReceipt{
					"uuid",
					sender,
					recipient,
					"txhash",
					"documentuuid",
					"proof",
				},
			},
			DocumentReceipt{
				"uuid",
				sender,
				recipient,
				"txhash",
				"documentuuid",
				"proof",
			},
			DocumentReceipts{
				DocumentReceipt{
					"uuid",
					sender,
					recipient,
					"txhash",
					"documentuuid",
					"proof",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			val, present := tt.us.AppendIfMissing(tt.newData)

			assert.Equal(t, tt.alreadyPresent, present)
			assert.Equal(t, tt.want, val)
		})
	}
}

func TestDocumentReceipts_AppendAllIfMissing(t *testing.T) {
	tests := []struct {
		name        string
		receipts    DocumentReceipts
		newReceipts DocumentReceipts
		want        DocumentReceipts
	}{
		{
			"append DocumentReceipts to an existing one",
			DocumentReceipts{
				DocumentReceipt{
					"uuid",
					sender,
					recipient,
					"txhash",
					"documentuuid",
					"proof",
				},
			},
			DocumentReceipts{
				DocumentReceipt{
					"uuid2",
					sender,
					recipient,
					"txhash2",
					"documentuuid2",
					"proof2",
				},
			},
			DocumentReceipts{
				DocumentReceipt{
					"uuid2",
					sender,
					recipient,
					"txhash2",
					"documentuuid2",
					"proof2",
				},
				DocumentReceipt{
					"uuid",
					sender,
					recipient,
					"txhash",
					"documentuuid",
					"proof",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.want, tt.receipts.AppendAllIfMissing(tt.newReceipts))
		})
	}
}

func TestDocumentReceipts_FindByDocumentID(t *testing.T) {
	tests := []struct {
		name     string
		receipts DocumentReceipts
		id       string
		want     DocumentReceipts
	}{
		{
			"find a document with given id",
			DocumentReceipts{
				DocumentReceipt{
					"uuid",
					sender,
					recipient,
					"txhash",
					"documentuuid",
					"proof",
				},
			},
			"documentuuid",
			DocumentReceipts{
				DocumentReceipt{
					"uuid",
					sender,
					recipient,
					"txhash",
					"documentuuid",
					"proof",
				},
			},
		},
		{
			"did not find document with given id",
			DocumentReceipts{
				DocumentReceipt{
					"uuid",
					sender,
					recipient,
					"txhash",
					"documentuuid",
					"proof",
				},
			},
			"documentuuidd",
			DocumentReceipts{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.want, tt.receipts.FindByDocumentID(tt.id))
		})
	}
}
