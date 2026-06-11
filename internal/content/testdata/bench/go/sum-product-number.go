func is_sum_product_number(n int) bool {
	digitSum := 0
	digitProd := 1
	temp := n
	for temp > 0 {
		d := temp % 10
		digitSum += d
		digitProd *= d
		temp /= 10
	}
	return digitSum*digitProd == n
}
