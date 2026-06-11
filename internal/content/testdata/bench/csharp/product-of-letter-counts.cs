using System.Collections.Generic;

public class Solution {
    public long product_of_counts(string s) {
        const long MOD = 1_000_000_007;
        var freq = new Dictionary<char, long>();
        foreach (var ch in s) {
            if (!freq.ContainsKey(ch)) freq[ch] = 0;
            freq[ch]++;
        }
        long result = 1;
        foreach (var cnt in freq.Values) {
            result = (result * cnt) % MOD;
        }
        return result;
    }
}
