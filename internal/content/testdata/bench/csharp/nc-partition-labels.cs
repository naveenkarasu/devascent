using System.Collections.Generic;

public class Solution {
    public long[] partition_labels(string s) {
        int[] last = new int[26];
        for (int i = 0; i < s.Length; i++) last[s[i] - 'a'] = i;
        var result = new List<long>();
        int start = 0, end = 0;
        for (int i = 0; i < s.Length; i++) {
            end = Math.Max(end, last[s[i] - 'a']);
            if (i == end) {
                result.Add((long)(end - start + 1));
                start = i + 1;
            }
        }
        return result.ToArray();
    }
}
