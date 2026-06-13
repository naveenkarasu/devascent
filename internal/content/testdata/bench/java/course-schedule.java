import java.util.*;

class Solution {
    public boolean can_finish(long num_courses, long[][] prerequisites) {
        int n = (int) num_courses;
        List<List<Integer>> graph = new ArrayList<>();
        for (int i = 0; i < n; i++) graph.add(new ArrayList<>());
        for (long[] pre : prerequisites) {
            graph.get((int) pre[0]).add((int) pre[1]);
        }
        int[] state = new int[n]; // 0=unvisited, 1=visiting, 2=done

        for (int i = 0; i < n; i++) {
            if (!dfs(graph, state, i)) return false;
        }
        return true;
    }
    private boolean dfs(List<List<Integer>> graph, int[] state, int c) {
        if (state[c] == 1) return false;
        if (state[c] == 2) return true;
        state[c] = 1;
        for (int nxt : graph.get(c)) {
            if (!dfs(graph, state, nxt)) return false;
        }
        state[c] = 2;
        return true;
    }
}
