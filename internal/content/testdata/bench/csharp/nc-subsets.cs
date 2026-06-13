public class Solution {
    public long[][] subsets(long[] nums) {
        var res = new List<long[]>();
        Backtrack(nums, 0, new List<long>(), res);
        res.Sort((a, b) => {
            int len = Math.Min(a.Length, b.Length);
            for (int i = 0; i < len; i++) {
                if (a[i] != b[i]) return a[i].CompareTo(b[i]);
            }
            return a.Length.CompareTo(b.Length);
        });
        return res.ToArray();
    }

    private void Backtrack(long[] nums, int start, List<long> current, List<long[]> res) {
        long[] sub = current.ToArray();
        Array.Sort(sub);
        res.Add(sub);
        for (int i = start; i < nums.Length; i++) {
            current.Add(nums[i]);
            Backtrack(nums, i + 1, current, res);
            current.RemoveAt(current.Count - 1);
        }
    }
}
