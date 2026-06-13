use std::collections::BinaryHeap;
use std::cmp::Reverse;

fn find_kth_largest(nums: Vec<i64>, k: i64) -> i64 {
    let mut heap: BinaryHeap<Reverse<i64>> = BinaryHeap::new();
    for n in nums {
        heap.push(Reverse(n));
        if heap.len() as i64 > k {
            heap.pop();
        }
    }
    heap.peek().map(|Reverse(v)| *v).unwrap_or(0)
}
