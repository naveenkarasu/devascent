import java.util.*;

class GraphNode {
    int val;
    List<GraphNode> neighbors;
    GraphNode(int val) { this.val = val; this.neighbors = new ArrayList<>(); }
}

class Solution {
    private GraphNode cloneNode(GraphNode node, Map<GraphNode, GraphNode> mp) {
        if (node == null) return null;
        if (mp.containsKey(node)) return mp.get(node);
        GraphNode copy = new GraphNode(node.val);
        mp.put(node, copy);
        for (GraphNode nb : node.neighbors) {
            copy.neighbors.add(cloneNode(nb, mp));
        }
        return copy;
    }

    public long[][] clone_graph(long[][] adj) {
        if (adj == null || adj.length == 0) return new long[0][];
        int n = adj.length;
        Map<Integer, GraphNode> nodes = new HashMap<>();
        for (int i = 0; i < n; i++) nodes.put(i + 1, new GraphNode(i + 1));
        for (int i = 0; i < n; i++) {
            for (long nb : adj[i]) {
                nodes.get(i + 1).neighbors.add(nodes.get((int) nb));
            }
        }
        Map<GraphNode, GraphNode> mp = new HashMap<>();
        GraphNode cloned = cloneNode(nodes.get(1), mp);

        long[][] out = new long[n][];
        Set<Integer> seen = new HashSet<>();
        Deque<GraphNode> stack = new ArrayDeque<>();
        stack.push(cloned);
        while (!stack.isEmpty()) {
            GraphNode nd = stack.pop();
            if (seen.contains(nd.val)) continue;
            seen.add(nd.val);
            List<Integer> nbVals = new ArrayList<>();
            for (GraphNode nb : nd.neighbors) {
                nbVals.add(nb.val);
                stack.push(nb);
            }
            Collections.sort(nbVals);
            out[nd.val - 1] = new long[nbVals.size()];
            for (int i = 0; i < nbVals.size(); i++) out[nd.val - 1][i] = nbVals.get(i);
        }
        return out;
    }
}
