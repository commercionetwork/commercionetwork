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

func AppendDidIfMissing(slice []types.Did, i types.Did) []types.Did {
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

func DidInSlice(a types.Did, list []types.Did) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
