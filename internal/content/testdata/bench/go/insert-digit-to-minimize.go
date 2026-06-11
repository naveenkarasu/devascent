func insert_to_minimize(a string) string {
	if a[0] == '1' {
		return string(a[0]) + "0" + a[1:]
	}
	return "1" + a
}
