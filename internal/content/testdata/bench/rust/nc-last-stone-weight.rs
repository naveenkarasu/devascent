use std::collections::BinaryHeap;

fn last_stone_weight(stones: Vec<i64>) -> i64 {
    let mut heap: BinaryHeap<i64> = BinaryHeap::new();
    for s in stones {
        heap.push(s);
    }
    while heap.len() > 1 {
        let y = heap.pop().unwrap();
        let x = heap.pop().unwrap();
        if x != y {
            heap.push(y - x);
        }
    }
    *heap.peek().unwrap_or(&0)
}
