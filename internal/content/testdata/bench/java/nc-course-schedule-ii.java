import java.util.*;

class Solution {
    public long[] find_order(long num_courses, long[][] prerequisites) {
        int n = (int) num_courses;
        long[] inDegree = new long[n];
        List<List<Integer>> adj = new ArrayList<>();
        for (int i = 0; i < n; i++) adj.add(new ArrayList<>());
        for (long[] pre : prerequisites) {
            int a = (int) pre[0], b = (int) pre[1];
            adj.get(b).add(a);
            inDegree[a]++;
        }
        PriorityQueue<Integer> heap = new PriorityQueue<>();
        for (int i = 0; i < n; i++) {
            if (inDegree[i] == 0) heap.add(i);
        }
        long[] order = new long[n];
        int idx = 0;
        while (!heap.isEmpty()) {
            int node = heap.poll();
            order[idx++] = node;
            for (int nei : adj.get(node)) {
                inDegree[nei]--;
                if (inDegree[nei] == 0) heap.add(nei);
            }
        }
        if (idx == n) return order;
        return new long[0];
    }
}
