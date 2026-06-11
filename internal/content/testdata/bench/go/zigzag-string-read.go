func zigzag_convert(s string, num_rows int) string {
	if num_rows == 1 {
		return s
	}
	rows := make([][]byte, num_rows)
	for i := range rows {
		rows[i] = []byte{}
	}
	row := 0
	direction := -1
	for i := 0; i < len(s); i++ {
		rows[row] = append(rows[row], s[i])
		if row == 0 || row == num_rows-1 {
			direction = -direction
		}
		row += direction
	}
	result := []byte{}
	for _, r := range rows {
		result = append(result, r...)
	}
	return string(result)
}
