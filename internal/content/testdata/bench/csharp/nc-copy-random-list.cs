using System.Collections.Generic;

public class Solution {
    private class Node {
        public long val;
        public Node next, random;
        public Node(long v) { val = v; }
    }

    private Node CopyRandomList(Node head) {
        if (head == null) return null;
        var mp = new Dictionary<Node, Node>(ReferenceEqualityComparer.Instance);
        Node cur = head;
        while (cur != null) { mp[cur] = new Node(cur.val); cur = cur.next; }
        cur = head;
        while (cur != null) {
            mp[cur].next = cur.next != null ? mp[cur.next] : null;
            mp[cur].random = cur.random != null ? mp[cur.random] : null;
            cur = cur.next;
        }
        return mp[head];
    }

    public object[][] copy_list(object[][] arr) {
        if (arr == null || arr.Length == 0) return new object[0][];
        int n = arr.Length;
        Node[] nodes = new Node[n];
        for (int i = 0; i < n; i++) {
            nodes[i] = new Node((long)arr[i][0]);
        }
        for (int i = 0; i < n; i++) {
            nodes[i].next = (i + 1 < n) ? nodes[i + 1] : null;
            object r = arr[i][1];
            nodes[i].random = (r != null) ? nodes[(int)(long)r] : null;
        }
        Node copy = CopyRandomList(nodes[0]);
        var order = new Dictionary<Node, int>(ReferenceEqualityComparer.Instance);
        Node c = copy;
        int idx = 0;
        while (c != null) { order[c] = idx++; c = c.next; }
        var out_list = new List<object[]>();
        c = copy;
        while (c != null) {
            object rnd = (c.random != null) ? (object)(long)order[c.random] : null;
            out_list.Add(new object[]{c.val, rnd});
            c = c.next;
        }
        return out_list.ToArray();
    }
}
