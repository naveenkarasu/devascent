public class Solution {
    public long my_pow(long x, long n) {
        if (n == 0) return 1;
        if (n < 0) {
            x = 1 / x;
            n = -n;
        }
        long half = my_pow(x, n / 2);
        if (n % 2 == 0) return half * half;
        else return half * half * x;
    }
}
