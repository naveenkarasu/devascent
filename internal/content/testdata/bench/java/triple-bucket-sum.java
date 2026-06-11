import java.util.*;
class Solution {
    public long triple_bucket_sum(long[] counts) {
        long MOD = 1_000_000_007L;
        long ans = 0;
        int n = counts.length;
        for (int i = 0; i < n - 2; i++) {
            for (int j = i + 1; j < n - 1; j++) {
                for (int k = j + 1; k < n; k++) {
                    ans = (ans + counts[i] * counts[j] % MOD * counts[k]) % MOD;
                }
            }
        }
        return ans;
    }
}
