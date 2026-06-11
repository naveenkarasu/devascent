function find_unique_among_k(k: number, arr: number[]): number {
    const counts: Map<number, number> = new Map();
    for (const val of arr) {
        counts.set(val, (counts.get(val) || 0) + 1);
    }
    for (const [val, cnt] of counts) {
        if (cnt % k !== 0) return val;
    }
    return -1;
}
