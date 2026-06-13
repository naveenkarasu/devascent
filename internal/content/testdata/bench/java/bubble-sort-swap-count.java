import java.util.*;

class Solution {
    public Map<String, Object> bubble_sort_info(long[] arr) {
        long[] a = arr.clone();
        long swaps = 0;
        boolean isSorted = false;
        while (!isSorted) {
            isSorted = true;
            for (int i = 0; i < a.length - 1; i++) {
                if (a[i] > a[i + 1]) {
                    long t = a[i];
                    a[i] = a[i + 1];
                    a[i + 1] = t;
                    swaps++;
                    isSorted = false;
                }
            }
        }
        Map<String, Object> out = new LinkedHashMap<>();
        out.put("swaps", swaps);
        out.put("first", a[0]);
        out.put("last", a[a.length - 1]);
        return out;
    }
}
