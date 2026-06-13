fn find(parent: &mut Vec<usize>, x: usize) -> usize {
    let mut x = x;
    while parent[x] != x {
        parent[x] = parent[parent[x]];
        x = parent[x];
    }
    x
}

fn count_components(n: i64, edges: Vec<Vec<i64>>) -> i64 {
    let nu = n as usize;
    let mut parent: Vec<usize> = (0..nu).collect();
    let mut merges = 0i64;
    for e in &edges {
        let a = e[0] as usize;
        let b = e[1] as usize;
        let pa = find(&mut parent, a);
        let pb = find(&mut parent, b);
        if pa != pb {
            parent[pa] = pb;
            merges += 1;
        }
    }
    n - merges
}
