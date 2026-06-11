public class Solution {
    public long find_unique_among_k(long k, long[] arr) {
        var counts = new Dictionary<long, long>();
        foreach (long x in arr) {
            if (!counts.ContainsKey(x)) counts[x] = 0;
            counts[x]++;
        }
        foreach (var kv in counts) {
            if (kv.Value % k != 0) return kv.Key;
        }
        return -1;
    }
}
