import java.util.*;
class Solution {
    public long[] max_sliding_window(long[] nums, long k) {
        int n = nums.length;
        int ki = (int) k;
        long[] result = new long[n - ki + 1];
        Deque<Integer> dq = new ArrayDeque<>();
        int ri = 0;
        for (int i = 0; i < n; i++) {
            while (!dq.isEmpty() && dq.peekFirst() < i - ki + 1) dq.pollFirst();
            while (!dq.isEmpty() && nums[dq.peekLast()] < nums[i]) dq.pollLast();
            dq.addLast(i);
            if (i >= ki - 1) result[ri++] = nums[dq.peekFirst()];
        }
        return result;
    }
}
