type lruNode struct {
	key, val   int
	prev, next *lruNode
}

type LRUCache struct {
	cap        int
	m          map[int]*lruNode
	head, tail *lruNode
}

func newLRU(cap int) *LRUCache {
	head := &lruNode{}
	tail := &lruNode{}
	head.next = tail
	tail.prev = head
	return &LRUCache{cap: cap, m: map[int]*lruNode{}, head: head, tail: tail}
}

func (c *LRUCache) removeNode(n *lruNode) {
	n.prev.next = n.next
	n.next.prev = n.prev
}

func (c *LRUCache) insertFront(n *lruNode) {
	n.next = c.head.next
	n.prev = c.head
	c.head.next.prev = n
	c.head.next = n
}

func (c *LRUCache) lget(key int) int {
	if n, ok := c.m[key]; ok {
		c.removeNode(n)
		c.insertFront(n)
		return n.val
	}
	return -1
}

func (c *LRUCache) lput(key, value int) {
	if n, ok := c.m[key]; ok {
		c.removeNode(n)
		n.val = value
		c.insertFront(n)
		return
	}
	n := &lruNode{key: key, val: value}
	c.m[key] = n
	c.insertFront(n)
	if len(c.m) > c.cap {
		lru := c.tail.prev
		c.removeNode(lru)
		delete(c.m, lru.key)
	}
}

func lru_cache_ops(capacity int, operations [][]interface{}) []interface{} {
	cache := newLRU(capacity)
	out := make([]interface{}, len(operations))
	for i, op := range operations {
		switch op[0].(string) {
		case "get":
			out[i] = cache.lget(op[1].(int))
		default: // put
			cache.lput(op[1].(int), op[2].(int))
			out[i] = nil
		}
	}
	return out
}
