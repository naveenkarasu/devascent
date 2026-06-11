public class Solution {
    public string longest_palindrome(string s) {
        int resStart = 0, resLen = 1;

        void Expand(int left, int right) {
            while (left >= 0 && right < s.Length && s[left] == s[right]) {
                int length = right - left + 1;
                if (length > resLen) {
                    resStart = left;
                    resLen = length;
                }
                left--;
                right++;
            }
        }

        for (int i = 0; i < s.Length; i++) {
            Expand(i, i);
            Expand(i, i + 1);
        }

        return s.Substring(resStart, resLen);
    }
}
