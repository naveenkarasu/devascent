fn backtrack(current: &mut Vec<i64>, remaining: &Vec<i64>, res: &mut Vec<Vec<i64>>) {
    if remaining.is_empty() {
        res.push(current.clone());
        return;
    }
    for i in 0..remaining.len() {
        current.push(remaining[i]);
        let mut rest = remaining.clone();
        rest.remove(i);
        backtrack(current, &rest, res);
        current.pop();
    }
}

fn permute(nums: Vec<i64>) -> Vec<Vec<i64>> {
    let mut res: Vec<Vec<i64>> = Vec::new();
    let mut current: Vec<i64> = Vec::new();
    backtrack(&mut current, &nums, &mut res);
    res.sort();
    res
}
