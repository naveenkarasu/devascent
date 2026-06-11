import java.util.*;

class Solution {
    public boolean is_factorion(long n) {
        long original = n;
        long s = 0;
        while (n > 0) {
            long d = n % 10;
            s += factorial(d);
            n /= 10;
        }
        return s == original;
    }

    private long factorial(long x) {
        long f = 1;
        for (long i = 2; i <= x; i++) f *= i;
        return f;
    }
}
