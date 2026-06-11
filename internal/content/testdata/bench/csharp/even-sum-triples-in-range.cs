public class Solution {
    public long count_even_sum_triples(long[] arr, long l, long r) {
        // subarray is arr[l-1..r-1] (0-indexed)
        long e = 0, o = 0;
        for (long i = l - 1; i < r; i++) {
            if (arr[i] % 2 == 0) e++;
            else o++;
        }
        return (e * (e - 1) * (e - 2)) / 6 + (o * (o - 1) / 2) * e;
    }
}
