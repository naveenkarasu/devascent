function set_union(a: number[], b: number[]): number[] {
    const seen = new Set([...a, ...b]);
    return Array.from(seen).sort((x, y) => x - y);
}
