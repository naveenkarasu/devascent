import java.util.*;

class Solution {
    public long[] days_breakdown(long d) {
        long y = d / 365;
        d = d % 365;
        long w = d / 7;
        d = d % 7;
        return new long[]{y, w, d};
    }
}
