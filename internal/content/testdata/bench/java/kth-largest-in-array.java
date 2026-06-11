import java.util.*;

class Solution {
    public long find_kth_largest(long[] nums, long k) {
        PriorityQueue<Long> heap = new PriorityQueue<>();
        for (long x : nums) {
            heap.offer(x);
            if (heap.size() > (int) k) {
                heap.poll();
            }
        }
        return heap.peek();
    }
}
