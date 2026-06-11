using System.Collections.Generic;

public class Solution {
    public string[][] partition(string s) {
        var res = new List<List<string>>();

        bool IsPalindrome(string t) {
            int l = 0, r = t.Length - 1;
            while (l < r) {
                if (t[l] != t[r]) return false;
                l++; r--;
            }
            return true;
        }

        void Backtrack(int start, List<string> current) {
            if (start == s.Length) {
                res.Add(new List<string>(current));
                return;
            }
            for (int end = start + 1; end <= s.Length; end++) {
                string substr = s.Substring(start, end - start);
                if (IsPalindrome(substr)) {
                    current.Add(substr);
                    Backtrack(end, current);
                    current.RemoveAt(current.Count - 1);
                }
            }
        }

        Backtrack(0, new List<string>());

        res.Sort((a, b) => {
            int len = Math.Min(a.Count, b.Count);
            for (int i = 0; i < len; i++) {
                int cmp = string.CompareOrdinal(a[i], b[i]);
                if (cmp != 0) return cmp;
            }
            return a.Count.CompareTo(b.Count);
        });

        return res.ConvertAll(p => p.ToArray()).ToArray();
    }
}
