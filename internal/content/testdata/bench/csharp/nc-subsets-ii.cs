using System.Collections.Generic;

public class Solution {
    public long[][] subsets_with_dup(long[] nums) {
        var res = new List<List<long>>();
        var sorted = (long[])nums.Clone();
        System.Array.Sort(sorted);

        void Backtrack(int start, List<long> current) {
            res.Add(new List<long>(current));
            for (int i = start; i < sorted.Length; i++) {
                if (i > start && sorted[i] == sorted[i - 1]) continue;
                current.Add(sorted[i]);
                Backtrack(i + 1, current);
                current.RemoveAt(current.Count - 1);
            }
        }

        Backtrack(0, new List<long>());

        // Sort each subset ascending, then sort the list of subsets
        foreach (var s in res) s.Sort();
        res.Sort((a, b) => {
            int len = Math.Min(a.Count, b.Count);
            for (int i = 0; i < len; i++) {
                int cmp = a[i].CompareTo(b[i]);
                if (cmp != 0) return cmp;
            }
            return a.Count.CompareTo(b.Count);
        });

        return res.ConvertAll(s => s.ToArray()).ToArray();
    }
}
