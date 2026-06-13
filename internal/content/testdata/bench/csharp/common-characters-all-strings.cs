public class Solution {
    public long count_common_characters(string[] strings) {
        if (strings.Length == 0) return 0;
        var common = new HashSet<char>(strings[0]);
        for (int i = 1; i < strings.Length; i++) {
            common.IntersectWith(new HashSet<char>(strings[i]));
        }
        return common.Count;
    }
}
