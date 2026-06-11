public class Solution {
    public long[] selection_sort_k_passes(long[] arr, long k) {
        long[] result = (long[])arr.Clone();
        int n = result.Length;
        int passes = (int)System.Math.Min(k, n - 1);
        for (int i = 0; i < passes; i++) {
            int minIdx = i;
            for (int j = i + 1; j < n; j++) {
                if (result[j] < result[minIdx]) minIdx = j;
            }
            long tmp = result[i];
            result[i] = result[minIdx];
            result[minIdx] = tmp;
        }
        return result;
    }
}
