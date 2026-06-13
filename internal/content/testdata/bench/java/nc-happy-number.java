import java.util.*;

class Solution {
    public boolean is_happy(long n) {
        Set<Long> seen = new HashSet<>();
        while (n != 1) {
            if (seen.contains(n)) return false;
            seen.add(n);
            long total = 0;
            while (n > 0) {
                long digit = n % 10;
                total += digit * digit;
                n /= 10;
            }
            n = total;
        }
        return true;
    }
}
