public class Solution {
    public long[] gcd_lcm(long a, long b) {
        long x = a, y = b;
        while (y != 0) {
            long t = x % y;
            x = y;
            y = t;
        }
        long g = x;
        long l = a / g * b;
        return new long[] { g, l };
    }
}
