package array

// It returns the index of the first element
func IndexOfInt64(items []int64, element int64) (int, bool) {
	var index = -1
	var ok bool

	for i, el := range items {
		if el == element {
			index = i
			ok = true
			break
		}
	}

	return index, ok
}

//Removes value from items.
func PullInt64(items []int64, value int64) []int64 {
	if i,ok := IndexOfInt64(items, value); ok {
		return append(items[:i], items[i+1:]...)
	}
	return items
}