import java.util.*;

class Solution {
    public long[][] merge_intervals(long[][] intervals) {
        Arrays.sort(intervals, (a, b) -> Long.compare(a[0], b[0]));
        List<long[]> res = new ArrayList<>();
        for (long[] iv : intervals) {
            if (!res.isEmpty() && iv[0] <= res.get(res.size() - 1)[1]) {
                res.get(res.size() - 1)[1] = Math.max(res.get(res.size() - 1)[1], iv[1]);
            } else {
                res.add(new long[]{iv[0], iv[1]});
            }
        }
        return res.toArray(new long[0][]);
    }
}
