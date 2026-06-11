public class Solution {
    public long[] filter_even_digit_numbers(long[] nums) {
        var result = new List<long>();
        foreach (long x in nums) {
            if (DigitCount(x) % 2 == 0) result.Add(x);
        }
        return result.ToArray();
    }
    private int DigitCount(long n) {
        if (n == 0) return 1;
        int count = 0;
        while (n > 0) { count++; n /= 10; }
        return count;
    }
}
