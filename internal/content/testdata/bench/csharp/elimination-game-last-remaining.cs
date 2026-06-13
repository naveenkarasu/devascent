public class Solution {
    public long last_remaining(long n) {
        if (n == 1) return 1;
        return 2 * (1 + n / 2 - last_remaining(n / 2));
    }
}
