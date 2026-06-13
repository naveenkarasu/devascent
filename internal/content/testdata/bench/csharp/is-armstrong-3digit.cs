public class Solution {
    public bool is_armstrong_3digit(long n) {
        if (n < 100 || n > 999) return false;
        long s = 0, temp = n;
        while (temp > 0) {
            long r = temp % 10;
            s += r * r * r;
            temp /= 10;
        }
        return s == n;
    }
}
