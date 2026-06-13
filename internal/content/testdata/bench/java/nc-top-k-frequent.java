import java.util.*;

class Solution {
    public long[] top_k_frequent(long[] nums, long k) {
        Map<Long, Integer> freq = new HashMap<>();
        for (long n : nums) freq.merge(n, 1, Integer::sum);
        List<Long> keys = new ArrayList<>(freq.keySet());
        keys.sort((a, b) -> {
            int cmp = freq.get(b) - freq.get(a);
            if (cmp != 0) return cmp;
            return Long.compare(a, b);
        });
        List<Long> top = keys.subList(0, (int) k);
        Collections.sort(top);
        long[] res = new long[top.size()];
        for (int i = 0; i < top.size(); i++) res[i] = top.get(i);
        return res;
    }
}
