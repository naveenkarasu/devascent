public class Solution {
    public bool can_finish(long num_courses, long[][] prerequisites) {
        int n = (int)num_courses;
        var graph = new List<List<int>>();
        for (int i = 0; i < n; i++) graph.Add(new List<int>());
        foreach (var pre in prerequisites) {
            graph[(int)pre[0]].Add((int)pre[1]);
        }
        int[] state = new int[n]; // 0=unvisited, 1=visiting, 2=done
        for (int i = 0; i < n; i++) {
            if (!Dfs(graph, state, i)) return false;
        }
        return true;
    }

    private bool Dfs(List<List<int>> graph, int[] state, int c) {
        if (state[c] == 1) return false;
        if (state[c] == 2) return true;
        state[c] = 1;
        foreach (int nxt in graph[c]) {
            if (!Dfs(graph, state, nxt)) return false;
        }
        state[c] = 2;
        return true;
    }
}
