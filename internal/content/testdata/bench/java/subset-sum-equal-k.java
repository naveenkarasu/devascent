import java.util.*;

class Solution {
    public boolean subset_sum_to_k(long[] nums, long k) {
        int K = (int) k;
        boolean[] prev = new boolean[K + 1];
        prev[0] = true;
        if (nums[0] <= K) prev[(int) nums[0]] = true;
        for (int i = 1; i < nums.length; i++) {
            boolean[] curr = new boolean[K + 1];
            curr[0] = true;
            for (int j = 1; j <= K; j++) {
                boolean notTake = prev[j];
                boolean take = (nums[i] <= j) && prev[j - (int) nums[i]];
                curr[j] = take || notTake;
            }
            prev = curr;
        }
        return prev[K];
    }
}
