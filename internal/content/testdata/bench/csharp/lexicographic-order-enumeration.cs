using System.Collections.Generic;

public class Solution {
    public long[] lexical_order(long n) {
        var result = new List<long>();
        long cur = 1;
        while (result.Count < (int)n) {
            result.Add(cur);
            if (cur * 10 <= n) {
                cur *= 10;
            } else {
                while (cur % 10 == 9 || cur >= n) {
                    cur /= 10;
                }
                cur += 1;
            }
        }
        return result.ToArray();
    }
}
