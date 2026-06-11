public class Solution {
    public long max_and_below_k(long n, long k) {
        if (((k - 1) | k) <= n) return k - 1;
        return k - 2;
    }
}
