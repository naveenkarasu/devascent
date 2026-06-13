func days_breakdown(d int) []int {
	y := d / 365
	d = d % 365
	w := d / 7
	d = d % 7
	return []int{y, w, d}
}
