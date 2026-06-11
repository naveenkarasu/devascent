public class Solution {
    public bool is_factorion(long n) {
        long original = n;
        long s = 0;
        long temp = n;
        while (temp > 0) {
            long d = temp % 10;
            long fact = 1;
            for (long i = 2; i <= d; i++) fact *= i;
            s += fact;
            temp /= 10;
        }
        return s == original;
    }
}
