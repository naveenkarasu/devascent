use std::collections::BinaryHeap;
use std::cmp::Reverse;

fn running_medians(stream: Vec<i64>) -> Vec<i64> {
    let mut lo: BinaryHeap<i64> = BinaryHeap::new(); // max-heap (lower half)
    let mut hi: BinaryHeap<Reverse<i64>> = BinaryHeap::new(); // min-heap (upper half)
    let mut result: Vec<i64> = Vec::new();
    for num in stream {
        lo.push(num);
        let top = lo.pop().unwrap();
        hi.push(Reverse(top));
        if hi.len() > lo.len() {
            let Reverse(t) = hi.pop().unwrap();
            lo.push(t);
        }
        if lo.len() == hi.len() {
            let l = *lo.peek().unwrap();
            let Reverse(h) = *hi.peek().unwrap();
            result.push((l + h) / 2);
        } else {
            let l = *lo.peek().unwrap();
            result.push(l);
        }
    }
    result
}
