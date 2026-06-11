public class Solution {
    public long max_partition_sum(long[] arr) {
        long total = 0;
        foreach (var x in arr) total += x;
        long current = 0;
        long ans = 0;
        for (int i = 0; i < arr.Length - 1; i++) {
            current += arr[i];
            long candidate = current > total - current ? current : total - current;
            if (candidate > ans) ans = candidate;
        }
        return ans;
    }
}
