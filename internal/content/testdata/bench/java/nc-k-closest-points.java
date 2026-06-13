import java.util.*;

class Solution {
    public long[][] k_closest(long[][] points, long k) {
        PriorityQueue<long[]> heap = new PriorityQueue<>(
            (a, b) -> Long.compare(a[0]*a[0]+a[1]*a[1], b[0]*b[0]+b[1]*b[1])
        );
        for (long[] p : points) heap.add(p);
        long[][] result = new long[(int)k][2];
        for (int i = 0; i < k; i++) {
            long[] p = heap.poll();
            result[i][0] = p[0];
            result[i][1] = p[1];
        }
        return result;
    }
}
