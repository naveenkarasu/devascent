import java.util.*;

class Solution {
    public long[] perfect_numbers_up_to(long n) {
        List<Long> result = new ArrayList<>();
        for (long i = 2; i <= n; i++) {
            long s = 1;
            long r = (long) Math.sqrt((double) i);
            for (long d = 2; d <= r; d++) {
                if (i % d == 0) {
                    s += d + i / d;
                    if (d * d == i) s -= d;
                }
            }
            if (s == i) result.add(i);
        }
        long[] arr = new long[result.size()];
        for (int i = 0; i < result.size(); i++) arr[i] = result.get(i);
        return arr;
    }
}
