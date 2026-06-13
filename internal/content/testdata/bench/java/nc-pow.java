class Solution {
    public long my_pow(long x, long n) {
        if (n == 0) return 1L;
        if (n < 0) {
            // For integer pow with negative n, result would be 0 for |x|>1
            // but handle as double logic then cast
            return 0L;
        }
        long half = my_pow(x, n / 2);
        if (n % 2 == 0) {
            return half * half;
        } else {
            return half * half * x;
        }
    }
}
