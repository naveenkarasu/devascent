func ladder_length(begin_word string, end_word string, word_list []string) int {
	wordSet := make(map[string]bool)
	for _, w := range word_list {
		wordSet[w] = true
	}
	if !wordSet[end_word] {
		return 0
	}
	type item struct {
		word   string
		length int
	}
	queue := []item{{begin_word, 1}}
	visited := map[string]bool{begin_word: true}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		word := cur.word
		for i := 0; i < len(word); i++ {
			for c := byte('a'); c <= 'z'; c++ {
				candidate := word[:i] + string(c) + word[i+1:]
				if candidate == end_word {
					return cur.length + 1
				}
				if wordSet[candidate] && !visited[candidate] {
					visited[candidate] = true
					queue = append(queue, item{candidate, cur.length + 1})
				}
			}
		}
	}
	return 0
}
