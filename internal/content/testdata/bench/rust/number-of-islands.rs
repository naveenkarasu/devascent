fn dfs(grid: &Vec<Vec<i64>>, seen: &mut Vec<Vec<bool>>, i: i64, j: i64, m: i64, n: i64) {
    if i < 0 || j < 0 || i >= m || j >= n {
        return;
    }
    let (ui, uj) = (i as usize, j as usize);
    if grid[ui][uj] == 0 || seen[ui][uj] {
        return;
    }
    seen[ui][uj] = true;
    dfs(grid, seen, i + 1, j, m, n);
    dfs(grid, seen, i - 1, j, m, n);
    dfs(grid, seen, i, j + 1, m, n);
    dfs(grid, seen, i, j - 1, m, n);
}

fn num_islands(grid: Vec<Vec<i64>>) -> i64 {
    if grid.is_empty() {
        return 0;
    }
    let m = grid.len();
    let n = grid[0].len();
    let mut seen = vec![vec![false; n]; m];
    let mut count = 0i64;
    for i in 0..m {
        for j in 0..n {
            if grid[i][j] == 1 && !seen[i][j] {
                count += 1;
                dfs(&grid, &mut seen, i as i64, j as i64, m as i64, n as i64);
            }
        }
    }
    count
}
