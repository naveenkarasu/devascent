import java.util.*;

class Solution {
    public long network_delay_time(long[][] times, long n, long k) {
        Map<Integer, List<long[]>> graph = new HashMap<>();
        for (long[] edge : times) {
            int u = (int) edge[0];
            graph.computeIfAbsent(u, x -> new ArrayList<>()).add(new long[]{edge[2], edge[1]});
        }
        Map<Integer, Long> dist = new HashMap<>();
        PriorityQueue<long[]> heap = new PriorityQueue<>(Comparator.comparingLong(a -> a[0]));
        heap.offer(new long[]{0, k});
        while (!heap.isEmpty()) {
            long[] curr = heap.poll();
            long d = curr[0];
            int node = (int) curr[1];
            if (dist.containsKey(node)) continue;
            dist.put(node, d);
            for (long[] nb : graph.getOrDefault(node, new ArrayList<>())) {
                int nbNode = (int) nb[1];
                if (!dist.containsKey(nbNode)) {
                    heap.offer(new long[]{d + nb[0], nbNode});
                }
            }
        }
        if (dist.size() != n) return -1;
        long max = 0;
        for (long v : dist.values()) max = Math.max(max, v);
        return max;
    }
}
