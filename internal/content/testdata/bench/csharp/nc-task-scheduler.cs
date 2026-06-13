using System.Collections.Generic;

public class Solution {
    public long least_interval(string[] tasks, long n) {
        var counts = new Dictionary<char, long>();
        foreach (string t in tasks) {
            char c = t[0];
            if (!counts.ContainsKey(c)) counts[c] = 0;
            counts[c]++;
        }
        long maxCount = 0;
        foreach (long v in counts.Values) if (v > maxCount) maxCount = v;
        long numMax = 0;
        foreach (long v in counts.Values) if (v == maxCount) numMax++;
        return Math.Max(tasks.Length, (maxCount - 1) * (n + 1) + numMax);
    }
}
