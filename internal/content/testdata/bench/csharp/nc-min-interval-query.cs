using System.Collections.Generic;

public class Solution {
    public long[] min_interval(long[][] intervals, long[] queries) {
        // Sort intervals by start
        System.Array.Sort(intervals, (a, b) => a[0].CompareTo(b[0]));

        int qLen = queries.Length;
        // indexed queries sorted by query value
        var indexedQueries = new (int origIdx, long q)[qLen];
        for (int i = 0; i < qLen; i++) indexedQueries[i] = (i, queries[i]);
        System.Array.Sort(indexedQueries, (a, b) => a.q.CompareTo(b.q));

        long[] res = new long[qLen];
        for (int i = 0; i < qLen; i++) res[i] = -1;

        // min-heap: (size, end)
        var heap = new PriorityQueue<(long size, long end), long>();
        int idx = 0;

        foreach (var (origIdx, q) in indexedQueries) {
            // Add all intervals whose start <= q
            while (idx < intervals.Length && intervals[idx][0] <= q) {
                long l = intervals[idx][0], r = intervals[idx][1];
                long size = r - l + 1;
                heap.Enqueue((size, r), size);
                idx++;
            }
            // Remove intervals that end before q
            while (heap.Count > 0 && heap.Peek().end < q) {
                heap.Dequeue();
            }
            if (heap.Count > 0) {
                res[origIdx] = heap.Peek().size;
            }
        }

        return res;
    }
}
