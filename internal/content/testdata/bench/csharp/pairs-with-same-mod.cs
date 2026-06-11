using System.Collections.Generic;

public class Solution {
    public long count_pairs_same_mod(long[] nums, long divisor) {
        var freq = new Dictionary<long, long>();
        foreach (var x in nums) {
            long r = x % divisor;
            if (!freq.ContainsKey(r)) freq[r] = 0;
            freq[r]++;
        }
        long pairs = 0;
        foreach (var cnt in freq.Values) {
            pairs += cnt * (cnt - 1) / 2;
        }
        return pairs;
    }
}
