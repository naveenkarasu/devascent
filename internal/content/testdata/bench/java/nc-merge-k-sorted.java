import java.util.*;

class Solution {
    public long[] merge_k_lists(long[][] lists) {
        if (lists == null || lists.length == 0) return new long[0];
        // min-heap: [value, listIndex, elemIndex]
        PriorityQueue<long[]> heap = new PriorityQueue<>((a, b) -> Long.compare(a[0], b[0]));
        for (int i = 0; i < lists.length; i++) {
            if (lists[i] != null && lists[i].length > 0) {
                heap.offer(new long[]{lists[i][0], i, 0});
            }
        }
        List<Long> out = new ArrayList<>();
        while (!heap.isEmpty()) {
            long[] top = heap.poll();
            out.add(top[0]);
            int i = (int) top[1];
            int j = (int) top[2] + 1;
            if (j < lists[i].length) {
                heap.offer(new long[]{lists[i][j], i, j});
            }
        }
        long[] res = new long[out.size()];
        for (int i = 0; i < out.size(); i++) res[i] = out.get(i);
        return res;
    }
}
