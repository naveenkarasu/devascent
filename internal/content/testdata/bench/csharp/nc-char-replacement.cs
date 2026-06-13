using System.Collections.Generic;

public class Solution {
    public long character_replacement(string s, long k) {
        var counts = new Dictionary<char, int>();
        int left = 0, maxCount = 0, best = 0;
        for (int right = 0; right < s.Length; right++) {
            char ch = s[right];
            counts[ch] = counts.GetValueOrDefault(ch, 0) + 1;
            if (counts[ch] > maxCount) maxCount = counts[ch];
            while ((right - left + 1) - maxCount > k) {
                char lc = s[left];
                counts[lc]--;
                left++;
            }
            if (right - left + 1 > best) best = right - left + 1;
        }
        return best;
    }
}
