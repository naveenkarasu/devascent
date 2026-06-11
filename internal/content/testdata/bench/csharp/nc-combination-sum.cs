public class Solution {
    public long[][] combination_sum(long[] candidates, long target) {
        Array.Sort(candidates);
        var res = new List<long[]>();
        Backtrack(candidates, 0, target, new List<long>(), res);
        res.Sort((a, b) => {
            int len = Math.Min(a.Length, b.Length);
            for (int i = 0; i < len; i++) {
                if (a[i] != b[i]) return a[i].CompareTo(b[i]);
            }
            return a.Length.CompareTo(b.Length);
        });
        return res.ToArray();
    }

    private void Backtrack(long[] candidates, int start, long remaining, List<long> current, List<long[]> res) {
        if (remaining == 0) {
            long[] combo = current.ToArray();
            Array.Sort(combo);
            res.Add(combo);
            return;
        }
        for (int i = start; i < candidates.Length; i++) {
            if (candidates[i] > remaining) break;
            current.Add(candidates[i]);
            Backtrack(candidates, i, remaining - candidates[i], current, res);
            current.RemoveAt(current.Count - 1);
        }
    }
}
