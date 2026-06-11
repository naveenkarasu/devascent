public class Solution {
    public long get_sum(long a, long b) {
        long mask = 0xFFFFFFFFL;
        while ((b & mask) != 0) {
            long carry = ((a & b) << 1) & mask;
            a = (a ^ b) & mask;
            b = carry;
        }
        a &= mask;
        return a <= 0x7FFFFFFFL ? a : ~(a ^ mask);
    }
}
