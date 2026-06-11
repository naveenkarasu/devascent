public class Solution {
    public long trap_node(long n, long[][] edges) {
        var degree = new long[n];
        foreach (var edge in edges) {
            long u = edge[0];
            long v = edge[1];
            degree[u]++;
            degree[v]++;
        }
        long maxDeg = 0;
        foreach (var d in degree) {
            if (d > maxDeg) maxDeg = d;
        }
        return maxDeg - 1;
    }
}
