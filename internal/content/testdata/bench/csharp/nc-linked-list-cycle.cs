public class Solution {
    private class CycleNode {
        public long val;
        public CycleNode next;
        public CycleNode(long v) { val = v; }
    }

    public bool has_cycle(long[] values, long pos) {
        if (values == null || values.Length == 0) return false;
        var nodes = new CycleNode[values.Length];
        for (int i = 0; i < values.Length; i++) nodes[i] = new CycleNode(values[i]);
        for (int i = 0; i < values.Length - 1; i++) nodes[i].next = nodes[i + 1];
        if (pos >= 0) nodes[values.Length - 1].next = nodes[(int)pos];

        var slow = nodes[0]; var fast = nodes[0];
        while (fast != null && fast.next != null) {
            slow = slow.next;
            fast = fast.next.next;
            if (slow == fast) return true;
        }
        return false;
    }
}
