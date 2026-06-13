fn merge_k_lists(lists: Vec<Vec<i64>>) -> Vec<i64> {
    use std::collections::BinaryHeap;
    use std::cmp::Reverse;
    // heap of (value, list_index, elem_index)
    let mut heap: BinaryHeap<Reverse<(i64, usize, usize)>> = BinaryHeap::new();
    for (i, lst) in lists.iter().enumerate() {
        if !lst.is_empty() {
            heap.push(Reverse((lst[0], i, 0)));
        }
    }
    let mut out: Vec<i64> = Vec::new();
    while let Some(Reverse((val, i, j))) = heap.pop() {
        out.push(val);
        if j + 1 < lists[i].len() {
            heap.push(Reverse((lists[i][j + 1], i, j + 1)));
        }
    }
    out
}
