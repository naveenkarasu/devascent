import "sort"

type timeEntry struct {
	timestamp int
	value     string
}

type TimeMap struct {
	store map[string][]timeEntry
}

func newTimeMap() *TimeMap {
	return &TimeMap{store: map[string][]timeEntry{}}
}

func (t *TimeMap) set(key, value string, timestamp int) {
	t.store[key] = append(t.store[key], timeEntry{timestamp, value})
}

func (t *TimeMap) get(key string, timestamp int) string {
	entries := t.store[key]
	if len(entries) == 0 {
		return ""
	}
	i := sort.Search(len(entries), func(j int) bool {
		return entries[j].timestamp > timestamp
	})
	if i == 0 {
		return ""
	}
	return entries[i-1].value
}

func time_map_ops(operations [][]interface{}) []interface{} {
	tm := newTimeMap()
	out := make([]interface{}, len(operations))
	for i, op := range operations {
		switch op[0].(string) {
		case "set":
			tm.set(op[1].(string), op[2].(string), op[3].(int))
			out[i] = nil
		default: // get
			out[i] = tm.get(op[1].(string), op[2].(int))
		}
	}
	return out
}
