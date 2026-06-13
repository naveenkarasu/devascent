import java.util.*;

class Solution {
    public boolean is_perfect_square(long n) {
        if (n < 0) return false;
        long r = (long) Math.sqrt((double) n);
        // adjust for floating point errors
        while (r * r > n) r--;
        while ((r + 1) * (r + 1) <= n) r++;
        return r * r == n;
    }
}
