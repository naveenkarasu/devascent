class Solution {
    private int resStart = 0, resLen = 1;

    private void expand(String s, int left, int right) {
        while (left >= 0 && right < s.length() && s.charAt(left) == s.charAt(right)) {
            int length = right - left + 1;
            if (length > resLen) {
                resStart = left;
                resLen = length;
            }
            left--;
            right++;
        }
    }

    public String longest_palindrome(String s) {
        resStart = 0;
        resLen = 1;
        for (int i = 0; i < s.length(); i++) {
            expand(s, i, i);
            expand(s, i, i + 1);
        }
        return s.substring(resStart, resStart + resLen);
    }
}
