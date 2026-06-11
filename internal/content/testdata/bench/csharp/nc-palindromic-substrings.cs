public class Solution {
    public long count_substrings(string s) {
        long count = 0;

        void Expand(int left, int right) {
            while (left >= 0 && right < s.Length && s[left] == s[right]) {
                count++;
                left--;
                right++;
            }
        }

        for (int i = 0; i < s.Length; i++) {
            Expand(i, i);
            Expand(i, i + 1);
        }

        return count;
    }
}
