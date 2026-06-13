using System.Collections.Generic;

public class Solution {
    public string min_window(string s, string t) {
        if (t.Length == 0 || s.Length == 0) return "";
        var need = new Dictionary<char, int>();
        foreach (char c in t) need[c] = need.GetValueOrDefault(c, 0) + 1;
        var have = new Dictionary<char, int>();
        int formed = 0, required = need.Count;
        int left = 0, bestLen = int.MaxValue, bestStart = 0;
        for (int right = 0; right < s.Length; right++) {
            char ch = s[right];
            have[ch] = have.GetValueOrDefault(ch, 0) + 1;
            if (need.ContainsKey(ch) && have[ch] == need[ch]) formed++;
            while (formed == required) {
                if (right - left + 1 < bestLen) { bestLen = right - left + 1; bestStart = left; }
                char lc = s[left];
                have[lc]--;
                if (need.ContainsKey(lc) && have[lc] < need[lc]) formed--;
                left++;
            }
        }
        return bestLen == int.MaxValue ? "" : s.Substring(bestStart, bestLen);
    }
}
