class Solution {
    public long divide_integers(long dividend, long divisor) {
        long INT_MAX = (1L << 31) - 1;
        long INT_MIN = -(1L << 31);
        if (dividend == INT_MIN && divisor == -1) return INT_MAX;
        boolean negative = (dividend < 0) != (divisor < 0);
        long a = Math.abs(dividend), b = Math.abs(divisor);
        long result = 0;
        while (a >= b) {
            long temp = b, multiple = 1;
            while (a >= (temp << 1)) {
                temp <<= 1;
                multiple <<= 1;
            }
            a -= temp;
            result += multiple;
        }
        return negative ? -result : result;
    }
}
