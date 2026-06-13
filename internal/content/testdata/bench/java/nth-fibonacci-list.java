class Solution {
    public long[] fibonacci_list(long n) {
        if (n == 0) return new long[0];
        long[] result = new long[(int)n];
        result[0] = 0;
        if (n == 1) return result;
        result[1] = 1;
        for (int i = 2; i < (int)n; i++) {
            result[i] = result[i - 1] + result[i - 2];
        }
        return result;
    }
}
