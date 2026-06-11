import java.util.*;

class Solution {
    public long[] set_union(long[] a, long[] b) {
        Set<Long> seen = new TreeSet<>();
        for (long x : a) seen.add(x);
        for (long x : b) seen.add(x);
        long[] result = new long[seen.size()];
        int i = 0;
        for (long x : seen) result[i++] = x;
        return result;
    }
}
