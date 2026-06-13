using System.Collections.Generic;
public class Solution {
    public Dictionary<string, long> bubble_sort_info(long[] arr) {
        long[] a = (long[])arr.Clone();
        long swaps = 0;
        bool isSorted = false;
        while (!isSorted) {
            isSorted = true;
            for (int i = 0; i < a.Length - 1; i++) {
                if (a[i] > a[i + 1]) {
                    long t = a[i]; a[i] = a[i + 1]; a[i + 1] = t;
                    swaps++;
                    isSorted = false;
                }
            }
        }
        return new Dictionary<string, long> {
            { "swaps", swaps },
            { "first", a[0] },
            { "last", a[a.Length - 1] }
        };
    }
}
