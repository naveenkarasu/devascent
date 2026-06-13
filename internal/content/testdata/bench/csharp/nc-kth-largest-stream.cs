using System.Collections.Generic;

public class Solution {
    public long[] kth_largest_stream_ops(long k, long[] nums, object[][] operations) {
        var heap = new SortedList<long, int>(); // min-heap via SortedList with count
        // Use priority queue instead
        var pq = new PriorityQueue<long, long>(); // min-heap
        foreach (long n in nums) pq.Enqueue(n, n);
        while (pq.Count > k) pq.Dequeue();
        var out_ = new long[operations.Length];
        for (int i = 0; i < operations.Length; i++) {
            long val = (long)operations[i][1];
            pq.Enqueue(val, val);
            while (pq.Count > k) pq.Dequeue();
            out_[i] = pq.Peek();
        }
        return out_;
    }
}
