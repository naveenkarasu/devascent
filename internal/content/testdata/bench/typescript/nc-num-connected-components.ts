function count_components(n: number, edges: number[][]): number {
    const parent = Array.from({length: n}, (_, i) => i);

    function find(x: number): number {
        while (parent[x] !== x) {
            parent[x] = parent[parent[x]];
            x = parent[x];
        }
        return x;
    }

    function union(a: number, b: number): number {
        const pa = find(a), pb = find(b);
        if (pa === pb) return 0;
        parent[pa] = pb;
        return 1;
    }

    let components = n;
    for (const [u, v] of edges) {
        components -= union(u, v);
    }
    return components;
}
