class Solution {
    public long[] update_pair(long a, long b) {
        return new long[]{a + b, Math.abs(a - b)};
    }
}
