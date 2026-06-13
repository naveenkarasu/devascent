use std::collections::BinaryHeap;
use std::cmp::Reverse;

fn min_interval(intervals: Vec<Vec<i64>>, queries: Vec<i64>) -> Vec<i64> {
    let mut intervals = intervals;
    intervals.sort();
    let mut indexed: Vec<(usize, i64)> = queries.iter().cloned().enumerate().collect();
    indexed.sort_by_key(|x| x.1);
    let mut res = vec![-1i64; queries.len()];
    // min-heap of (size, end)
    let mut heap: BinaryHeap<Reverse<(i64, i64)>> = BinaryHeap::new();
    let mut i = 0usize;
    for &(orig_idx, q) in indexed.iter() {
        while i < intervals.len() && intervals[i][0] <= q {
            let l = intervals[i][0];
            let r = intervals[i][1];
            heap.push(Reverse((r - l + 1, r)));
            i += 1;
        }
        while let Some(&Reverse((_, end))) = heap.peek() {
            if end < q {
                heap.pop();
            } else {
                break;
            }
        }
        if let Some(&Reverse((size, _))) = heap.peek() {
            res[orig_idx] = size;
        }
    }
    res
}
