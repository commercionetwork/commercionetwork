package types

import "testing"

func TestGenesisState_Validate(t *testing.T) {
	invalidDocument := ValidDocument
	invalidDocument.Sender = ""

	invalidDocumentReceipt := ValidDocumentReceiptRecipient1
	invalidDocumentReceipt.Sender = ""

	type fields struct {
		Documents []*Document
		Receipts  []*DocumentReceipt
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "empty",
			fields: fields{
				Documents: []*Document{},
				Receipts:  []*DocumentReceipt{},
			},
			wantErr: false,
		},
		{
			name: "document with no receipt",
			fields: fields{
				Documents: []*Document{&ValidDocument},
				Receipts:  []*DocumentReceipt{},
			},
			wantErr: false,
		},
		{
			name: "document and receipt",
			fields: fields{
				Documents: []*Document{&ValidDocument},
				Receipts:  []*DocumentReceipt{&ValidDocumentReceiptRecipient1},
			},
			wantErr: false,
		},
		{
			name: "invalid document",
			fields: fields{
				Documents: []*Document{&invalidDocument},
				Receipts:  []*DocumentReceipt{},
			},
			wantErr: true,
		},
		{
			name: "invalid document",
			fields: fields{
				Documents: []*Document{&ValidDocument},
				Receipts:  []*DocumentReceipt{&invalidDocumentReceipt},
			},
			wantErr: true,
		},
		{
			name: "invalid receipt",
			fields: fields{
				Documents: []*Document{&ValidDocument},
				Receipts:  []*DocumentReceipt{&invalidDocumentReceipt},
			},
			wantErr: true,
		},
		{
			name: "receipt without corresponding document",
			fields: fields{
				Documents: []*Document{},
				Receipts:  []*DocumentReceipt{&ValidDocumentReceiptRecipient1},
			},
			wantErr: true,
		},
		{
			name: "documents with same ID",
			fields: fields{
				Documents: []*Document{&ValidDocument, &ValidDocument},
				Receipts:  []*DocumentReceipt{&ValidDocumentReceiptRecipient1},
			},
			wantErr: true,
		},
		{
			name: "receipts with same ID",
			fields: fields{
				Documents: []*Document{&ValidDocument},
				Receipts:  []*DocumentReceipt{&ValidDocumentReceiptRecipient1, &ValidDocumentReceiptRecipient1},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := GenesisState{
				Documents: tt.fields.Documents,
				Receipts:  tt.fields.Receipts,
			}
			if err := gs.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("GenesisState.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
