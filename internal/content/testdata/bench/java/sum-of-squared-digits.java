class Solution {
    public long sum_squared_digits(long n) {
        long total = 0;
        while (n > 0) {
            long d = n % 10;
            total += d * d;
            n /= 10;
        }
        return total;
    }
}
