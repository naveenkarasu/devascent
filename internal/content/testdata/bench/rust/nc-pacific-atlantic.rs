use std::collections::VecDeque;
use std::collections::HashSet;

fn bfs(heights: &Vec<Vec<i64>>, starts: &Vec<(i64, i64)>, m: i64, n: i64) -> HashSet<(i64, i64)> {
    let mut visited: HashSet<(i64, i64)> = HashSet::new();
    let mut queue: VecDeque<(i64, i64)> = VecDeque::new();
    for &s in starts {
        visited.insert(s);
        queue.push_back(s);
    }
    let dirs = [(1i64, 0i64), (-1, 0), (0, 1), (0, -1)];
    while let Some((i, j)) = queue.pop_front() {
        for (di, dj) in dirs.iter() {
            let ni = i + di;
            let nj = j + dj;
            if ni >= 0 && ni < m && nj >= 0 && nj < n
                && !visited.contains(&(ni, nj))
                && heights[ni as usize][nj as usize] >= heights[i as usize][j as usize]
            {
                visited.insert((ni, nj));
                queue.push_back((ni, nj));
            }
        }
    }
    visited
}

fn pacific_atlantic(heights: Vec<Vec<i64>>) -> Vec<Vec<i64>> {
    if heights.is_empty() {
        return Vec::new();
    }
    let m = heights.len() as i64;
    let n = heights[0].len() as i64;
    let mut pac_starts: Vec<(i64, i64)> = Vec::new();
    for j in 0..n {
        pac_starts.push((0, j));
    }
    for i in 1..m {
        pac_starts.push((i, 0));
    }
    let mut atl_starts: Vec<(i64, i64)> = Vec::new();
    for j in 0..n {
        atl_starts.push((m - 1, j));
    }
    for i in 0..(m - 1) {
        atl_starts.push((i, n - 1));
    }
    let pac = bfs(&heights, &pac_starts, m, n);
    let atl = bfs(&heights, &atl_starts, m, n);
    let mut both: Vec<(i64, i64)> = pac.intersection(&atl).cloned().collect();
    both.sort();
    both.into_iter().map(|(r, c)| vec![r, c]).collect()
}
