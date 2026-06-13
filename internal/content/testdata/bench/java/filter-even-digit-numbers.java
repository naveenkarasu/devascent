class Solution {
    public long[] filter_even_digit_numbers(long[] nums) {
        java.util.List<Long> result = new java.util.ArrayList<>();
        for (long n : nums) {
            if (digitCount(n) % 2 == 0) {
                result.add(n);
            }
        }
        long[] out = new long[result.size()];
        for (int i = 0; i < out.length; i++) out[i] = result.get(i);
        return out;
    }

    private int digitCount(long n) {
        if (n == 0) return 1;
        int count = 0;
        while (n > 0) {
            count++;
            n /= 10;
        }
        return count;
    }
}
