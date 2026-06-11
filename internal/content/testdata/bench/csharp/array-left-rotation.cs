public class Solution {
    public long[] left_rotate(long[] arr, long k) {
        int n = arr.Length;
        if (n == 0) return arr;
        int ki = (int)(k % n);
        long[] result = new long[n];
        for (int i = 0; i < n; i++) {
            result[i] = arr[(i + ki) % n];
        }
        return result;
    }
}
