function valid_tree(n: number, edges: number[][]): boolean {
    if (edges.length !== n - 1) return false;
    const adj: number[][] = Array.from({length: n}, () => []);
    for (const [u, v] of edges) {
        adj[u].push(v);
        adj[v].push(u);
    }
    const visited = new Set<number>();

    function dfs(node: number): void {
        visited.add(node);
        for (const nb of adj[node]) {
            if (!visited.has(nb)) {
                dfs(nb);
            }
        }
    }

    dfs(0);
    return visited.size === n;
}
