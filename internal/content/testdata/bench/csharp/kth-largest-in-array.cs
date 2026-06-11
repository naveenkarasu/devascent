using System.Collections.Generic;

public class Solution {
    public long find_kth_largest(long[] nums, long k) {
        var heap = new PriorityQueue<long, long>();
        foreach (long x in nums) {
            heap.Enqueue(x, x);
            if (heap.Count > (int)k) {
                heap.Dequeue();
            }
        }
        return heap.Peek();
    }
}
