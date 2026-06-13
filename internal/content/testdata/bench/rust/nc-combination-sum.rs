fn backtrack(
    start: usize,
    candidates: &Vec<i64>,
    current: &mut Vec<i64>,
    remaining: i64,
    res: &mut Vec<Vec<i64>>,
) {
    if remaining == 0 {
        res.push(current.clone());
        return;
    }
    for i in start..candidates.len() {
        if candidates[i] > remaining {
            break;
        }
        current.push(candidates[i]);
        backtrack(i, candidates, current, remaining - candidates[i], res);
        current.pop();
    }
}

fn combination_sum(candidates: Vec<i64>, target: i64) -> Vec<Vec<i64>> {
    let mut candidates = candidates;
    candidates.sort();
    let mut res: Vec<Vec<i64>> = Vec::new();
    let mut current: Vec<i64> = Vec::new();
    backtrack(0, &candidates, &mut current, target, &mut res);
    for c in res.iter_mut() {
        c.sort();
    }
    res.sort();
    res
}
