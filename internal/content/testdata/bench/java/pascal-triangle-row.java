class Solution {
    public long[] pascal_row(long n) {
        long[] row = new long[(int)n];
        for (int k = 0; k < (int)n; k++) {
            row[k] = comb((int)n - 1, k);
        }
        return row;
    }

    private long comb(int n, int k) {
        if (k > n - k) k = n - k;
        long result = 1;
        for (int i = 0; i < k; i++) {
            result = result * (n - i) / (i + 1);
        }
        return result;
    }
}
