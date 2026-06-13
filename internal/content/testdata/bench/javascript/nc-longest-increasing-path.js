function longest_increasing_path(matrix) {
    if (!matrix || !matrix[0]) return 0;
    const m = matrix.length, n = matrix[0].length;
    const memo = new Map();

    function dfs(r, c) {
        const key = r * n + c;
        if (memo.has(key)) return memo.get(key);
        let best = 1;
        const dirs = [[-1, 0], [1, 0], [0, -1], [0, 1]];
        for (const [dr, dc] of dirs) {
            const nr = r + dr, nc = c + dc;
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
