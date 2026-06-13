public class Solution {
    public bool is_anagram(string s, string t) {
        if (s.Length != t.Length) return false;
        var counts = new Dictionary<char, int>();
        foreach (char c in s) {
            if (!counts.ContainsKey(c)) counts[c] = 0;
            counts[c]++;
        }
        foreach (char c in t) {
            if (!counts.ContainsKey(c) || counts[c] == 0) return false;
            counts[c]--;
        }
        return true;
    }
}
