import "sort"

func majority_element_ii(nums []int) []int {
	cnt1, cnt2 := 0, 0
	num1, num2 := 0, 1
	for _, n := range nums {
		if num1 == n {
			cnt1++
		} else if num2 == n {
			cnt2++
		} else if cnt1 == 0 {
			num1, cnt1 = n, 1
		} else if cnt2 == 0 {
			num2, cnt2 = n, 1
		} else {
			cnt1--
			cnt2--
		}
	}
	cnt1, cnt2 = 0, 0
	for _, n := range nums {
		if n == num1 {
			cnt1++
		} else if n == num2 {
			cnt2++
		}
	}
	res := []int{}
	if cnt1 > len(nums)/3 {
		res = append(res, num1)
	}
	if cnt2 > len(nums)/3 {
		res = append(res, num2)
	}
	sort.Ints(res)
	return res
}
