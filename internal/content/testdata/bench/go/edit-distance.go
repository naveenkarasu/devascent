func edit_distance(str1 string, str2 string) int {
	m, n := len(str1), len(str2)
	cur := make([]int, n+1)
	for j := 0; j <= n; j++ {
		cur[j] = j
	}
	for i := 1; i <= m; i++ {
		pre := cur[0]
		cur[0] = i
		for j := 1; j <= n; j++ {
			temp := cur[j]
			if str1[i-1] == str2[j-1] {
				cur[j] = pre
			} else {
				minVal := pre
				if cur[j] < minVal {
					minVal = cur[j]
				}
				if cur[j-1] < minVal {
					minVal = cur[j-1]
				}
				cur[j] = 1 + minVal
			}
			pre = temp
		}
	}
	return cur[n]
}
