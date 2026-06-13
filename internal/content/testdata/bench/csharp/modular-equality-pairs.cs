public class Solution {
    public long count_modeq_pairs(long n, long m) {
        long count = 0;
        for (long a = 1; a <= n; a++) {
            for (long b = a + 1; b <= n; b++) {
                if ((m % a) % b == (m % b) % a) count++;
            }
        }
        return count;
    }
}
