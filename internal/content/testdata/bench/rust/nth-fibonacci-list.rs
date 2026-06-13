fn fibonacci_list(n: i64) -> Vec<i64> {
    if n == 0 {
        return vec![];
    }
    if n == 1 {
        return vec![0];
    }
    let mut result: Vec<i64> = vec![0, 1];
    for _ in 2..n {
        let len = result.len();
        result.push(result[len - 1] + result[len - 2]);
    }
    result
}
