import java.util.*;

class Solution {
    public long last_stone_weight(long[] stones) {
        PriorityQueue<Long> heap = new PriorityQueue<>(Collections.reverseOrder());
        for (long s : stones) heap.add(s);
        while (heap.size() > 1) {
            long y = heap.poll();
            long x = heap.poll();
            if (x != y) heap.add(y - x);
        }
        return heap.isEmpty() ? 0 : heap.poll();
    }
}
