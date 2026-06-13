using System.Collections.Generic;

public class Solution {
    public long[][] combination_sum2(long[] candidates, long target) {
        var res = new List<List<long>>();
        var sorted = (long[])candidates.Clone();
        System.Array.Sort(sorted);

        void Backtrack(int start, List<long> current, long remaining) {
            if (remaining == 0) {
                res.Add(new List<long>(current));
                return;
            }
            for (int i = start; i < sorted.Length; i++) {
                if (sorted[i] > remaining) break;
                if (i > start && sorted[i] == sorted[i - 1]) continue;
                current.Add(sorted[i]);
                Backtrack(i + 1, current, remaining - sorted[i]);
                current.RemoveAt(current.Count - 1);
            }
        }

        Backtrack(0, new List<long>(), target);

        foreach (var c in res) c.Sort();
        res.Sort((a, b) => {
            int len = Math.Min(a.Count, b.Count);
            for (int i = 0; i < len; i++) {
                int cmp = a[i].CompareTo(b[i]);
                if (cmp != 0) return cmp;
            }
            return a.Count.CompareTo(b.Count);
        });

        return res.ConvertAll(c => c.ToArray()).ToArray();
    }
}
