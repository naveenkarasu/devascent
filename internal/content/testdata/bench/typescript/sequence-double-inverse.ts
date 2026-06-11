function sequence_equation(p: number[]): number[] {
    const n = p.length;
    const inv: number[] = new Array(n + 1).fill(0);
    for (let i = 0; i < n; i++) {
        inv[p[i]] = i + 1;
    }
    const result: number[] = [];
    for (let x = 1; x <= n; x++) {
        result.push(inv[inv[x]]);
    }
    return result;
}
