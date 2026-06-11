use std::collections::BinaryHeap;
use std::cmp::Reverse;

fn median_stream_ops(operations: Vec<Vec<J>>) -> Vec<J> {
    // lo: max-heap (lower half), hi: min-heap (upper half).
    let mut lo: BinaryHeap<i64> = BinaryHeap::new();
    let mut hi: BinaryHeap<Reverse<i64>> = BinaryHeap::new();
    let mut out: Vec<J> = Vec::new();

    for op in &operations {
        let name = match &op[0] {
            J::Str(s) => s.as_str(),
            _ => "",
        };
        if name == "addNum" {
            let v = if let J::Int(x) = &op[1] { *x } else { 0 };
            lo.push(v);
            // move top of lo to hi
            let t = lo.pop().unwrap();
            hi.push(Reverse(t));
            if hi.len() > lo.len() {
                let Reverse(t2) = hi.pop().unwrap();
                lo.push(t2);
            }
            out.push(J::Null);
        } else {
            // findMedian
            if lo.len() > hi.len() {
                out.push(J::Int(*lo.peek().unwrap()));
            } else {
                let a = *lo.peek().unwrap();
                let b = hi.peek().unwrap().0;
                let sum = a + b;
                if sum % 2 == 0 {
                    out.push(J::Int(sum / 2));
                } else {
                    out.push(J::Flt(sum as f64 / 2.0));
                }
            }
        }
    }
    out
}
