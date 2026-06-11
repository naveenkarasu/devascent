class Solution {
    public long max_product(long[] nums) {
        long maxProd = nums[0], minProd = nums[0], res = nums[0];
        for (int i = 1; i < nums.length; i++) {
            long n = nums[i];
            long tempMax = Math.max(n, Math.max(maxProd * n, minProd * n));
            long tempMin = Math.min(n, Math.min(maxProd * n, minProd * n));
            maxProd = tempMax;
            minProd = tempMin;
            res = Math.max(res, maxProd);
        }
        return res;
    }
}
