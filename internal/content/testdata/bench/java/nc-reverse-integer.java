class Solution {
    public long reverse_integer(long x) {
        long INT_MIN = -(1L << 31);
        long INT_MAX = (1L << 31) - 1;
        long sign = x < 0 ? -1 : 1;
        String digits = new StringBuilder(Long.toString(Math.abs(x))).reverse().toString();
        long result = sign * Long.parseLong(digits);
        if (result < INT_MIN || result > INT_MAX) return 0;
        return result;
    }
}
