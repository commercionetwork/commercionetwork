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
		if ele.Checksum.Value == i.Checksum.Value {
			return slice
		}
	}
	return append(slice, i)
}
