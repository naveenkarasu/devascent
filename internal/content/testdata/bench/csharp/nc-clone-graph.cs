public class Solution {
    private class Node {
        public int val;
        public List<Node> neighbors;
        public Node(int v) { val = v; neighbors = new List<Node>(); }
    }

    private Node Clone(Node node, Dictionary<int, Node> mp) {
        if (node == null) return null;
        if (mp.ContainsKey(node.val)) return mp[node.val];
        var copy = new Node(node.val);
        mp[node.val] = copy;
        foreach (var nb in node.neighbors)
            copy.neighbors.Add(Clone(nb, mp));
        return copy;
    }

    public long[][] clone_graph(long[][] adj) {
        if (adj == null || adj.Length == 0) return new long[0][];
        int n = adj.Length;
        var nodes = new Node[n + 1];
        for (int i = 1; i <= n; i++) nodes[i] = new Node(i);
        for (int i = 0; i < n; i++)
            foreach (long nb in adj[i])
                nodes[i + 1].neighbors.Add(nodes[(int)nb]);

        var mp = new Dictionary<int, Node>();
        var cloned = Clone(nodes[1], mp);

        // Re-serialize
        var out_adj = new long[n][];
        var seen = new HashSet<int>();
        var stack = new Stack<Node>();
        stack.Push(cloned);
        while (stack.Count > 0) {
            var nd = stack.Pop();
            if (seen.Contains(nd.val)) continue;
            seen.Add(nd.val);
            var nbVals = nd.neighbors.Select(nb => (long)nb.val).OrderBy(x => x).ToArray();
            out_adj[nd.val - 1] = nbVals;
            foreach (var nb in nd.neighbors) stack.Push(nb);
        }
        return out_adj;
    }
}
