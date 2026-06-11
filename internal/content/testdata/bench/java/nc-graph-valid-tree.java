import java.util.*;

class Solution {
    public boolean valid_tree(long n, long[][] edges) {
        if (edges.length != n - 1) return false;
        int N = (int) n;
        List<List<Integer>> adj = new ArrayList<>();
        for (int i = 0; i < N; i++) adj.add(new ArrayList<>());
        for (long[] e : edges) {
            int u = (int) e[0], v = (int) e[1];
            adj.get(u).add(v);
            adj.get(v).add(u);
        }
        boolean[] visited = new boolean[N];
        dfs(0, adj, visited);
        for (boolean v : visited) {
            if (!v) return false;
        }
        return true;
    }

    private void dfs(int node, List<List<Integer>> adj, boolean[] visited) {
        visited[node] = true;
        for (int nb : adj.get(node)) {
            if (!visited[nb]) dfs(nb, adj, visited);
        }
    }
}
