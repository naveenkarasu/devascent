use std::collections::BinaryHeap;
use std::cmp::Reverse;

fn k_closest(points: Vec<Vec<i64>>, k: i64) -> Vec<Vec<i64>> {
    let mut heap: BinaryHeap<Reverse<(i64, i64, i64)>> = BinaryHeap::new();
    for p in &points {
        let x = p[0];
        let y = p[1];
        let dist2 = x * x + y * y;
        heap.push(Reverse((dist2, x, y)));
    }
    let mut result: Vec<Vec<i64>> = Vec::new();
    for _ in 0..k {
        if let Some(Reverse((_d, x, y))) = heap.pop() {
            result.push(vec![x, y]);
        }
    }
    result
}
