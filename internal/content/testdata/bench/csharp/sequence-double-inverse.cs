public class Solution {
    public long[] sequence_equation(long[] p) {
        int n = p.Length;
        long[] inv = new long[n + 1];
        for (int i = 0; i < n; i++) {
            inv[p[i]] = i + 1;
        }
        long[] result = new long[n];
        for (int x = 1; x <= n; x++) {
            result[x - 1] = inv[inv[x]];
        }
        return result;
    }
}
