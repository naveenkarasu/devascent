import java.util.*;

class Solution {
    public long count_even_sum_triples(long[] arr, long l, long r) {
        int li = (int)(l - 1);
        int ri = (int)(r - 1);
        long e = 0;
        for (int i = li; i <= ri; i++) {
            if (arr[i] % 2 == 0) e++;
        }
        long o = (ri - li + 1) - e;
        return (e * (e - 1) * (e - 2)) / 6 + (o * (o - 1) / 2) * e;
    }
}
