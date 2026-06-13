use std::collections::{HashMap, BinaryHeap};
use std::cmp::Reverse;

fn network_delay_time(times: Vec<Vec<i64>>, n: i64, k: i64) -> i64 {
    let mut graph: HashMap<i64, Vec<(i64, i64)>> = HashMap::new();
    for t in &times {
        let (u, v, w) = (t[0], t[1], t[2]);
        graph.entry(u).or_insert_with(Vec::new).push((w, v));
    }
    let mut dist: HashMap<i64, i64> = HashMap::new();
    let mut heap: BinaryHeap<Reverse<(i64, i64)>> = BinaryHeap::new();
    heap.push(Reverse((0, k)));
    while let Some(Reverse((d, node))) = heap.pop() {
        if dist.contains_key(&node) {
            continue;
        }
        dist.insert(node, d);
        if let Some(neighbors) = graph.get(&node) {
            for &(weight, nb) in neighbors {
                if !dist.contains_key(&nb) {
                    heap.push(Reverse((d + weight, nb)));
                }
            }
        }
    }
    if dist.len() as i64 == n {
        *dist.values().max().unwrap()
    } else {
        -1
    }
}
