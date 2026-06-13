fn combination_sum2(candidates: Vec<i64>, target: i64) -> Vec<Vec<i64>> {
    let mut candidates = candidates;
    candidates.sort();
    let mut res: Vec<Vec<i64>> = Vec::new();
    let mut current: Vec<i64> = Vec::new();

    fn backtrack(
        start: usize,
        remaining: i64,
        candidates: &Vec<i64>,
        current: &mut Vec<i64>,
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
            if i > start && candidates[i] == candidates[i - 1] {
                continue;
            }
            current.push(candidates[i]);
            backtrack(i + 1, remaining - candidates[i], candidates, current, res);
            current.pop();
        }
    }

    backtrack(0, target, &candidates, &mut current, &mut res);
    for c in res.iter_mut() {
        c.sort();
    }
    res.sort();
    res
}
