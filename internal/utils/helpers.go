package utils

func GetNextIndex(current *int, max_len int) {
	next_index := *current + 1

	if next_index >= max_len {
		*current = 0
	} else {
		*current = next_index
	}
}
