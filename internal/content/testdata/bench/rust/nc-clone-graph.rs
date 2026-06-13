fn clone_graph(adj: Vec<Vec<i64>>) -> Vec<Vec<i64>> {
    if adj.is_empty() {
        return Vec::new();
    }
    // Deep-copy semantics are value-equivalent to returning a normalized
    // adjacency list: for each node i (1-indexed), the sorted neighbor labels.
    let n = adj.len();
    let mut out: Vec<Vec<i64>> = vec![Vec::new(); n];
    for i in 0..n {
        let mut nbrs = adj[i].clone();
        nbrs.sort();
        out[i] = nbrs;
    }
    out
}
