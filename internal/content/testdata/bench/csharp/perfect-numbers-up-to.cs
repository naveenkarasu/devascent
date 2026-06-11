using System.Collections.Generic;

public class Solution {
    public long[] perfect_numbers_up_to(long n) {
        var result = new List<long>();
        for (long i = 2; i <= n; i++) {
            long s = 1;
            long r = (long)System.Math.Sqrt((double)i);
            while (r * r > i) r--;
            for (long d = 2; d <= r; d++) {
                if (i % d == 0) {
                    s += d + i / d;
                    if (d * d == i) s -= d;
                }
            }
            if (s == i) result.Add(i);
        }
        return result.ToArray();
    }
}
