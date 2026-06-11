public class Solution {
    public bool subset_sum_to_k(long[] nums, long k) {
        int K = (int)k;
        bool[] prev = new bool[K + 1];
        prev[0] = true;
        if (nums[0] <= K) prev[(int)nums[0]] = true;
        for (int i = 1; i < nums.Length; i++) {
            bool[] curr = new bool[K + 1];
            curr[0] = true;
            for (int j = 1; j <= K; j++) {
                bool notTake = prev[j];
                bool take = (nums[i] <= j) ? prev[j - (int)nums[i]] : false;
                curr[j] = take || notTake;
            }
            prev = curr;
        }
        return prev[K];
    }
}
