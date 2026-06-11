public class Solution {
    public long reverse_integer(long x) {
        long sign = x < 0 ? -1 : 1;
        x = System.Math.Abs(x);
        long rev = 0;
        while (x != 0) {
            rev = rev * 10 + x % 10;
            x /= 10;
        }
        rev *= sign;
        if (rev > 2147483647L || rev < -2147483648L) return 0;
        return rev;
    }
}
