fn trap_node(n: i64, edges: Vec<Vec<i64>>) -> i64 {
    let mut degree: Vec<i64> = vec![0; n as usize];
    for e in &edges {
        degree[e[0] as usize] += 1;
        degree[e[1] as usize] += 1;
    }
    degree.iter().cloned().max().unwrap_or(0) - 1
}
