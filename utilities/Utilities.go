package utilities

import "github.com/commercionetwork/commercionetwork/types"

// ----------------------------------
// --- Utility functions
// ----------------------------------

func AppendStringIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func HasDocument(doc types.Document, docList []types.Document) bool {

	for _, currentDoc := range docList {
		if doc.Checksum.Value == currentDoc.Checksum.Value {
			return true
		}
	}
	return false
}

func AppendDocIfMissing(slice []types.Document, i types.Document) []types.Document {
	for _, ele := range slice {
		if ele.Equals(i) {
			return slice
		}
	}
	return append(slice, i)
}

func AppendReceiptIfMissing(slice []types.DocumentReceipt, receipt types.DocumentReceipt) []types.DocumentReceipt {
	for _, ele := range slice {
		if ele.Equals(receipt) {
			return slice
		}
	}
	return append(slice, receipt)
}
