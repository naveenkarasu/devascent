function oranges_rotting(grid: number[][]): number {
    const m = grid.length;
    const n = grid[0].length;
    const queue: [number, number, number][] = [];
    let fresh = 0;
    for (let i = 0; i < m; i++) {
        for (let j = 0; j < n; j++) {
            if (grid[i][j] === 2) queue.push([i, j, 0]);
            else if (grid[i][j] === 1) fresh++;
        }
    }
    if (fresh === 0) return 0;
    let minutes = 0;
    let qi = 0;
    while (qi < queue.length) {
        const [i, j, t] = queue[qi++];
        for (const [di, dj] of [[1,0],[-1,0],[0,1],[0,-1]]) {
            const ni = i + di, nj = j + dj;
            if (ni >= 0 && ni < m && nj >= 0 && nj < n && grid[ni][nj] === 1) {
                grid[ni][nj] = 2;
                fresh--;
                minutes = t + 1;
                queue.push([ni, nj, t + 1]);
            }
        }
    }
    return fresh === 0 ? minutes : -1;
}
