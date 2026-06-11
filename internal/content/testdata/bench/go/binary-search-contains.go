func contains_element(arr []int, key int) bool {
	lo, hi := 0, len(arr)-1
	for lo <= hi {
		mid := (lo + hi) / 2
		if arr[mid] == key {
			return true
		} else if arr[mid] < key {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return false
}
