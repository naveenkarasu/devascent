function clone_graph(adj) {
    if (!adj || adj.length === 0) return [];
    // Build adjacency list output: BFS from node 1
    const out = Array.from({ length: adj.length }, () => []);
    const seen = new Set();
    const stack = [1];
    while (stack.length > 0) {
        const val = stack.pop();
        if (seen.has(val)) continue;
        seen.add(val);
        const nbrs = adj[val - 1];
        out[val - 1] = nbrs.slice().sort((a, b) => a - b);
        for (const nb of nbrs) {
            stack.push(nb);
        }
    }
    return out;
}
