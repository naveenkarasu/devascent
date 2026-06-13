fn dfs(c: usize, graph: &Vec<Vec<usize>>, state: &mut Vec<u8>) -> bool {
    if state[c] == 1 {
        return false;
    }
    if state[c] == 2 {
        return true;
    }
    state[c] = 1;
    for &nxt in &graph[c] {
        if !dfs(nxt, graph, state) {
            return false;
        }
    }
    state[c] = 2;
    true
}

fn can_finish(num_courses: i64, prerequisites: Vec<Vec<i64>>) -> bool {
    let n = num_courses as usize;
    let mut graph: Vec<Vec<usize>> = vec![Vec::new(); n];
    for pre in &prerequisites {
        let a = pre[0] as usize;
        let b = pre[1] as usize;
        graph[a].push(b);
    }
    let mut state = vec![0u8; n];
    for c in 0..n {
        if !dfs(c, &graph, &mut state) {
            return false;
        }
    }
    true
}
