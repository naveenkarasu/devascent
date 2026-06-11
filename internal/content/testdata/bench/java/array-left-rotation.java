class Solution {
    public long[] left_rotate(long[] arr, long k) {
        int n = arr.length;
        if (n == 0) return arr;
        int shift = (int)(k % n);
        long[] result = new long[n];
        for (int i = 0; i < n; i++) {
            result[i] = arr[(i + shift) % n];
        }
        return result;
    }
}
