public class Solution {
    public long network_delay_time(long[][] times, long n, long k) {
        var graph = new Dictionary<int, List<(long weight, int nb)>>();
        foreach (var edge in times) {
            int u = (int)edge[0];
            if (!graph.ContainsKey(u)) graph[u] = new List<(long, int)>();
            graph[u].Add((edge[2], (int)edge[1]));
        }
        var dist = new Dictionary<int, long>();
        var heap = new PriorityQueue<(long d, int node), long>();
        heap.Enqueue((0, (int)k), 0);
        while (heap.Count > 0) {
            var (d, node) = heap.Dequeue();
            if (dist.ContainsKey(node)) continue;
            dist[node] = d;
            if (graph.ContainsKey(node)) {
                foreach (var (weight, nb) in graph[node]) {
                    if (!dist.ContainsKey(nb)) {
                        heap.Enqueue((d + weight, nb), d + weight);
                    }
                }
            }
        }
        if (dist.Count != n) return -1;
        return dist.Values.Max();
    }
}
