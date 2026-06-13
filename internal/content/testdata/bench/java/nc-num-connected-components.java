import java.util.*;

class Solution {
    public long count_components(long n, long[][] edges) {
        int N = (int) n;
        int[] parent = new int[N];
        for (int i = 0; i < N; i++) parent[i] = i;

        long components = n;
        for (long[] e : edges) {
            int pa = find(parent, (int) e[0]);
            int pb = find(parent, (int) e[1]);
            if (pa != pb) {
                parent[pa] = pb;
                components--;
            }
        }
        return components;
    }

    private int find(int[] parent, int x) {
        while (parent[x] != x) {
            parent[x] = parent[parent[x]];
            x = parent[x];
        }
        return x;
    }
}
