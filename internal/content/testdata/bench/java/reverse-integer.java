import java.util.*;

class Solution {
    public long reverse_integer(long x) {
        long sign = x < 0 ? -1 : 1;
        x = Math.abs(x);
        long rev = 0;
        while (x != 0) {
            rev = rev * 10 + x % 10;
            x /= 10;
        }
        rev *= sign;
        if (rev > Integer.MAX_VALUE || rev < Integer.MIN_VALUE) return 0;
        return rev;
    }
}
