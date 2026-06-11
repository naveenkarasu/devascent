import "sort"

func car_fleet(target int, position []int, speed []int) int {
	type car struct{ pos, spd int }
	n := len(position)
	cars := make([]car, n)
	for i := range position {
		cars[i] = car{position[i], speed[i]}
	}
	sort.Slice(cars, func(i, j int) bool {
		return cars[i].pos > cars[j].pos
	})
	stack := []float64{}
	for _, c := range cars {
		time := float64(target-c.pos) / float64(c.spd)
		if len(stack) == 0 || time > stack[len(stack)-1] {
			stack = append(stack, time)
		}
	}
	return len(stack)
}
