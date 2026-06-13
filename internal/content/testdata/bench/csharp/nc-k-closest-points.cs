public class Solution {
    public long[][] k_closest(long[][] points, long k) {
        var heap = new PriorityQueue<long[], long>();
        foreach (long[] p in points) {
            long dist2 = p[0]*p[0] + p[1]*p[1];
            heap.Enqueue(p, dist2);
        }
        long[][] result = new long[(int)k][];
        for (int i = 0; i < (int)k; i++) {
            result[i] = heap.Dequeue();
        }
        return result;
    }
}
