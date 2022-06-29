package types

// Strings represents a slice of strings.
type Strings []string

// AppendIfMissing returns a new slice containing the given element.
// If the element was already present, it won't be appended and false will be returned.
// If it wasn't present, the element will be appended to the list and true will be returned.
func (elements Strings) AppendIfMissing(element string) (Strings, bool) {
	if elements.Contains(element) {
		return elements, false
	}
	return append(elements, element), true
}

// RemoveIfExisting returns a new Addresses instance that does not contain the
// given address.
func (elements Strings) RemoveIfExisting(address string) (Strings, bool) {
	indexOf := elements.IndexOf(address)
	if indexOf == -1 {
		return elements, false
	}
	return append(elements[:indexOf], elements[indexOf+1:]...), true
}

// IndexOf returns the index of the given address inside the addresses array,
// or -1 if such an address was not found
func (elements Strings) IndexOf(address string) int {
	for i, a := range elements {
		if a == address {
			return i
		}
	}
	return -1
}

// Contains returns true iff the given element is present inside the elements slice
func (elements Strings) Contains(element string) bool {
	for _, ele := range elements {
		if ele == element {
			return true
		}
	}
	return false
}

// Equals returns true if elements contain the same data of the other slice, in the same order.
func (elements Strings) Equals(other Strings) bool {
	if len(elements) != len(other) {
		return false
	}

	for index, element := range elements {
		if element != other[index] {
			return false
		}
	}

	return true
}

// Empty returns true if this slice does not contain any address
func (addresses Strings) Empty() bool {
	return len(addresses) == 0
}

// IsSet returns true if this slice does not contain any duplicate value
func (elements Strings) IsSet() bool {
	appears := map[string]struct{}{}
	for _, s := range elements {
		if _, found := appears[s]; found {
			return false
		}
		appears[s] = struct{}{}
	}
	return true
}
