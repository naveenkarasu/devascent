class Solution {
    public boolean is_divisible_by_six(String s) {
        char lastChar = s.charAt(s.length() - 1);
        if ((lastChar - '0') % 2 != 0) return false;
        long digitSum = 0;
        for (char ch : s.toCharArray()) {
            digitSum += (ch - '0');
        }
        return digitSum % 3 == 0;
    }
}
