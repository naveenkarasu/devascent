public class Solution {
    public long[] plus_one(long[] digits) {
        for (int i = digits.Length - 1; i >= 0; i--) {
            if (digits[i] < 9) {
                digits[i]++;
                return digits;
            }
            digits[i] = 0;
        }
        long[] result = new long[digits.Length + 1];
        result[0] = 1;
        return result;
    }
}
