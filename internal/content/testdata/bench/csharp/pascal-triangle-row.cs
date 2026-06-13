public class Solution {
    public long[] pascal_row(long n) {
        long[] row = new long[n];
        for (long k = 0; k < n; k++) {
            row[k] = Comb(n - 1, k);
        }
        return row;
    }
    private long Comb(long n, long k) {
        if (k > n - k) k = n - k;
        long result = 1;
        for (long i = 0; i < k; i++) {
            result = result * (n - i) / (i + 1);
        }
        return result;
    }
}
