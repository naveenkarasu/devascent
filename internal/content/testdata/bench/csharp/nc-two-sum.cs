public class Solution {
    public long[] two_sum(long[] nums, long target) {
        var seen = new Dictionary<long, int>();
        for (int i = 0; i < nums.Length; i++) {
            long need = target - nums[i];
            if (seen.ContainsKey(need)) return new long[] { seen[need], i };
            seen[nums[i]] = i;
        }
        return new long[] {};
    }
}
