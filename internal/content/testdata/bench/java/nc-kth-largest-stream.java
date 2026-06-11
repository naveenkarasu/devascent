import java.util.*;

class Solution {
    public long[] kth_largest_stream_ops(long k, long[] nums, Object[][] operations) {
        PriorityQueue<Long> heap = new PriorityQueue<>(); // min-heap
        for (long n : nums) heap.add(n);
        while (heap.size() > k) heap.poll();
        long[] out = new long[operations.length];
        for (int i = 0; i < operations.length; i++) {
            Object[] op = operations[i];
            long val = ((Number) op[1]).longValue();
            heap.add(val);
            while (heap.size() > k) heap.poll();
            out[i] = heap.peek();
        }
        return out;
    }
}
