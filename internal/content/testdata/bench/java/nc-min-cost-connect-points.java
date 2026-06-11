import java.util.*;

class Solution {
    public long min_cost_connect_points(long[][] points) {
        int n = points.length;
        boolean[] visited = new boolean[n];
        // min-heap: (cost, point_index)
        PriorityQueue<long[]> heap = new PriorityQueue<>(Comparator.comparingLong(a -> a[0]));
        heap.offer(new long[]{0, 0});
        long total = 0;
        int count = 0;
        while (count < n) {
            long[] curr = heap.poll();
            long cost = curr[0];
            int i = (int) curr[1];
            if (visited[i]) continue;
            visited[i] = true;
            total += cost;
            count++;
            long xi = points[i][0], yi = points[i][1];
            for (int j = 0; j < n; j++) {
                if (!visited[j]) {
                    long dist = Math.abs(xi - points[j][0]) + Math.abs(yi - points[j][1]);
                    heap.offer(new long[]{dist, j});
                }
            }
        }
        return total;
    }
}
