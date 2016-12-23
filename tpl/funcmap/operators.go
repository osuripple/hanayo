package funcmap

// BAnd is a bitwise AND between some ints.
func BAnd(i1 int, i ...int) int {
	for _, el := range i {
		i1 &= el
	}
	return i1
}
