class Solution {
    private long[] parent;
    private long[] rank;

    public long[] find_redundant_connection(long[][] edges) {
        int n = edges.length;
        parent = new long[n + 1];
        rank = new long[n + 1];
        for (int i = 0; i <= n; i++) parent[i] = i;
        for (long[] edge : edges) {
            int u = (int) edge[0], v = (int) edge[1];
            if (!union(u, v)) {
                return new long[]{u, v};
            }
        }
        return new long[0];
    }

    private long find(long x) {
        while (parent[(int)x] != x) {
            parent[(int)x] = parent[(int)parent[(int)x]];
            x = parent[(int)x];
        }
        return x;
    }

    private boolean union(int a, int b) {
        long pa = find(a), pb = find(b);
        if (pa == pb) return false;
        if (rank[(int)pa] < rank[(int)pb]) {
            long tmp = pa; pa = pb; pb = tmp;
        }
        parent[(int)pb] = pa;
        if (rank[(int)pa] == rank[(int)pb]) rank[(int)pa]++;
        return true;
    }
}
