public class Solution {
    public long max_product(long[] nums) {
        long maxProd = nums[0];
        long minProd = nums[0];
        long res = nums[0];
        for (int i = 1; i < nums.Length; i++) {
            long n = nums[i];
            long a = n;
            long b = maxProd * n;
            long c = minProd * n;
            long newMax = Math.Max(a, Math.Max(b, c));
            long newMin = Math.Min(a, Math.Min(b, c));
            maxProd = newMax;
            minProd = newMin;
            res = Math.Max(res, maxProd);
        }
        return res;
    }
}
