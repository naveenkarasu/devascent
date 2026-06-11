public class Solution {
    public long count_vowels(string s) {
        const string vowels = "aeiouAEIOU";
        long count = 0;
        foreach (var ch in s) {
            if (vowels.IndexOf(ch) >= 0) count++;
        }
        return count;
    }
}
