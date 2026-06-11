using System.Collections.Generic;

public class Solution {
    public long[] max_sliding_window(long[] nums, long k) {
        int n = nums.Length;
        int ki = (int)k;
        long[] result = new long[n - ki + 1];
        var dq = new LinkedList<int>();
        int ri = 0;
        for (int i = 0; i < n; i++) {
            while (dq.Count > 0 && dq.First.Value < i - ki + 1) dq.RemoveFirst();
            while (dq.Count > 0 && nums[dq.Last.Value] < nums[i]) dq.RemoveLast();
            dq.AddLast(i);
            if (i >= ki - 1) result[ri++] = nums[dq.First.Value];
        }
        return result;
    }
}
