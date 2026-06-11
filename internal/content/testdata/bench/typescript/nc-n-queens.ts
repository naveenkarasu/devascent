function total_n_queens(n: number): number {
    let count = 0;
    const cols: Set<number> = new Set();
    const diag1: Set<number> = new Set(); // row - col
    const diag2: Set<number> = new Set(); // row + col

    function backtrack(row: number): void {
        if (row === n) {
            count++;
            return;
        }
        for (let col = 0; col < n; col++) {
            if (cols.has(col) || diag1.has(row - col) || diag2.has(row + col)) continue;
            cols.add(col);
            diag1.add(row - col);
            diag2.add(row + col);
            backtrack(row + 1);
            cols.delete(col);
            diag1.delete(row - col);
            diag2.delete(row + col);
        }
    }

    backtrack(0);
    return count;
}
