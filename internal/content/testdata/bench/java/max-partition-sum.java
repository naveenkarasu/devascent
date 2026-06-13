class Solution {
    public long max_partition_sum(long[] arr) {
        long total = 0;
        for (long x : arr) total += x;
        long current = 0;
        long ans = 0;
        for (int i = 0; i < arr.length - 1; i++) {
            current += arr[i];
            long best = Math.max(current, total - current);
            if (best > ans) ans = best;
        }
        return ans;
    }
}
