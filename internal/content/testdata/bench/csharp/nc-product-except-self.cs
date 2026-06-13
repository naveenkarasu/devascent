public class Solution {
    public long[] product_except_self(long[] nums) {
        int n = nums.Length;
        long[] res = new long[n];
        long prefix = 1;
        for (int i = 0; i < n; i++) {
            res[i] = prefix;
            prefix *= nums[i];
        }
        long suffix = 1;
        for (int i = n - 1; i >= 0; i--) {
            res[i] *= suffix;
            suffix *= nums[i];
        }
        return res;
    }
}
