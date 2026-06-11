use std::collections::BinaryHeap;
use std::cmp::Reverse;

fn swim_in_water(grid: Vec<Vec<i64>>) -> i64 {
    let n = grid.len();
    let mut visited = vec![vec![false; n]; n];
    let mut heap: BinaryHeap<Reverse<(i64, usize, usize)>> = BinaryHeap::new();
    heap.push(Reverse((grid[0][0], 0, 0)));
    visited[0][0] = true;
    let dirs: [(i64, i64); 4] = [(1, 0), (-1, 0), (0, 1), (0, -1)];
    while let Some(Reverse((t, i, j))) = heap.pop() {
        if i == n - 1 && j == n - 1 {
            return t;
        }
        for &(di, dj) in dirs.iter() {
            let ni = i as i64 + di;
            let nj = j as i64 + dj;
            if ni >= 0 && ni < n as i64 && nj >= 0 && nj < n as i64 {
                let (ui, uj) = (ni as usize, nj as usize);
                if !visited[ui][uj] {
                    visited[ui][uj] = true;
                    heap.push(Reverse((t.max(grid[ui][uj]), ui, uj)));
                }
            }
        }
    }
    -1
}
