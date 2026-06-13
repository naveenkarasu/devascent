public class Solution {
    public bool merge_triplets(long[][] triplets, long[] target) {
        long[] result = {0, 0, 0};
        foreach (long[] t in triplets) {
            if (t[0] <= target[0] && t[1] <= target[1] && t[2] <= target[2]) {
                result[0] = Math.Max(result[0], t[0]);
                result[1] = Math.Max(result[1], t[1]);
                result[2] = Math.Max(result[2], t[2]);
            }
        }
        return result[0] == target[0] && result[1] == target[1] && result[2] == target[2];
    }
}
