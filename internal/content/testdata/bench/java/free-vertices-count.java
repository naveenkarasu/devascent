import java.util.*;
class Solution {
    public long free_vertices(long n, long m) {
        long k = 0;
        while (k * (k - 1) / 2 < m) {
            k++;
        }
        return n - k;
    }
}
