function max_area_of_island(grid) {
    if (!grid || grid.length === 0) return 0;
    const m = grid.length, n = grid[0].length;
    const seen = new Set();

    function dfs(i, j) {
        if (i < 0 || j < 0 || i >= m || j >= n) return 0;
        if (grid[i][j] === 0 || seen.has(i * n + j)) return 0;
        seen.add(i * n + j);
        return 1 + dfs(i+1,j) + dfs(i-1,j) + dfs(i,j+1) + dfs(i,j-1);
    }

    let best = 0;
    for (let i = 0; i < m; i++) {
        for (let j = 0; j < n; j++) {
            if (grid[i][j] === 1 && !seen.has(i * n + j)) {
                best = Math.max(best, dfs(i, j));
            }
        }
    }
    return best;
}
