public class Solution {
    public long digit_sum(long n) {
        long total = 0;
        while (n > 0) {
            total += n % 10;
            n /= 10;
        }
        return total;
    }
}
