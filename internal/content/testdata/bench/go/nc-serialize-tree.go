func codec_roundtrip(root *TreeNode) *TreeNode {
	// serialize using preorder DFS
	var parts []string
	var ser func(n *TreeNode)
	ser = func(n *TreeNode) {
		if n == nil {
			parts = append(parts, "#")
			return
		}
		// convert int to string without strconv
		v := n.Val
		if v == 0 {
			parts = append(parts, "0")
		} else {
			neg := false
			if v < 0 {
				neg = true
				v = -v
			}
			digits := []byte{}
			for v > 0 {
				digits = append([]byte{byte('0' + v%10)}, digits...)
				v /= 10
			}
			if neg {
				digits = append([]byte{'-'}, digits...)
			}
			parts = append(parts, string(digits))
		}
		ser(n.Left)
		ser(n.Right)
	}
	ser(root)
	// join with comma
	data := ""
	for i, p := range parts {
		if i > 0 {
			data += ","
		}
		data += p
	}
	// deserialize
	idx := 0
	var atoi func(s string) int
	atoi = func(s string) int {
		neg := false
		start := 0
		if len(s) > 0 && s[0] == '-' {
			neg = true
			start = 1
		}
		result := 0
		for i := start; i < len(s); i++ {
			result = result*10 + int(s[i]-'0')
		}
		if neg {
			return -result
		}
		return result
	}
	// split by comma
	tokens := []string{}
	cur := ""
	for _, ch := range data {
		if ch == ',' {
			tokens = append(tokens, cur)
			cur = ""
		} else {
			cur += string(ch)
		}
	}
	tokens = append(tokens, cur)
	var build func() *TreeNode
	build = func() *TreeNode {
		if idx >= len(tokens) {
			return nil
		}
		v := tokens[idx]
		idx++
		if v == "#" {
			return nil
		}
		node := &TreeNode{Val: atoi(v)}
		node.Left = build()
		node.Right = build()
		return node
	}
	return build()
}
