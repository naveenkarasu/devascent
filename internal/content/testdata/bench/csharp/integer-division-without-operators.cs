public class Solution {
    public long divide_integers(long dividend, long divisor) {
        long INT_MAX = 2147483647L;
        long INT_MIN = -2147483648L;
        if (dividend == INT_MIN && divisor == -1) return INT_MAX;
        bool negative = (dividend < 0) != (divisor < 0);
        long a = Math.Abs(dividend), b = Math.Abs(divisor);
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
