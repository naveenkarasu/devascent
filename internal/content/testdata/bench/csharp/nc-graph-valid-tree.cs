using System.Collections.Generic;

public class Solution {
    public bool valid_tree(long n, long[][] edges) {
        if (edges.Length != n - 1) return false;
        int N = (int)n;
        var adj = new List<List<int>>();
        for (int i = 0; i < N; i++) adj.Add(new List<int>());
        foreach (long[] e in edges) {
            int u = (int)e[0], v = (int)e[1];
            adj[u].Add(v);
            adj[v].Add(u);
        }
        bool[] visited = new bool[N];
        Dfs(0, adj, visited);
        foreach (bool v in visited) {
            if (!v) return false;
        }
        return true;
    }

    private void Dfs(int node, List<List<int>> adj, bool[] visited) {
        visited[node] = true;
        foreach (int nb in adj[node]) {
            if (!visited[nb]) Dfs(nb, adj, visited);
        }
    }
}
