import java.util.HashMap;
import java.util.Map;

class Solution {
    public long find_unique_among_k(long k, long[] arr) {
        Map<Long, Long> counts = new HashMap<>();
        for (long x : arr) {
            counts.put(x, counts.getOrDefault(x, 0L) + 1);
        }
        for (Map.Entry<Long, Long> e : counts.entrySet()) {
            if (e.getValue() % k != 0) return e.getKey();
        }
        return -1L;
    }
}
