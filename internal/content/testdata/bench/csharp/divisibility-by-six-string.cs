public class Solution {
    public bool is_divisible_by_six(string s) {
        int lastDigit = s[s.Length - 1] - '0';
        if (lastDigit % 2 != 0) return false;
        long digitSum = 0;
        foreach (var ch in s) {
            digitSum += ch - '0';
        }
        return digitSum % 3 == 0;
    }
}
