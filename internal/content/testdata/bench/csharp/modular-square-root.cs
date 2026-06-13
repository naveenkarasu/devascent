public class Solution {
    public long modular_sqrt(long n, long m) {
        for (long x = 0; x < m; x++) {
            if ((x * x) % m == n) return x;
        }
        return -1;
    }
}
