public class Solution {
    public long min_cost_connect_points(long[][] points) {
        int n = points.Length;
        bool[] visited = new bool[n];
        // min-heap: (cost, point_index)
        var heap = new PriorityQueue<(long cost, int idx), long>();
        heap.Enqueue((0, 0), 0);
        long total = 0;
        int count = 0;
        while (count < n) {
            var (cost, i) = heap.Dequeue();
            if (visited[i]) continue;
            visited[i] = true;
            total += cost;
            count++;
            long xi = points[i][0], yi = points[i][1];
            for (int j = 0; j < n; j++) {
                if (!visited[j]) {
                    long dist = Math.Abs(xi - points[j][0]) + Math.Abs(yi - points[j][1]);
                    heap.Enqueue((dist, j), dist);
                }
            }
        }
        return total;
    }
}
