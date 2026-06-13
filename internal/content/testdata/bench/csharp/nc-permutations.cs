public class Solution {
    public long[][] permute(long[] nums) {
        var res = new List<long[]>();
        bool[] used = new bool[nums.Length];
        Backtrack(nums, used, new List<long>(), res);
        res.Sort((a, b) => {
            for (int i = 0; i < a.Length; i++) {
                if (a[i] != b[i]) return a[i].CompareTo(b[i]);
            }
            return 0;
        });
        return res.ToArray();
    }

    private void Backtrack(long[] nums, bool[] used, List<long> current, List<long[]> res) {
        if (current.Count == nums.Length) {
            res.Add(current.ToArray());
            return;
        }
        for (int i = 0; i < nums.Length; i++) {
            if (!used[i]) {
                used[i] = true;
                current.Add(nums[i]);
                Backtrack(nums, used, current, res);
                current.RemoveAt(current.Count - 1);
                used[i] = false;
            }
        }
    }
}
