use std::collections::VecDeque;

fn oranges_rotting(grid: Vec<Vec<i64>>) -> i64 {
    let mut grid = grid;
    let m = grid.len() as i64;
    let n = grid[0].len() as i64;
    let mut queue: VecDeque<(i64, i64, i64)> = VecDeque::new();
    let mut fresh = 0i64;
    for i in 0..m {
        for j in 0..n {
            let c = grid[i as usize][j as usize];
            if c == 2 {
                queue.push_back((i, j, 0));
            } else if c == 1 {
                fresh += 1;
            }
        }
    }
    if fresh == 0 {
        return 0;
    }
    let mut minutes = 0i64;
    let dirs = [(1i64, 0i64), (-1, 0), (0, 1), (0, -1)];
    while let Some((i, j, t)) = queue.pop_front() {
        for (di, dj) in dirs.iter() {
            let ni = i + di;
            let nj = j + dj;
            if ni >= 0 && ni < m && nj >= 0 && nj < n && grid[ni as usize][nj as usize] == 1 {
                grid[ni as usize][nj as usize] = 2;
                fresh -= 1;
                minutes = t + 1;
                queue.push_back((ni, nj, t + 1));
            }
        }
    }
    if fresh == 0 { minutes } else { -1 }
}
