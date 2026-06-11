type trieNode struct {
	children [26]*trieNode
	isEnd    bool
}

type trie struct {
	root *trieNode
}

func newTrie() *trie {
	return &trie{root: &trieNode{}}
}

func (t *trie) insert(word string) {
	node := t.root
	for _, c := range word {
		idx := c - 'a'
		if node.children[idx] == nil {
			node.children[idx] = &trieNode{}
		}
		node = node.children[idx]
	}
	node.isEnd = true
}

func (t *trie) search(word string) bool {
	node := t.root
	for _, c := range word {
		idx := c - 'a'
		if node.children[idx] == nil {
			return false
		}
		node = node.children[idx]
	}
	return node.isEnd
}

func (t *trie) startsWith(prefix string) bool {
	node := t.root
	for _, c := range prefix {
		idx := c - 'a'
		if node.children[idx] == nil {
			return false
		}
		node = node.children[idx]
	}
	return true
}

func trie_ops(operations [][]string) []interface{} {
	t := newTrie()
	out := make([]interface{}, 0, len(operations))
	for _, op := range operations {
		opName := op[0]
		arg := op[1]
		switch opName {
		case "insert":
			t.insert(arg)
			out = append(out, nil)
		case "search":
			out = append(out, t.search(arg))
		default: // startsWith
			out = append(out, t.startsWith(arg))
		}
	}
	return out
}
