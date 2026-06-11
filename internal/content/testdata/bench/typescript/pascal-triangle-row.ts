function pascal_row(n: number): number[] {
    function comb(n: number, k: number): number {
        if (k === 0 || k === n) return 1;
        let result = 1;
        for (let i = 0; i < k; i++) {
            result = result * (n - i) / (i + 1);
        }
        return Math.round(result);
    }
    const row: number[] = [];
    for (let k = 0; k < n; k++) {
        row.push(comb(n - 1, k));
    }
    return row;
}
