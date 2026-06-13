using System.Collections.Generic;

public class Solution {
    public long[] fibonacci_list(long n) {
        if (n == 0) return new long[0];
        if (n == 1) return new long[] { 0 };
        var result = new List<long> { 0, 1 };
        for (long i = 2; i < n; i++) {
            result.Add(result[(int)(i - 1)] + result[(int)(i - 2)]);
        }
        return result.ToArray();
    }
}
