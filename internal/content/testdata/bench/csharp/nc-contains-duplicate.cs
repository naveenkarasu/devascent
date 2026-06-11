public class Solution {
    public bool contains_duplicate(long[] nums) {
        var seen = new HashSet<long>();
        foreach (long n in nums) {
            if (!seen.Add(n)) return true;
        }
        return false;
    }
}
