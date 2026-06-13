function longest_increasing_path(matrix: number[][]): number {
    if (!matrix || matrix.length === 0 || matrix[0].length === 0) return 0;
    const m = matrix.length;
    const n = matrix[0].length;
    const memo = new Map<string, number>();

    function dfs(r: number, c: number): number {
        const key = `${r},${c}`;
        if (memo.has(key)) return memo.get(key)!;
        let best = 1;
        for (const [dr, dc] of [[-1,0],[1,0],[0,-1],[0,1]]) {
            const nr = r + dr;
            const nc = c + dc;
            if (nr >= 0 && nr < m && nc >= 0 && nc < n && matrix[nr][nc] > matrix[r][c]) {
                best = Math.max(best, 1 + dfs(nr, nc));
            }
        }
        memo.set(key, best);
        return best;
    }

    let result = 0;
    for (let r = 0; r < m; r++) {
        for (let c = 0; c < n; c++) {
            result = Math.max(result, dfs(r, c));
        }
    }
    return result;
}
