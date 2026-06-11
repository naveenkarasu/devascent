public class Solution {
    public long sum_of_factorials(long n) {
        long total = 0;
        long fact = 1;
        for (long i = 1; i <= n; i++) {
            fact *= i;
            total += fact;
        }
        return total;
    }
}
