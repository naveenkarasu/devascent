import java.util.*;

class Solution {
    public long bounded_sequence_sum(long n, long b, long x, long y) {
        long[] a = new long[(int)(n + 1)];
        a[0] = 0;
        for (int i = 1; i <= n; i++) {
            if (a[i - 1] + x <= b) {
                a[i] = a[i - 1] + x;
            } else {
                a[i] = a[i - 1] - y;
            }
        }
        long sum = 0;
        for (long v : a) sum += v;
        return sum;
    }
}
