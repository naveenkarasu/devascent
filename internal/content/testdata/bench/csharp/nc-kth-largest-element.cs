public class Solution {
    public long find_kth_largest(long[] nums, long k) {
        var heap = new PriorityQueue<long, long>();
        foreach (long n in nums) {
            heap.Enqueue(n, n); // min-heap
            if (heap.Count > k) heap.Dequeue();
        }
        return heap.Dequeue();
    }
}
