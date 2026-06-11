public class Solution {
    public long[] find_redundant_connection(long[][] edges) {
        int n = edges.Length;
        int[] parent = new int[n + 1];
        int[] rank = new int[n + 1];
        for (int i = 0; i <= n; i++) parent[i] = i;

        int Find(int x) {
            while (parent[x] != x) {
                parent[x] = parent[parent[x]];
                x = parent[x];
            }
            return x;
        }

        bool Union(int a, int b) {
            int pa = Find(a), pb = Find(b);
            if (pa == pb) return false;
            if (rank[pa] < rank[pb]) { int t = pa; pa = pb; pb = t; }
            parent[pb] = pa;
            if (rank[pa] == rank[pb]) rank[pa]++;
            return true;
        }

        foreach (var e in edges) {
            int u = (int)e[0], v = (int)e[1];
            if (!Union(u, v)) return new long[] { e[0], e[1] };
        }
        return System.Array.Empty<long>();
    }
}
