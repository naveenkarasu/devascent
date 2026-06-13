use std::collections::HashSet;

fn dfs(adj: &Vec<Vec<usize>>, visited: &mut HashSet<usize>, node: usize) {
    visited.insert(node);
    for &nb in &adj[node] {
        if !visited.contains(&nb) {
            dfs(adj, visited, nb);
        }
    }
}

fn valid_tree(n: i64, edges: Vec<Vec<i64>>) -> bool {
    if edges.len() as i64 != n - 1 {
        return false;
    }
    let nu = n as usize;
    let mut adj: Vec<Vec<usize>> = vec![Vec::new(); nu];
    for e in &edges {
        let u = e[0] as usize;
        let v = e[1] as usize;
        adj[u].push(v);
        adj[v].push(u);
    }
    let mut visited: HashSet<usize> = HashSet::new();
    if nu > 0 {
        dfs(&adj, &mut visited, 0);
    }
    visited.len() == nu
}
