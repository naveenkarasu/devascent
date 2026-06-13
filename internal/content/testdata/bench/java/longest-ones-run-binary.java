class Solution {
    public long longest_ones_run(long n) {
        String binary = Long.toBinaryString(n);
        String[] parts = binary.split("0");
        long maxLen = 0;
        for (String p : parts) {
            if (p.length() > maxLen) maxLen = p.length();
        }
        return maxLen;
    }
}
