import "sort"

func majority_elements_n3(nums []int) []int {
	cnt1, cnt2 := 0, 0
	var num1, num2 *int
	for _, n := range nums {
		nCopy := n
		if num1 != nil && n == *num1 {
			cnt1++
		} else if num2 != nil && n == *num2 {
			cnt2++
		} else if cnt1 == 0 {
			num1 = &nCopy
			cnt1 = 1
		} else if cnt2 == 0 {
			num2 = &nCopy
			cnt2 = 1
		} else {
			cnt1--
			cnt2--
		}
	}
	threshold := len(nums) / 3
	result := []int{}
	if num1 != nil {
		count := 0
		for _, x := range nums {
			if x == *num1 {
				count++
			}
		}
		if count > threshold {
			result = append(result, *num1)
		}
	}
	if num2 != nil && (num1 == nil || *num2 != *num1) {
		count := 0
		for _, x := range nums {
			if x == *num2 {
				count++
			}
		}
		if count > threshold {
			result = append(result, *num2)
		}
	}
	sort.Ints(result)
	return result
}
