using System.Collections.Generic;

public class Solution {
    public bool check_inclusion(string s1, string s2) {
        if (s1.Length > s2.Length) return false;
        var need = new Dictionary<char, int>();
        foreach (char c in s1) need[c] = need.GetValueOrDefault(c, 0) + 1;
        var window = new Dictionary<char, int>();
        int k = s1.Length;
        for (int i = 0; i < s2.Length; i++) {
            char ch = s2[i];
            window[ch] = window.GetValueOrDefault(ch, 0) + 1;
            if (i >= k) {
                char lc = s2[i - k];
                window[lc]--;
                if (window[lc] == 0) window.Remove(lc);
            }
            if (window.Count == need.Count) {
                bool match = true;
                foreach (var kv in need) {
                    if (!window.ContainsKey(kv.Key) || window[kv.Key] != kv.Value) { match = false; break; }
                }
                if (match) return true;
            }
        }
        return false;
    }
}
