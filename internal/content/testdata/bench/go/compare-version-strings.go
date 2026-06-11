import "strings"
import "strconv"

func compare_versions(v1 string, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")
	length := len(parts1)
	if len(parts2) > length {
		length = len(parts2)
	}
	for i := 0; i < length; i++ {
		a, b := 0, 0
		if i < len(parts1) {
			a, _ = strconv.Atoi(parts1[i])
		}
		if i < len(parts2) {
			b, _ = strconv.Atoi(parts2[i])
		}
		if a > b {
			return 1
		} else if a < b {
			return -1
		}
	}
	return 0
}
