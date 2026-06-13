import (
	"strconv"
	"strings"
)

func encode_decode(strs []string) []string {
	// encode
	var sb strings.Builder
	for _, s := range strs {
		sb.WriteString(strconv.Itoa(len(s)))
		sb.WriteByte('#')
		sb.WriteString(s)
	}
	encoded := sb.String()

	// decode
	res := []string{}
	i := 0
	for i < len(encoded) {
		j := i
		for encoded[j] != '#' {
			j++
		}
		length, _ := strconv.Atoi(encoded[i:j])
		res = append(res, encoded[j+1:j+1+length])
		i = j + 1 + length
	}
	return res
}
