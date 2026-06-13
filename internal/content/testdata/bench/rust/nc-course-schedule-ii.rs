use std::collections::BinaryHeap;
use std::cmp::Reverse;

fn find_order(num_courses: i64, prerequisites: Vec<Vec<i64>>) -> Vec<i64> {
    let nc = num_courses as usize;
    let mut in_degree = vec![0i64; nc];
    let mut adj: Vec<Vec<usize>> = vec![Vec::new(); nc];
    for pre in prerequisites.iter() {
        let a = pre[0] as usize;
        let b = pre[1] as usize;
        adj[b].push(a);
        in_degree[a] += 1;
    }
    let mut heap: BinaryHeap<Reverse<i64>> = BinaryHeap::new();
    for i in 0..nc {
        if in_degree[i] == 0 {
            heap.push(Reverse(i as i64));
        }
    }
    let mut order: Vec<i64> = Vec::new();
    while let Some(Reverse(node)) = heap.pop() {
        order.push(node);
        for &nei in adj[node as usize].iter() {
            in_degree[nei] -= 1;
            if in_degree[nei] == 0 {
                heap.push(Reverse(nei as i64));
            }
        }
    }
    if order.len() == nc {
        order
    } else {
        Vec::new()
    }
}
