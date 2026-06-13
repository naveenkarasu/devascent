class Solution {
    public long off_bulbs_between_lit(long[] bulbs) {
        int n = bulbs.length;
        int firstOn = -1;
        int lastOn = -1;
        for (int i = 0; i < n; i++) {
            if (bulbs[i] == 1) {
                firstOn = i;
                break;
            }
        }
        for (int i = n - 1; i >= 0; i--) {
            if (bulbs[i] == 1) {
                lastOn = i;
                break;
            }
        }
        if (firstOn == -1) return 0;
        long count = 0;
        for (int j = firstOn; j <= lastOn; j++) {
            if (bulbs[j] == 0) count++;
        }
        return count;
    }
}
