function find_redundant_connection(edges: number[][]): number[] {
    const n = edges.length;
    const parent: number[] = Array.from({length: n + 1}, (_, i) => i);
    const rank: number[] = new Array(n + 1).fill(0);

    function find(x: number): number {
        while (parent[x] !== x) {
            parent[x] = parent[parent[x]];
            x = parent[x];
        }
        return x;
    }

    function union(a: number, b: number): boolean {
        let pa = find(a);
        let pb = find(b);
        if (pa === pb) return false;
        if (rank[pa] < rank[pb]) {
            const tmp = pa; pa = pb; pb = tmp;
        }
        parent[pb] = pa;
        if (rank[pa] === rank[pb]) rank[pa]++;
        return true;
    }

    for (const [u, v] of edges) {
        if (!union(u, v)) return [u, v];
    }
    return [];
}
