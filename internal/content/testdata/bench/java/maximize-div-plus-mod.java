import java.util.*;

class Solution {
    public long max_div_plus_mod(long l, long r, long a) {
        long best = r / a + r % a;
        long multiple = (r / a) * a;
        long candidate = multiple - 1;
        if (l <= candidate && candidate <= r) {
            long val = candidate / a + candidate % a;
            if (val > best) best = val;
        }
        return best;
    }
}
