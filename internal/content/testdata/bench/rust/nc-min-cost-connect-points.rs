use std::collections::BinaryHeap;
use std::cmp::Reverse;

fn min_cost_connect_points(points: Vec<Vec<i64>>) -> i64 {
    let n = points.len();
    let mut visited = vec![false; n];
    let mut heap: BinaryHeap<Reverse<(i64, usize)>> = BinaryHeap::new();
    heap.push(Reverse((0, 0)));
    let mut total = 0i64;
    let mut count = 0;
    while count < n {
        let Reverse((cost, i)) = match heap.pop() {
            Some(x) => x,
            None => break,
        };
        if visited[i] {
            continue;
        }
        visited[i] = true;
        count += 1;
        total += cost;
        let (xi, yi) = (points[i][0], points[i][1]);
        for j in 0..n {
            if !visited[j] {
                let (xj, yj) = (points[j][0], points[j][1]);
                let d = (xi - xj).abs() + (yi - yj).abs();
                heap.push(Reverse((d, j)));
            }
        }
    }
    total
}
