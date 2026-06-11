type wordDictNode struct {
	children [26]*wordDictNode
	isEnd    bool
}

type wordDictionary struct {
	root *wordDictNode
}

func newWordDictionary() *wordDictionary {
	return &wordDictionary{root: &wordDictNode{}}
}

func (wd *wordDictionary) addWord(word string) {
	node := wd.root
	for _, c := range word {
		idx := c - 'a'
		if node.children[idx] == nil {
			node.children[idx] = &wordDictNode{}
		}
		node = node.children[idx]
	}
	node.isEnd = true
}

func (wd *wordDictionary) search(word string) bool {
	var dfs func(node *wordDictNode, i int) bool
	dfs = func(node *wordDictNode, i int) bool {
		if i == len(word) {
			return node.isEnd
		}
		c := word[i]
		if c == '.' {
			for _, child := range node.children {
				if child != nil && dfs(child, i+1) {
					return true
				}
			}
			return false
		}
		idx := c - 'a'
		if node.children[idx] == nil {
			return false
		}
		return dfs(node.children[idx], i+1)
	}
	return dfs(wd.root, 0)
}

func word_dictionary_ops(operations [][]string) []interface{} {
	wd := newWordDictionary()
	out := make([]interface{}, 0, len(operations))
	for _, op := range operations {
		opName := op[0]
		arg := op[1]
		if opName == "addWord" {
			wd.addWord(arg)
			out = append(out, nil)
		} else {
			out = append(out, wd.search(arg))
		}
	}
	return out
}
