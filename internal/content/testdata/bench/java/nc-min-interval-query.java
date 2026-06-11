import java.util.*;

class Solution {
    public long[] min_interval(long[][] intervals, long[] queries) {
        // Sort intervals by start
        Arrays.sort(intervals, (a, b) -> Long.compare(a[0], b[0]));
        int qLen = queries.length;
        // indexed queries sorted by value
        Long[] indices = new Long[qLen];
        for (int i = 0; i < qLen; i++) indices[i] = (long) i;
        Arrays.sort(indices, (a, b) -> Long.compare(queries[(int)(long)a], queries[(int)(long)b]));
        long[] res = new long[qLen];
        Arrays.fill(res, -1);
        // heap: (size, end)
        PriorityQueue<long[]> heap = new PriorityQueue<>((a, b) -> {
            if (a[0] != b[0]) return Long.compare(a[0], b[0]);
            return Long.compare(a[1], b[1]);
        });
        int i = 0;
        for (long origIdx : indices) {
            long q = queries[(int)origIdx];
            while (i < intervals.length && intervals[i][0] <= q) {
                long l = intervals[i][0], r = intervals[i][1];
                heap.offer(new long[]{r - l + 1, r});
                i++;
            }
            while (!heap.isEmpty() && heap.peek()[1] < q) {
                heap.poll();
            }
            if (!heap.isEmpty()) {
                res[(int)origIdx] = heap.peek()[0];
            }
        }
        return res;
    }
}
