import (
	"strconv"
	"strings"
)

func longest_ones_run(n int) int {
	binary := strconv.FormatInt(int64(n), 2)
	runs := strings.Split(binary, "0")
	maxLen := 0
	for _, r := range runs {
		if len(r) > maxLen {
			maxLen = len(r)
		}
	}
	return maxLen
}
