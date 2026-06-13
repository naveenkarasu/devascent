fn count_grid_colorings(n: i64, m: i64) -> i64 {
    const MOD: i64 = 1_000_000_007;
    let n = n as usize;
    let m = m as usize;
    let mut grid: Vec<Vec<i64>> = vec![vec![0; m]; n];

    fn backtrack(pos: usize, n: usize, m: usize, grid: &mut Vec<Vec<i64>>) -> i64 {
        const MOD: i64 = 1_000_000_007;
        if pos == n * m {
            return 1;
        }
        let r = pos / m;
        let c = pos % m;
        let mut total: i64 = 0;
        for color in 1..=3i64 {
            if (r == 0 || grid[r - 1][c] != color) && (c == 0 || grid[r][c - 1] != color) {
                grid[r][c] = color;
                total = (total + backtrack(pos + 1, n, m, grid)) % MOD;
                grid[r][c] = 0;
            }
        }
        total
    }

    if n == 0 || m == 0 {
        return 1;
    }
    backtrack(0, n, m, &mut grid)
}
