import java.util.*;

class Solution {
    public long[] gcd_lcm(long a, long b) {
        long g = gcd(a, b);
        long l = a / g * b;
        return new long[]{g, l};
    }

    private long gcd(long x, long y) {
        while (y != 0) {
            long tmp = y;
            y = x % y;
            x = tmp;
        }
        return x;
    }
}
