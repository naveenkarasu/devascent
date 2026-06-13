import java.util.*;

class Solution {
    static class Node {
        long val;
        Node next, random;
        Node(long v) { val = v; }
    }

    private Node copyRandomList(Node head) {
        if (head == null) return null;
        Map<Node, Node> mp = new IdentityHashMap<>();
        Node cur = head;
        while (cur != null) { mp.put(cur, new Node(cur.val)); cur = cur.next; }
        cur = head;
        while (cur != null) {
            mp.get(cur).next = (cur.next != null) ? mp.get(cur.next) : null;
            mp.get(cur).random = (cur.random != null) ? mp.get(cur.random) : null;
            cur = cur.next;
        }
        return mp.get(head);
    }

    public Object[][] copy_list(Object[][] arr) {
        if (arr == null || arr.length == 0) return new Object[0][];
        int n = arr.length;
        Node[] nodes = new Node[n];
        for (int i = 0; i < n; i++) {
            nodes[i] = new Node(((Number) arr[i][0]).longValue());
        }
        for (int i = 0; i < n; i++) {
            nodes[i].next = (i + 1 < n) ? nodes[i + 1] : null;
            Object r = arr[i][1];
            nodes[i].random = (r != null) ? nodes[(int) ((Number) r).longValue()] : null;
        }
        Node copy = copyRandomList(nodes[0]);
        // index map for the copied list
        Map<Node, Integer> order = new IdentityHashMap<>();
        Node cur = copy;
        int idx = 0;
        while (cur != null) { order.put(cur, idx++); cur = cur.next; }
        List<Object[]> out = new ArrayList<>();
        cur = copy;
        while (cur != null) {
            Object rnd = (cur.random != null) ? (Object) (long) order.get(cur.random) : null;
            out.add(new Object[]{cur.val, rnd});
            cur = cur.next;
        }
        return out.toArray(new Object[0][]);
    }
}
