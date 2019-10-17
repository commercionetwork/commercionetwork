package types

type Strings []string

func (elements Strings) AppendIfMissing(element string) (Strings, bool) {
	for _, ele := range elements {
		if ele == element {
			return elements, false
		}
	}

	return append(elements, element), true
}

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
