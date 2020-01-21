package docs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/commercionetwork/commercionetwork/x/docs/internal/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Test_exportDocuments(t *testing.T) {
	tests := []struct {
		name      string
		documents []Document
	}{
		{
			"no documents",
			[]Document{},
		},
		{
			"some documents",
			[]Document{{UUID: "first"}, {UUID: "second"}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k := keeper.SetupTestInput()

			for _, doc := range tt.documents {
				err := k.SaveDocument(ctx, doc)
				assert.NoError(t, err)
			}

			for _, doc := range exportDocuments(ctx, k) {
				assert.Contains(t, tt.documents, doc)
			}
		})
	}
}

func Test_exportMetadataSchemes(t *testing.T) {
	tests := []struct {
		name    string
		schemes []MetadataSchema
	}{
		{
			"no metadata schemes",
			[]MetadataSchema{},
		},
		{
			"some metadata schemase",
			[]MetadataSchema{{Type: "first"}, {Type: "second"}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k := keeper.SetupTestInput()

			for _, schema := range tt.schemes {
				k.AddSupportedMetadataScheme(ctx, schema)
			}

			for _, schema := range exportMetadataSchemes(ctx, k) {
				assert.Contains(t, tt.schemes, schema)
			}
		})
	}
}

func Test_exportReceipts(t *testing.T) {
	doc1 := keeper.TestingDocument
	doc2 := keeper.TestingDocument

	doc2.UUID = doc1.UUID + "new doc!"

	tests := []struct {
		name           string
		receipts       []DocumentReceipt
		associatedDocs []Document
	}{
		{
			"no receipts",
			[]DocumentReceipt{},
			[]Document{},
		},
		{
			"some receipts",
			[]DocumentReceipt{{DocumentUUID: doc1.UUID, UUID: "doc1"}, {DocumentUUID: doc2.UUID, UUID: "doc2"}},
			[]Document{doc1, doc2},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k := keeper.SetupTestInput()

			for _, doc := range tt.associatedDocs {
				err := k.SaveDocument(ctx, doc)
				assert.NoError(t, err)
			}

			for _, receipt := range tt.receipts {
				err := k.SaveReceipt(ctx, receipt)
				assert.NoError(t, err)
			}

			er := exportReceipts(ctx, k)
			for _, receipt := range er {
				assert.Contains(t, tt.receipts, receipt)
			}

			assert.Len(t, er, len(tt.receipts))
		})
	}
}

func Test_exportTrustedSchemaProviders(t *testing.T) {
	tsp1 := keeper.TestingSender
	tsp2 := keeper.TestingSender2

	tests := []struct {
		name string
		tsps []sdk.AccAddress
	}{
		{
			"no tsps",
			[]sdk.AccAddress{},
		},
		{
			"some tsps",
			[]sdk.AccAddress{tsp1, tsp2},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k := keeper.SetupTestInput()

			for _, tsp := range tt.tsps {
				k.AddTrustedSchemaProposer(ctx, tsp)
			}

			for _, tsp := range exportTrustedSchemaProviders(ctx, k) {
				assert.Contains(t, tt.tsps, tsp)
			}

		})
	}
}
