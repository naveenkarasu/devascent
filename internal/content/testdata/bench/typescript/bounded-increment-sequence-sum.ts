function bounded_sequence_sum(n: number, b: number, x: number, y: number): number {
    const a: number[] = new Array(n + 1).fill(0);
    for (let i = 1; i <= n; i++) {
        if (a[i - 1] + x <= b) {
            a[i] = a[i - 1] + x;
        } else {
            a[i] = a[i - 1] - y;
        }
    }
    return a.reduce((sum, val) => sum + val, 0);
}
