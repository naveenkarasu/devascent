function trap_node(n: number, edges: number[][]): number {
    const degree: number[] = new Array(n).fill(0);
    for (const [u, v] of edges) {
        degree[u]++;
        degree[v]++;
    }
    return Math.max(...degree) - 1;
}
