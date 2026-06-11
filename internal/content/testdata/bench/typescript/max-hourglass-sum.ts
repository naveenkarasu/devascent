function max_hourglass_sum(grid: number[][]): number {
    let best = -Infinity;
    for (let i = 0; i < 4; i++) {
        for (let j = 0; j < 4; j++) {
            const s = (grid[i][j] + grid[i][j+1] + grid[i][j+2]
                + grid[i+1][j+1]
                + grid[i+2][j] + grid[i+2][j+1] + grid[i+2][j+2]);
            if (s > best) best = s;
        }
    }
    return best;
}
