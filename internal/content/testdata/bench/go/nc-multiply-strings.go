func multiply(num1 string, num2 string) string {
	if num1 == "0" || num2 == "0" {
		return "0"
	}
	m, n := len(num1), len(num2)
	buf := make([]int, m+n)
	for i := m - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			mul := int(num1[i]-'0') * int(num2[j]-'0')
			p1, p2 := i+j, i+j+1
			total := mul + buf[p2]
			buf[p2] = total % 10
			buf[p1] += total / 10
		}
	}
	result := []byte{}
	for _, d := range buf {
		if !(len(result) == 0 && d == 0) {
			result = append(result, byte('0'+d))
		}
	}
	if len(result) == 0 {
		return "0"
	}
	return string(result)
}
