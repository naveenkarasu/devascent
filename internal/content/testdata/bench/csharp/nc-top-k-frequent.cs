public class Solution {
    public long[] top_k_frequent(long[] nums, long k) {
        var count = new Dictionary<long, long>();
        foreach (var n in nums) {
            if (!count.ContainsKey(n)) count[n] = 0;
            count[n]++;
        }
        var top = count.Keys
            .OrderBy(x => (-count[x], x))
            .Take((int)k)
            .OrderBy(x => x)
            .ToArray();
        return top;
    }
}
