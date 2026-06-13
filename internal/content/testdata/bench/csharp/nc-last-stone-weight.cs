using System.Collections.Generic;

public class Solution {
    public long last_stone_weight(long[] stones) {
        // Max-heap via negation in SortedSet won't handle duplicates; use PriorityQueue in reverse
        var heap = new PriorityQueue<long, long>();
        foreach (long s in stones) heap.Enqueue(s, -s); // max-heap by negative priority
        while (heap.Count > 1) {
            long y = heap.Dequeue();
            long x = heap.Dequeue();
            if (x != y) heap.Enqueue(y - x, -(y - x));
        }
        return heap.Count == 0 ? 0 : heap.Dequeue();
    }
}
