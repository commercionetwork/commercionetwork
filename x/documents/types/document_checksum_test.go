package types

import (
	"testing"
)

var validDocumentChecksum = DocumentChecksum{
	Value:     "93dfcaf3d923ec47edb8580667473987",
	Algorithm: "md5",
}

var anotherValidDocumentChecksum = DocumentChecksum{
	Value:     "D13519A356BEB6F2D993848AA29ECB8C07F4E80F",
	Algorithm: "sha-1",
}

func TestDocumentChecksum_Validate(t *testing.T) {
	type fields struct {
		Value     string
		Algorithm string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "invalid algorithm length",
			fields: fields(validDocumentChecksum),
			wantErr: false,
		},
		{
			name: "empty value",
			fields: fields{
				Value:     "",
				Algorithm: validDocumentChecksum.Algorithm,
			},
			wantErr: true,
		},
		{
			name: "empty algorithm",
			fields: fields{
				Value:     validDocumentChecksum.Value,
				Algorithm: "",
			},
			wantErr: true,
		},
		{
			name: "invalid value",
			fields: fields{
				Value:     validDocumentChecksum.Value + "x",
				Algorithm: validDocumentChecksum.Algorithm,
			},
			wantErr: true,
		},
		{
			name: "invalid algorithm type",
			fields: fields{
				Value:     validDocumentChecksum.Value,
				Algorithm: validDocumentChecksum.Algorithm + "x",
			},
			wantErr: true,
		},
		{
			name: "invalid value length",
			fields: fields{
				Value:     anotherValidDocumentChecksum.Value,
				Algorithm: validDocumentChecksum.Algorithm,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checksum := DocumentChecksum{
				Value:     tt.fields.Value,
				Algorithm: tt.fields.Algorithm,
			}
			if err := checksum.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("DocumentChecksum.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDocumentChecksum_Equals(t *testing.T) {
	type fields struct {
		Value     string
		Algorithm string
	}
	type args struct {
		other DocumentChecksum
	}
	tests := []struct {
		name string
		fields
		args args
		want bool
	}{
		{
			name: "equal",
			fields: fields(validDocumentChecksum),
			args: args{
				other: validDocumentChecksum,
			},
			want: true,
		},
		{
			name: "different value",
			fields: fields{
				Value:     validDocumentChecksum.Value,
				Algorithm: anotherValidDocumentChecksum.Algorithm,
			},
			args: args{
				other: validDocumentChecksum,
			},
			want: false,
		},
		{
			name: "different algorithm",
			fields: fields{
				Value:     anotherValidDocumentChecksum.Value,
				Algorithm: validDocumentChecksum.Algorithm,
			},
			args: args{
				other: validDocumentChecksum,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checksum := DocumentChecksum{
				Value:     tt.fields.Value,
				Algorithm: tt.fields.Algorithm,
			}
			if got := checksum.Equals(tt.args.other); got != tt.want {
				t.Errorf("DocumentChecksum.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}
