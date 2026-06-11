public class Solution {
    public string[][] group_anagrams(string[] strs) {
        var groups = new Dictionary<string, List<string>>();
        foreach (var s in strs) {
            var key = string.Concat(s.OrderBy(c => c));
            if (!groups.ContainsKey(key)) groups[key] = new List<string>();
            groups[key].Add(s);
        }
        var result = groups.Values
            .Select(g => g.OrderBy(x => x).ToArray())
            .OrderBy(g => g, Comparer<string[]>.Create((a, b) => {
                for (int i = 0; i < Math.Min(a.Length, b.Length); i++) {
                    int c = string.Compare(a[i], b[i], StringComparison.Ordinal);
                    if (c != 0) return c;
                }
                return a.Length.CompareTo(b.Length);
            }))
            .ToArray();
        return result;
    }
}
