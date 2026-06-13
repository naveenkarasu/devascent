public class Solution {
    public long[] insertion_sort(long[] arr) {
        long[] a = (long[])arr.Clone();
        int n = a.Length;
        for (int i = 1; i < n; i++) {
            long key = a[i];
            int j = i - 1;
            while (j >= 0 && a[j] > key) {
                a[j + 1] = a[j];
                j--;
            }
            a[j + 1] = key;
        }
        return a;
    }
}
