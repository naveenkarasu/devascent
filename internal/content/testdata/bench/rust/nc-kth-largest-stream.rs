use std::collections::BinaryHeap;
use std::cmp::Reverse;

fn kth_largest_stream_ops(k: i64, nums: Vec<i64>, operations: Vec<Vec<J>>) -> Vec<i64> {
    let k = k as usize;
    // Min-heap of the k largest seen so far (heap root = kth largest).
    let mut heap: BinaryHeap<Reverse<i64>> = BinaryHeap::new();
    for &v in &nums {
        heap.push(Reverse(v));
    }
    while heap.len() > k {
        heap.pop();
    }
    let mut out: Vec<i64> = Vec::new();
    for op in &operations {
        // op = ["add", val]
        let val = if let J::Int(v) = &op[1] { *v } else { 0 };
        heap.push(Reverse(val));
        while heap.len() > k {
            heap.pop();
        }
        out.push(heap.peek().unwrap().0);
    }
    out
}
