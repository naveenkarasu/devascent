using System.Collections.Generic;

public class Solution {
    public long[] find_order(long num_courses, long[][] prerequisites) {
        int n = (int)num_courses;
        long[] inDegree = new long[n];
        var adj = new List<int>[n];
        for (int i = 0; i < n; i++) adj[i] = new List<int>();

        foreach (var pre in prerequisites) {
            int a = (int)pre[0], b = (int)pre[1];
            adj[b].Add(a);
            inDegree[a]++;
        }

        var heap = new PriorityQueue<long, long>();
        for (int i = 0; i < n; i++) {
            if (inDegree[i] == 0) heap.Enqueue(i, i);
        }

        var order = new List<long>();
        while (heap.Count > 0) {
            long node = heap.Dequeue();
            order.Add(node);
            foreach (int nei in adj[(int)node]) {
                inDegree[nei]--;
                if (inDegree[nei] == 0) heap.Enqueue(nei, nei);
            }
        }

        return order.Count == n ? order.ToArray() : System.Array.Empty<long>();
    }
}
