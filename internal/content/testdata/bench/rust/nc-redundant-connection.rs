fn find_redundant_connection(edges: Vec<Vec<i64>>) -> Vec<i64> {
    let n = edges.len();
    let mut parent: Vec<usize> = (0..=n).collect();
    let mut rank = vec![0i64; n + 1];

    fn find(parent: &mut Vec<usize>, mut x: usize) -> usize {
        while parent[x] != x {
            parent[x] = parent[parent[x]];
            x = parent[x];
        }
        x
    }

    for edge in edges.iter() {
        let u = edge[0] as usize;
        let v = edge[1] as usize;
        let mut pa = find(&mut parent, u);
        let mut pb = find(&mut parent, v);
        if pa == pb {
            return vec![edge[0], edge[1]];
        }
        if rank[pa] < rank[pb] {
            std::mem::swap(&mut pa, &mut pb);
        }
        parent[pb] = pa;
        if rank[pa] == rank[pb] {
            rank[pa] += 1;
        }
    }
    Vec::new()
}
