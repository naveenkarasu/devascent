public class Solution {
    public long smallest_trailing_zero_multiple(long a, long b) {
        long power = 1;
        for (long i = 0; i < b; i++) {
            power *= 10;
        }
        long g = Gcd(a, power);
        return (a / g) * power;
    }

    private long Gcd(long x, long y) {
        while (y != 0) {
            long t = y;
            y = x % y;
            x = t;
        }
        return x;
    }
}
