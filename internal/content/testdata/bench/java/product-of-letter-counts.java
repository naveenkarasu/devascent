import java.util.*;
class Solution {
    public long product_of_counts(String s) {
        long MOD = 1_000_000_007L;
        Map<Character, Long> freq = new HashMap<>();
        for (char ch : s.toCharArray()) {
            freq.put(ch, freq.getOrDefault(ch, 0L) + 1);
        }
        long result = 1;
        for (long cnt : freq.values()) {
            result = (result * cnt) % MOD;
        }
        return result;
    }
}
