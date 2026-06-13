using System.Collections.Generic;

public class Solution {
    public long length_of_longest_substring(string s) {
        var lastSeen = new Dictionary<char, int>();
        int left = 0, best = 0;
        for (int right = 0; right < s.Length; right++) {
            char ch = s[right];
            if (lastSeen.ContainsKey(ch) && lastSeen[ch] >= left)
                left = lastSeen[ch] + 1;
            lastSeen[ch] = right;
            if (right - left + 1 > best) best = right - left + 1;
        }
        return best;
    }
}
