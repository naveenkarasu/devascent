public class Solution {
    public long digit_sum(long n) {
        long s = 0;
        while (n > 0) {
            s += n % 10;
            n /= 10;
        }
        return s;
    }
}
