func pascal_row(n int) []int {
	if n <= 0 {
		return []int{}
	}
	row := make([]int, n)
	row[0] = 1
	for k := 1; k < n; k++ {
		// C(n-1, k) = C(n-1, k-1) * (n-1-k+1) / k = C(n-1, k-1) * (n-k) / k
		row[k] = row[k-1] * (n - k) / k
	}
	return row
}
