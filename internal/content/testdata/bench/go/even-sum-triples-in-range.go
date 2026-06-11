func count_even_sum_triples(arr []int, l int, r int) int {
	sub := arr[l-1 : r]
	e := 0
	for _, x := range sub {
		if x%2 == 0 {
			e++
		}
	}
	o := len(sub) - e
	return (e*(e-1)*(e-2))/6 + (o*(o-1)/2)*e
}
