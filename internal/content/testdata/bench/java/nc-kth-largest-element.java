import java.util.*;

class Solution {
    public long find_kth_largest(long[] nums, long k) {
        PriorityQueue<Long> heap = new PriorityQueue<>();
        for (long n : nums) {
            heap.add(n);
            if (heap.size() > k) heap.poll();
        }
        return heap.poll();
    }
}
