function clone_graph(adj: number[][]): number[][] {
    if (!adj || adj.length === 0) return [];

    // Build original nodes
    const nodes: any[] = [];
    for (let i = 0; i < adj.length; i++) {
        nodes.push({ val: i + 1, neighbors: [] });
    }
    for (let i = 0; i < adj.length; i++) {
        nodes[i].neighbors = adj[i].map(x => nodes[x - 1]);
    }

    // Clone using DFS with map
    const mp = new Map<any, any>();
    function dfs(n: any): any {
        if (mp.has(n)) return mp.get(n);
        const copy = { val: n.val, neighbors: [] as any[] };
        mp.set(n, copy);
        for (const nb of n.neighbors) {
            copy.neighbors.push(dfs(nb));
        }
        return copy;
    }

    const cloned = dfs(nodes[0]);

    // Serialize clone back to adjacency list
    const out: number[][] = Array.from({ length: adj.length }, () => []);
    const seen = new Set<number>();
    const stack = [cloned];
    while (stack.length > 0) {
        const nd = stack.pop();
        if (seen.has(nd.val)) continue;
        seen.add(nd.val);
        out[nd.val - 1] = nd.neighbors.map((n: any) => n.val).sort((a: number, b: number) => a - b);
        for (const nb of nd.neighbors) {
            stack.push(nb);
        }
    }
    return out;
}
