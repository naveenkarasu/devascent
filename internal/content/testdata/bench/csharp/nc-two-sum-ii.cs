public class Solution {
    public long[] two_sum_ii(long[] numbers, long target) {
        int i = 0, j = numbers.Length - 1;
        while (i < j) {
            long total = numbers[i] + numbers[j];
            if (total == target) return new long[] { i + 1, j + 1 };
            if (total < target) i++;
            else j--;
        }
        return new long[] {};
    }
}
