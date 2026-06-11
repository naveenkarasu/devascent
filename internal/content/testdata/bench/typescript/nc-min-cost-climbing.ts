function min_cost_climbing_stairs(cost: number[]): number {
    const n = cost.length;
    const c = cost.slice();
    for (let i = 2; i < n; i++) {
        c[i] += Math.min(c[i - 1], c[i - 2]);
    }
    return Math.min(c[n - 1], c[n - 2]);
}
