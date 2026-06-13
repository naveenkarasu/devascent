using System.Collections.Generic;

public class Solution {
    public string alien_order(string[] words) {
        var adj = new Dictionary<char, HashSet<char>>();
        var indeg = new Dictionary<char, int>();
        foreach (var w in words)
            foreach (char c in w)
                if (!indeg.ContainsKey(c)) { indeg[c] = 0; adj[c] = new HashSet<char>(); }

        for (int i = 0; i < words.Length - 1; i++) {
            string a = words[i], b = words[i + 1];
            int m = Math.Min(a.Length, b.Length);
            if (a.Length > b.Length && a.Substring(0, m) == b) return "";
            for (int j = 0; j < m; j++) {
                if (a[j] != b[j]) {
                    if (!adj[a[j]].Contains(b[j])) {
                        adj[a[j]].Add(b[j]);
                        indeg[b[j]]++;
                    }
                    break;
                }
            }
        }

        var heap = new SortedSet<char>(indeg.Where(kv => kv.Value == 0).Select(kv => kv.Key));
        var res = new System.Text.StringBuilder();
        while (heap.Count > 0) {
            char c = heap.Min;
            heap.Remove(c);
            res.Append(c);
            foreach (char nxt in adj[c].OrderBy(x => x)) {
                indeg[nxt]--;
                if (indeg[nxt] == 0) heap.Add(nxt);
            }
        }
        return res.Length != indeg.Count ? "" : res.ToString();
    }
}
